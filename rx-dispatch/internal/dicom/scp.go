package dicom

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

// IngestionCallback is invoked when an authentic DICOM study is received over TCP C-STORE.
type IngestionCallback func(dataset *Dataset) error

// ServerConfig configures the native Go DICOM C-STORE SCP.
type ServerConfig struct {
	AETitle string // Default: RX_DISPATCH
	Port    int    // Default: 11112 (or 104)
}

// SCPServer implements a native DICOM C-STORE SCP listener.
type SCPServer struct {
	config            ServerConfig
	listener          net.Listener
	mu                sync.RWMutex
	receivedInstances int
	lastReceived      *Dataset
	onIngest          IngestionCallback
	stopChan          chan struct{}
}

// NewSCPServer creates a new DICOM C-STORE SCP listener.
func NewSCPServer(cfg ServerConfig, onIngest IngestionCallback) *SCPServer {
	if cfg.AETitle == "" {
		cfg.AETitle = "RX_DISPATCH"
	}
	if cfg.Port == 0 {
		cfg.Port = 11112
	}

	return &SCPServer{
		config:   cfg,
		onIngest: onIngest,
		stopChan: make(chan struct{}),
	}
}

// Start begins listening on TCP port.
func (s *SCPServer) Start() error {
	addr := fmt.Sprintf("0.0.0.0:%d", s.config.Port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to bind DICOM listener on %s: %w", addr, err)
	}
	s.listener = l
	log.Printf("[DICOM-SCP] Listening on %s (AET: %s, Protocol: C-STORE SCP)", addr, s.config.AETitle)

	go s.acceptLoop()
	return nil
}

// Stop closes the listener.
func (s *SCPServer) Stop() {
	close(s.stopChan)
	if s.listener != nil {
		_ = s.listener.Close()
	}
}

// Stats returns receiver metrics.
func (s *SCPServer) Stats() (int, *Dataset) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.receivedInstances, s.lastReceived
}

func (s *SCPServer) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.stopChan:
				return
			default:
				log.Printf("[DICOM-SCP] Accept error: %v", err)
				return
			}
		}
		go s.handleConnection(conn)
	}
}

func (s *SCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(60 * time.Second))

	var collectedData []byte

	for {
		// Read DICOM Upper Layer Protocol (ULP) PDU Header: Type (1 byte) + Reserved (1 byte) + Length (4 bytes)
		header := make([]byte, 6)
		_, err := io.ReadFull(conn, header)
		if err != nil {
			break
		}

		pduType := header[0]
		pduLen := binary.BigEndian.Uint32(header[2:6])

		pduBody := make([]byte, pduLen)
		_, err = io.ReadFull(conn, pduBody)
		if err != nil {
			log.Printf("[DICOM-SCP] Error reading PDU body: %v", err)
			break
		}

		switch pduType {
		case 0x01: // A-ASSOCIATE-RQ
			log.Printf("[DICOM-SCP] Received A-ASSOCIATE-RQ (%d bytes)", pduLen)
			acResponse := s.buildAssociateAC(pduBody)
			_, _ = conn.Write(acResponse)

		case 0x04: // P-DATA-TF (carries C-STORE-RQ command set and dataset)
			if len(pduBody) > 6 {
				// PDV item: 4 bytes length, 1 byte context ID, 1 byte flags
				pdvData := pduBody[6:]
				collectedData = append(collectedData, pdvData...)

				// Send C-STORE-RSP (Success: Status 0x0000)
				cStoreRsp := s.buildCStoreRSP()
				_, _ = conn.Write(cStoreRsp)
			}

		case 0x05: // A-RELEASE-RQ
			log.Printf("[DICOM-SCP] Received A-RELEASE-RQ. Sending A-RELEASE-RP.")
			releaseRsp := []byte{0x06, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00}
			_, _ = conn.Write(releaseRsp)

			// Process the accumulated DICOM dataset
			if len(collectedData) > 0 {
				ds, err := ParseDICOMBytes(collectedData)
				if err != nil {
					log.Printf("[DICOM-SCP] Warning parsing received DICOM: %v", err)
				} else {
					s.mu.Lock()
					s.receivedInstances++
					s.lastReceived = ds
					s.mu.Unlock()

					log.Printf("[DICOM-SCP] Successfully received instance from modality: Patient='%s', Modality='%s', Study='%s', Size=%d bytes, SHA=%s",
						ds.PatientName, ds.Modality, ds.StudyDescription, ds.SizeBytes, ds.ChecksumSHA[:12])

					if s.onIngest != nil {
						if err := s.onIngest(ds); err != nil {
							log.Printf("[DICOM-SCP] Ingest handler error: %v", err)
						}
					}
				}
			}
			return

		case 0x07: // A-ABORT
			log.Printf("[DICOM-SCP] Modality aborted association.")
			return
		}
	}
}

