package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"rx-dispatch/internal/contracts"
	"rx-dispatch/internal/models"
	"rx-dispatch/internal/shared"
	"rx-dispatch/internal/transport"
)

const (
	serviceName = "rx-reader"
	defaultPort = "8085" // Corrected port
	version     = "0.1.0-scaffold"
)

// server encapsula las dependencias de transporte y comunicaciÃ³n del microservicio.
type server struct {
	auditClient  transport.AuditClienter
	resultClient transport.ResultClienter
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	host := os.Getenv("HOST")
	log.Printf("[%s] Starting service on port %s", serviceName, port)

	// InicializaciÃ³n del servidor con sus clientes de transporte
	srv := &server{
		auditClient:  transport.NewAuditClient(),
		resultClient: transport.NewResultClient(),
	}

	// Registro de rutas HTTP
	http.HandleFunc("/healthz", srv.healthCheckHandler)
	http.HandleFunc("/reader/analyze", srv.analyzeHandler)

	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("[%s] Listening on %s", serviceName, addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("[%s] Server failed: %v", serviceName, err)
	}
}

// healthCheckHandler responde con el estado de salud del microservicio.
func (s *server) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

// analyzeHandler procesa la lectura radiolÃ³gica, aplica las invariantes clÃ­nicas y responde en JSON.
func (s *server) analyzeHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Validar mÃ©todo HTTP
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// Inyectar protecciÃ³n DoS
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB limit

	// 2. Decodificar la solicitud
	var req contracts.RequestGenericReadingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid or malformed JSON payload", http.StatusBadRequest)
		return
	}

	readingContent := fmt.Sprintf("Automated analysis for study [%s] on anatomical region [%s] completed.", req.StudyID, req.AnatomicalRegion)

	// 3. Construir respuesta con invariantes de dominio
	response := models.GenericReading{
		StudyID:          req.StudyID,
		AnatomicalRegion: req.AnatomicalRegion,
		Timestamp:        time.Now().UTC().Format(time.RFC3339),
		ReadingContent:   readingContent,
		Disclaimer: models.Disclaimer{
			ReadingType:               shared.ReadingTypeGenericAutomated,
			SignatureStatus:           shared.StatusWithoutMedicalSignature,
			MedicalReportStatus:       shared.MedicalReportNotIncluded,
			FullMandatoryLegalWarning: shared.MandatoryLegalDisclaimer, // âœ… Usando la constante autoritativa
		},
	}

	// 4. AuditorÃ­a (sÃ­ncrona, con manejo de errores)
	event := models.AuditEventEPHI{
		EventTimestamp:   time.Now().UTC(),
		EventAction:      "GENERIC_READING_GENERATED",
		EventOutcome:     "SUCCESS",
		UserID:           serviceName,
		SourceIP:         r.RemoteAddr,
		StudyInstanceUID: req.StudyID,
	}
	if err := s.auditClient.SendEvent(contracts.RecordAuditEventRequest{Event: event}); err != nil {
		// No fallar la solicitud principal, pero registrar el fallo de auditorÃ­a.
		log.Printf("[%s] ERROR: Failed to send audit event for study %s: %v", serviceName, req.StudyID, err)
	}

	// 5. ConsolidaciÃ³n de resultado (sÃ­ncrona, con manejo de errores)
	if err := s.resultClient.ConsolidateReading(response); err != nil {
		log.Printf("[%s] ERROR: Failed to consolidate reading for study %s: %v", serviceName, req.StudyID, err)
		http.Error(w, "Failed to process reading result", http.StatusInternalServerError)
		return
	}

	// 6. Escribir respuesta JSON completa para evitar errores EOF
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		// Si la escritura de la respuesta falla, solo podemos registrarla.
		log.Printf("[%s] ERROR: Failed to write response for study %s: %v", serviceName, req.StudyID, err)
	}
}
