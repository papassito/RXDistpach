package dicom

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

// SendCStore sends a raw DICOM dataset to a remote C-STORE SCP server via TCP.
func SendCStore(host string, port int, callingAE, calledAE string, dicomBytes []byte) error {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to DICOM SCP on %s: %w", addr, err)
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(10 * time.Second))

	// 1. Send A-ASSOCIATE-RQ
	assocRQ := buildAssociateRQ(callingAE, calledAE)
	if _, err := conn.Write(assocRQ); err != nil {
		return fmt.Errorf("failed to write A-ASSOCIATE-RQ: %w", err)
	}

	// 2. Read A-ASSOCIATE-AC
	pduHeader := make([]byte, 6)
	if _, err := io.ReadFull(conn, pduHeader); err != nil {
		return fmt.Errorf("failed to read A-ASSOCIATE-AC header: %w", err)
	}
	if pduHeader[0] != 0x02 {
		return fmt.Errorf("expected A-ASSOCIATE-AC (0x02), got 0x%02x", pduHeader[0])
	}
	acLen := binary.BigEndian.Uint32(pduHeader[2:6])
	acBody := make([]byte, acLen)
	if _, err := io.ReadFull(conn, acBody); err != nil {
		return fmt.Errorf("failed to read A-ASSOCIATE-AC body: %w", err)
	}

	// 3. Send P-DATA-TF with C-STORE-RQ command and Dataset
	pdataPDU := buildPDataTF(dicomBytes)
	if _, err := conn.Write(pdataPDU); err != nil {
		return fmt.Errorf("failed to write P-DATA-TF: %w", err)
	}

	// 4. Read C-STORE-RSP
	rspHeader := make([]byte, 6)
	if _, err := io.ReadFull(conn, rspHeader); err != nil {
		return fmt.Errorf("failed to read C-STORE-RSP header: %w", err)
	}
	rspLen := binary.BigEndian.Uint32(rspHeader[2:6])
	rspBody := make([]byte, rspLen)
	if _, err := io.ReadFull(conn, rspBody); err != nil {
		return fmt.Errorf("failed to read C-STORE-RSP body: %w", err)
	}

	// 5. Send A-RELEASE-RQ
	releaseRQ := []byte{0x05, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00}
	if _, err := conn.Write(releaseRQ); err != nil {
		return fmt.Errorf("failed to write A-RELEASE-RQ: %w", err)
	}

	// 6. Read A-RELEASE-RP
	relHeader := make([]byte, 6)
	if _, err := io.ReadFull(conn, relHeader); err != nil {
		return fmt.Errorf("failed to read A-RELEASE-RP header: %w", err)
	}
	if relHeader[0] != 0x06 {
		return fmt.Errorf("expected A-RELEASE-RP (0x06), got 0x%02x", relHeader[0])
	}
	relLen := binary.BigEndian.Uint32(relHeader[2:6])
	relBody := make([]byte, relLen)
	_, _ = io.ReadFull(conn, relBody)

	return nil
}