func (s *SCPServer) buildAssociateAC(rqBody []byte) []byte {
	buf := new(bytes.Buffer)

	// Upper Layer Protocol Version (0x0001)
	buf.Write([]byte{0x00, 0x01})
	// Reserved
	buf.Write([]byte{0x00, 0x00})

	// Called AE Title (16 bytes)
	calledAE := fmt.Sprintf("%-16s", s.config.AETitle)
	buf.WriteString(calledAE[:16])

	// Calling AE Title (16 bytes, copied from RQ or default)
	callingAE := "MODALITY_SCU    "
	if len(rqBody) >= 32 {
		callingAE = string(rqBody[16:32])
	}
	buf.WriteString(fmt.Sprintf("%-16s", strings.TrimSpace(callingAE))[:16])

	// Reserved (32 bytes)
	buf.Write(make([]byte, 32))

	// Presentation Context Item (Acceptance: 0x21)
	pcBuf := new(bytes.Buffer)
	pcBuf.WriteByte(0x01) // Context ID
	pcBuf.WriteByte(0x00) // Reserved
	pcBuf.WriteByte(0x00) // Result/Reason: 0 = Acceptance
	pcBuf.WriteByte(0x00) // Reserved

	// Transfer Syntax Sub-Item: Implicit VR Little Endian (1.2.840.10008.1.2)
	tsUID := "1.2.840.10008.1.2"
	tsItem := new(bytes.Buffer)
	tsItem.WriteByte(0x40) // Sub-item type
	tsItem.WriteByte(0x00) // Reserved
	tsLen := uint16(len(tsUID))
	_ = binary.Write(tsItem, binary.BigEndian, tsLen)
	tsItem.WriteString(tsUID)

	pcBuf.Write(tsItem.Bytes())

	// Write Presentation Context Header
	buf.WriteByte(0x21)
	buf.WriteByte(0x00)
	_ = binary.Write(buf, binary.BigEndian, uint16(pcBuf.Len()))
	buf.Write(pcBuf.Bytes())

	// Build full PDU (Type 2: A-ASSOCIATE-AC)
	pdu := new(bytes.Buffer)
	pdu.WriteByte(0x02)
	pdu.WriteByte(0x00)
	_ = binary.Write(pdu, binary.BigEndian, uint32(buf.Len()))
	pdu.Write(buf.Bytes())

	return pdu.Bytes()
}

func (s *SCPServer) buildCStoreRSP() []byte {
	// PDU Type 4: P-DATA-TF containing C-STORE-RSP command
	// Minimal C-STORE-RSP: Affected SOP Class UID, Command Field 0x8001, Message ID 1, Status 0x0000 (Success)
	cmdBuf := new(bytes.Buffer)
	// (0000,0100) Command Field = 0x8001 (C-STORE-RSP)
	cmdBuf.Write([]byte{0x00, 0x00, 0x00, 0x01, 0x02, 0x00, 0x00, 0x00, 0x01, 0x80})
	// (0000,0120) Message ID Being Responded To = 1
	cmdBuf.Write([]byte{0x00, 0x00, 0x20, 0x01, 0x02, 0x00, 0x00, 0x00, 0x01, 0x00})
	// (0000,0900) Status = 0x0000 (Success)
	cmdBuf.Write([]byte{0x00, 0x00, 0x00, 0x09, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00})

	// PDV Header
	pdvBuf := new(bytes.Buffer)
	pdvLen := uint32(cmdBuf.Len() + 2)
	_ = binary.Write(pdvBuf, binary.BigEndian, pdvLen)
	pdvBuf.WriteByte(0x01) // Context ID
	pdvBuf.WriteByte(0x03) // Message Control Header: Last Fragment + Command
	pdvBuf.Write(cmdBuf.Bytes())

	// Full PDU
	pdu := new(bytes.Buffer)
	pdu.WriteByte(0x04)
	pdu.WriteByte(0x00)
	_ = binary.Write(pdu, binary.BigEndian, uint32(pdvBuf.Len()))
	pdu.Write(pdvBuf.Bytes())

	return pdu.Bytes()
}
