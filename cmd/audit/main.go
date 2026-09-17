package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"rx-dispatch/internal/contracts"
	"rx-dispatch/internal/models"
)

const (
	serviceName = "rx-audit"
	defaultPort = "8088"
	version     = "0.1.0-scaffold"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	host := os.Getenv("HOST")
	// Per REQ-SEC-001, internal services must default to loopback.
	if host == "" {
		host = "127.0.0.1"
	}

	log.Printf("[%s] Starting service scaffold on port %s", serviceName, port)

	http.HandleFunc("/healthz", healthCheckHandler)
	http.HandleFunc("/audit/event", eventHandler)

	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("[%s] Listening on %s", serviceName, addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("[%s] Failed to start server: %v", serviceName, err)
	}
}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB limit
	defer r.Body.Close()

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			http.Error(
				w,
				http.StatusText(http.StatusRequestEntityTooLarge),
				http.StatusRequestEntityTooLarge,
			)
			return
		}
		// Otro error de lectura, menos común pero posible (ej. error de red).
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var req contracts.RecordAuditEventRequest
	dec := json.NewDecoder(bytes.NewReader(payload))
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		http.Error(w, "Invalid or malformed JSON payload", http.StatusBadRequest)
		return
	}

	// Validación básica para asegurar que el evento no está vacío.
	if req.Event == (models.AuditEvent{}) || req.Event.EventAction == "" {
		http.Error(w, "Invalid audit event payload: event data is missing or incomplete", http.StatusBadRequest)
		return
	}

	// Registrar el evento recibido. En una implementación futura, esto se escribiría en un log persistente.
	log.Printf("[%s] AUDIT EVENT RECEIVED: Action=%s, User=%s, Outcome=%s, SourceIP=%s, StudyUID=%s",
		serviceName,
		req.Event.EventAction,
		req.Event.UserID,
		req.Event.EventOutcome,
		req.Event.SourceIP,
		req.Event.StudyInstanceUID,
	)

	w.WriteHeader(http.StatusAccepted)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	response := contracts.HealthResponse{
		Service:   serviceName,
		Status:    "UP",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Version:   version,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
