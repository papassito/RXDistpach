package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"rx-dispatch/internal/contracts"
	"rx-dispatch/internal/models"
	"rx-dispatch/internal/shared"
	"rx-dispatch/internal/transport"
)

const (
	serviceName = "rx-reader"
	servicePort = 8085
	version     = "0.1.0-scaffold" // Note: Version will be updated later
)

// server encapsula las dependencias para los manejadores HTTP, permitiendo la inyección de dependencias.
type server struct {
	auditClient  transport.AuditClienter
	resultClient transport.ResultClienter
}

func main() {
	log.Printf("[%s] Starting service on port %d", serviceName, servicePort)

	// Crear el servidor con sus dependencias.
	srv := &server{
		auditClient:  transport.NewAuditClient(),
		resultClient: transport.NewResultClient(),
	}

	http.HandleFunc("/healthz", srv.healthCheckHandler)
	http.HandleFunc("/reader/analyze", srv.analyzeHandler)

	addr := fmt.Sprintf("127.0.0.1:%d", servicePort)
	log.Printf("[%s] Listening on %s", serviceName, addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("[%s] FATAL: Failed to start server: %v", serviceName, err)
	}
}

func (s *server) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := contracts.HealthResponse{
		Service:   serviceName,
		Port:      servicePort,
		Status:    "UP",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Version:   version,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func (s *server) analyzeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req contracts.RequestGenericReadingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("[%s] Received analysis request for StudyID: %s", serviceName, req.StudyID)

	// Placeholder for the actual analysis logic.
	readingContent := fmt.Sprintf("Automated analysis for study [%s] on anatomical region [%s] completed.", req.StudyID, req.AnatomicalRegion)

	// Construct the response, strictly adhering to the domain model and clinical invariants.
	response := models.GenericReading{
		StudyID:          req.StudyID,
		AnatomicalRegion: req.AnatomicalRegion,
		Timestamp:        time.Now().UTC().Format(time.RFC3339),
		ReadingContent:   readingContent,
		Disclaimer: models.Disclaimer{
			ReadingType:               shared.ReadingTypeGenericAutomated,
			SignatureStatus:           shared.StatusWithoutMedicalSignature,
			MedicalReportStatus:       shared.MedicalReportNotIncluded,
			FullMandatoryLegalWarning: shared.MandatoryLegalDisclaimer,
		},
	}

	// Cumplir con REQ-AUD-002: Enviar un evento de auditoría.
	go func() {
		event := models.AuditEvent{
			EventTimestamp:   time.Now().UTC(),
			EventAction:      "GENERIC_READING_GENERATED",
			EventOutcome:     "SUCCESS",
			UserID:           serviceName, // The service is the actor.
			SourceIP:         r.RemoteAddr,
			StudyInstanceUID: req.StudyID,
		}
		s.auditClient.SendEvent(contracts.RecordAuditEventRequest{Event: event})
	}()

	// Enviar la lectura generada al servicio rx-result para su consolidación.
	go func() {
		s.resultClient.ConsolidateReading(response)
	}()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		log.Printf("[%s] ERROR: Failed to encode response: %v", serviceName, err)
		http.Error(w, "Failed to generate response", http.StatusInternalServerError)
	}
}