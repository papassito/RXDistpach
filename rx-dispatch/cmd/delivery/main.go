package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"rx-dispatch/internal/config"
	"rx-dispatch/internal/contracts"
	"rx-dispatch/internal/delivery"
	"rx-dispatch/internal/models"
	"rx-dispatch/internal/transport"
)

type DeliveryServer struct {
	mu         sync.RWMutex
	deliveries map[string][]models.DeliveryRecord
	port       int
	topology   config.TopologyConfig
	queue      *delivery.StoreAndForwardQueue
}

func NewDeliveryServer(port int, topology config.TopologyConfig) *DeliveryServer {
	srv := &DeliveryServer{
		deliveries: make(map[string][]models.DeliveryRecord),
		port:       port,
		topology:   topology,
	}

	dataDir := config.GetEnv("RX_DATA_DIR", "/app/applet/rx-dispatch/data")
	srv.queue = delivery.NewStoreAndForwardQueue(
		dataDir,
		func(item *delivery.QueueItem) {
			// Delivery success callback: save to in-memory records
			rec := item.ToDeliveryRecord()
			srv.mu.Lock()
			srv.deliveries[rec.StudyID] = append([]models.DeliveryRecord{rec}, srv.deliveries[rec.StudyID]...)
			srv.mu.Unlock()

			// Update study status to DELIVERED
			updateStatusReq := contracts.UpdateStudyStatusRequest{
				Status: models.StudyStatusDelivered,
				Reason: fmt.Sprintf("Dispatched via Store-and-Forward (%s to %s)", item.Channel, item.Recipient),
			}
			_ = transport.PostJSON(fmt.Sprintf("%s/study/status/%s", srv.topology.StudyURL, item.StudyID), updateStatusReq, nil)
		},
		func(eventType, studyID, actor, details string) {
			// Audit callback
			auditReq := contracts.RecordAuditEventRequest{
				EventType: eventType,
				StudyID:   studyID,
				Service:   "rx-delivery",
				Actor:     actor,
				Details:   details,
			}
			_ = transport.PostJSON(fmt.Sprintf("%s/audit/events", srv.topology.AuditURL), auditReq, nil)
		},
	)

	return srv
}

func (s *DeliveryServer) handleHealthz(w http.ResponseWriter, r *http.Request) {
	resp := contracts.HealthResponse{
		Service:   "rx-delivery",
		Port:      s.port,
		Status:    "UP",
		Timestamp: time.Now(),
		Version:   "2.0.0-decoupled",
	}
	transport.WriteJSON(w, http.StatusOK, resp)
}

func (s *DeliveryServer) handleDispatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		transport.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", "rx-delivery")
		return
	}

	var req contracts.DispatchDeliveryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		transport.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid payload: %v", err), "rx-delivery")
		return
	}

	deliveryID := fmt.Sprintf("DEL-%d", time.Now().UnixNano()/1e6)
	randBytes := make([]byte, 6)
	_, _ = rand.Read(randBytes)
	trackingToken := fmt.Sprintf("TRK-%s", hex.EncodeToString(randBytes))

	channel := req.Channel
	if channel == "" {
		channel = "EMAIL"
	}
	recipient := req.Recipient
	if recipient == "" {
		recipient = req.PackageID.Study.Patient.Email
	}

	record := models.DeliveryRecord{
		ID:            deliveryID,
		PackageID:     req.PackageID.ID,
		StudyID:       req.PackageID.StudyID,
		Recipient:     recipient,
		Channel:       channel,
		Status:        models.DeliverySent,
		SentAt:        time.Now(),
		TrackingToken: trackingToken,
	}

	s.mu.Lock()
	s.deliveries[record.StudyID] = append([]models.DeliveryRecord{record}, s.deliveries[record.StudyID]...)
	s.mu.Unlock()

	// Notify RX Audit
	auditReq := contracts.RecordAuditEventRequest{
		EventType: "delivery_dispatched",
		StudyID:   record.StudyID,
		Service:   "rx-delivery",
		Actor:     "DELIVERY_ENGINE",
		Details:   fmt.Sprintf("Resultado despachado vía %s a %s. Token de rastreo: %s", channel, recipient, trackingToken),
	}
	_ = transport.PostJSON(fmt.Sprintf("%s/audit/events", s.topology.AuditURL), auditReq, nil)

	// Update study status in RX Study to DELIVERED
	updateStatusReq := contracts.UpdateStudyStatusRequest{
		Status: models.StudyStatusDelivered,
		Reason: fmt.Sprintf("Dispatched via %s to %s", channel, recipient),
	}
	_ = transport.PostJSON(fmt.Sprintf("%s/studies/%s/status", s.topology.StudyURL, record.StudyID), updateStatusReq, nil)

	resp := contracts.DispatchDeliveryResponse{
		DeliveryID:    deliveryID,
		TrackingToken: trackingToken,
		Status:        models.DeliverySent,
		Recipient:     recipient,
		Channel:       channel,
	}

	log.Printf("[RX-DELIVERY] Delivery %s sent to %s (Token: %s)", deliveryID, recipient, trackingToken)
	transport.WriteJSON(w, http.StatusOK, resp)
}