func buildAssociateRQ(callingAE, calledAE string) []byte {
	buf := new(bytes.Buffer)
	// Protocol Version 1
	buf.Write([]byte{0x00, 0x01, 0x00, 0x00})
	// Called AE (16 chars)
	buf.WriteString(fmt.Sprintf("%-16s", calledAE)[:16])
	// Calling AE (16 chars)
	buf.WriteString(fmt.Sprintf("%-16s", callingAE)[:16])
	// Reserved 32 bytes
	buf.Write(make([]byte, 32))

	// Presentation Context: Digital X-Ray Image Storage (1.2.840.10008.5.1.4.1.1.1.1)
	sopClass := "1.2.840.10008.5.1.4.1.1.1.1"
	pcBuf := new(bytes.Buffer)
	pcBuf.WriteByte(0x01) // Context ID
	pcBuf.WriteByte(0x00) // Reserved
	pcBuf.WriteByte(0x00) // Reserved
	pcBuf.WriteByte(0x00) // Reserved

	// Abstract Syntax Sub-Item (0x30)
	asItem := new(bytes.Buffer)
	asItem.WriteByte(0x30)
	asItem.WriteByte(0x00)
	_ = binary.Write(asItem, binary.BigEndian, uint16(len(sopClass)))
	asItem.WriteString(sopClass)
	pcBuf.Write(asItem.Bytes())

	// Transfer Syntax Sub-Item: Implicit VR Little Endian (0x40)
	tsUID := "1.2.840.10008.1.2"
	tsItem := new(bytes.Buffer)
	tsItem.WriteByte(0x40)
	tsItem.WriteByte(0x00)
	_ = binary.Write(tsItem, binary.BigEndian, uint16(len(tsUID)))
	tsItem.WriteString(tsUID)
	pcBuf.Write(tsItem.Bytes())

	// Add PC Item to RQ
	buf.WriteByte(0x20) // PC Item Type
	buf.WriteByte(0x00)
	_ = binary.Write(buf, binary.BigEndian, uint16(pcBuf.Len()))
	buf.Write(pcBuf.Bytes())

	// Full PDU (Type 1)
	pdu := new(bytes.Buffer)
	pdu.WriteByte(0x01)
	pdu.WriteByte(0x00)
	_ = binary.Write(pdu, binary.BigEndian, uint32(buf.Len()))
	pdu.Write(buf.Bytes())

	return pdu.Bytes()
}

func buildPDataTF(data []byte) []byte {
	pdvBuf := new(bytes.Buffer)
	pdvLen := uint32(len(data) + 2)
	_ = binary.Write(pdvBuf, binary.BigEndian, pdvLen)
	pdvBuf.WriteByte(0x01) // Context ID
	pdvBuf.WriteByte(0x02) // Message Control Header: Last Fragment + Data
	pdvBuf.Write(data)

	pdu := new(bytes.Buffer)
	pdu.WriteByte(0x04)
	pdu.WriteByte(0x00)
	_ = binary.Write(pdu, binary.BigEndian, uint32(pdvBuf.Len()))
	pdu.Write(pdvBuf.Bytes())

	return pdu.Bytes()
}

// GenerateSyntheticDICOM creates a valid test DICOM payload with standard clinical tags.
func GenerateSyntheticDICOM(patientName, patientID, accession, studyDesc string) []byte {
	buf := new(bytes.Buffer)

	// 128 bytes preamble
	buf.Write(make([]byte, 128))
	// "DICM" magic prefix
	buf.WriteString("DICM")

	// Helper to write element (Implicit VR Little Endian)
	writeElem := func(group, element uint16, val string) {
		valBytes := []byte(val)
		if len(valBytes)%2 != 0 {
			valBytes = append(valBytes, ' ') // Pad to even bytes
		}
		_ = binary.Write(buf, binary.LittleEndian, group)
		_ = binary.Write(buf, binary.LittleEndian, element)
		_ = binary.Write(buf, binary.LittleEndian, uint32(len(valBytes)))
		buf.Write(valBytes)
	}

	writeElem(0x0010, 0x0010, patientName)
	writeElem(0x0010, 0x0020, patientID)
	writeElem(0x0008, 0x0050, accession)
	writeElem(0x0008, 0x0060, "DX")
	writeElem(0x0008, 0x1030, studyDesc)
	writeElem(0x0018, 0x0015, "CHEST")
	writeElem(0x0020, 0x000D, fmt.Sprintf("1.2.840.10008.%d", time.Now().UnixNano()))
	writeElem(0x0020, 0x000E, fmt.Sprintf("1.2.840.10008.2.%d", time.Now().UnixNano()))

	// Mock pixel data (7FE0, 0010)
	pixelGroup := uint16(0x7FE0)
	pixelElem := uint16(0x0010)
	mockPixels := bytes.Repeat([]byte{0x7A, 0x82, 0x9F, 0xB0}, 64) // 256 bytes pixel stream
	_ = binary.Write(buf, binary.LittleEndian, pixelGroup)
	_ = binary.Write(buf, binary.LittleEndian, pixelElem)
	_ = binary.Write(buf, binary.LittleEndian, uint32(len(mockPixels)))
	buf.Write(mockPixels)

	return buf.Bytes()
}
