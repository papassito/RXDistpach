package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"rx-dispatch/internal/config"
	"rx-dispatch/internal/contracts"
	"rx-dispatch/internal/models"
	"rx-dispatch/internal/transport"
)

type AuditServer struct {
	mu     sync.RWMutex
	events []models.AuditEvent
	port   int
}

func NewAuditServer(port int) *AuditServer {
	return &AuditServer{
		events: make([]models.AuditEvent, 0),
		port:   port,
	}
}

func (s *AuditServer) handleHealthz(w http.ResponseWriter, r *http.Request) {
	resp := contracts.HealthResponse{
		Service:   "rx-audit",
		Port:      s.port,
		Status:    "UP",
		Timestamp: time.Now(),
		Version:   "2.0.0-decoupled",
	}
	transport.WriteJSON(w, http.StatusOK, resp)
}

func (s *AuditServer) handleEvents(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var req contracts.RecordAuditEventRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			transport.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid payload: %v", err), "rx-audit")
			return
		}

		eventID := fmt.Sprintf("AUD-%d-%04d", time.Now().UnixNano()/1e6, len(s.events)+1)
		event := models.AuditEvent{
			ID:        eventID,
			Timestamp: time.Now(),
			EventType: req.EventType,
			StudyID:   req.StudyID,
			Service:   req.Service,
			Actor:     req.Actor,
			Details:   req.Details,
		}

		s.mu.Lock()
		// Prepend to maintain newest-first order
		s.events = append([]models.AuditEvent{event}, s.events...)
		s.mu.Unlock()

		log.Printf("[RX-AUDIT] Recorded event %s [%s] for study %s", event.ID, event.EventType, event.StudyID)
		transport.WriteJSON(w, http.StatusCreated, contracts.RecordAuditEventResponse{
			EventID: eventID,
			Success: true,
		})

	case http.MethodGet:
		studyID := r.URL.Query().Get("studyId")
		s.mu.RLock()
		defer s.mu.RUnlock()

		filtered := make([]models.AuditEvent, 0)
		for _, e := range s.events {
			if studyID == "" || e.StudyID == studyID {
				filtered = append(filtered, e)
			}
		}

		transport.WriteJSON(w, http.StatusOK, contracts.QueryAuditEventsResponse{
			Events: filtered,
			Count:  len(filtered),
		})

	default:
		transport.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", "rx-audit")
	}
}

func main() {
	port := config.GetEnvInt("RX_AUDIT_PORT", 8088)
	server := NewAuditServer(port)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", server.handleHealthz)
	mux.HandleFunc("/audit/events", server.handleEvents)

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	log.Printf("[RX-AUDIT] Starting on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("[RX-AUDIT] Server failed: %v", err)
	}
}