func (s *DeliveryServer) handleRecords(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		transport.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", "rx-delivery")
		return
	}

	studyID := r.URL.Query().Get("studyId")
	s.mu.RLock()
	defer s.mu.RUnlock()

	results := make([]models.DeliveryRecord, 0)
	if studyID != "" {
		results = append(results, s.deliveries[studyID]...)
	} else {
		for _, list := range s.deliveries {
			results = append(results, list...)
		}
	}

	transport.WriteJSON(w, http.StatusOK, contracts.QueryDeliveriesResponse{
		Deliveries: results,
		Count:      len(results),
	})
}

func (s *DeliveryServer) handleQueue(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		status := s.queue.GetStatus()
		transport.WriteJSON(w, http.StatusOK, status)

	case http.MethodPost:
		var req contracts.DispatchDeliveryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			transport.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid payload: %v", err), "rx-delivery")
			return
		}

		recipient := req.Recipient
		if recipient == "" {
			recipient = req.PackageID.Study.Patient.Email
		}
		channel := req.Channel
		if channel == "" {
			channel = "EMAIL"
		}

		item := s.queue.Enqueue(
			req.PackageID.StudyID,
			req.PackageID.ID,
			recipient,
			channel,
			req.PackageID.ChecksumSHA,
		)

		transport.WriteJSON(w, http.StatusAccepted, map[string]interface{}{
			"message":       "Delivery package accepted in resilient Store-and-Forward queue",
			"queueItem":     item,
			"trackingToken": item.TrackingToken,
			"status":        item.Status,
		})
	default:
		transport.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", "rx-delivery")
	}
}

func (s *DeliveryServer) handleNetworkToggle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		transport.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", "rx-delivery")
		return
	}

	var req struct {
		SimulateDrop bool `json:"simulateDrop"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	s.queue.SetNetworkSimulation(req.SimulateDrop)

	transport.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"simulatedDropActive": req.SimulateDrop,
		"status":              s.queue.GetStatus(),
	})
}

func (s *DeliveryServer) handleRetryAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		transport.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", "rx-delivery")
		return
	}

	retried := s.queue.RetryAll()
	transport.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"retriedCount": retried,
		"status":       s.queue.GetStatus(),
	})
}

func main() {
	port := config.GetEnvInt("RX_DELIVERY_PORT", 8087)
	topology := config.LoadTopologyConfig()
	server := NewDeliveryServer(port, topology)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", server.handleHealthz)
	mux.HandleFunc("/delivery/dispatch", server.handleDispatch)
	mux.HandleFunc("/delivery/records", server.handleRecords)
	mux.HandleFunc("/delivery/queue", server.handleQueue)
	mux.HandleFunc("/delivery/network-toggle", server.handleNetworkToggle)
	mux.HandleFunc("/delivery/retry-all", server.handleRetryAll)

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	log.Printf("[RX-DELIVERY] Starting on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("[RX-DELIVERY] Server failed: %v", err)
	}
}
