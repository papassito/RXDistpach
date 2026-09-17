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
	servicePort = 8083
	version     = "0.1.0-scaffold"
)

// server encapsula las dependencias de transporte y comunicación del microservicio.
type server struct {
	auditClient  transport.AuditClienter
	resultClient transport.ResultClienter
}

func main() {
	log.Printf("[%s] Starting service on port %d", serviceName, servicePort)

	// Inicialización del servidor con sus clientes de transporte
	srv := &server{
		auditClient:  transport.NewAuditClient(),
		resultClient: transport.NewResultClient(),
	}

	// Registro de rutas HTTP
	http.HandleFunc("/healthz", srv.healthCheckHandler)
	http.HandleFunc("/reader/analyze", srv.analyzeHandler)

	addr := fmt.Sprintf("127.0.0.1:%d", servicePort)
	log.Printf("[%s] Listening on %s", serviceName, addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("[%s] Server failed: %v", serviceName, err)
	}
}

// healthCheckHandler responde con el estado de salud del microservicio.
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

// analyzeHandler redirige las peticiones entrantes a la lógica de análisis.
func (s *server) analyzeHandler(w http.ResponseWriter, r *http.Request) {
	s.handleAnalyze(w, r)
}

// handleAnalyze procesa la lectura radiológica, aplica las invariantes clínicas y responde en JSON.
func (s *server) handleAnalyze(w http.ResponseWriter, r *http.Request) {
	// 1. Validar método HTTP
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Decodificar la solicitud
	var req contracts.RequestGenericReadingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
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
			FullMandatoryLegalWarning: shared.MandatoryLegalDisclaimer, // ✅ Usando la constante autoritativa
		},
	}

	// 4. Auditoría asíncrona (EventTimestamp corregido a tipo time.Time)
	go func() {
		event := models.AuditEvent{
			EventTimestamp:   time.Now().UTC(), // ✅ Tipo time.Time correcto
			EventAction:      "GENERIC_READING_GENERATED",
			EventOutcome:     "SUCCESS",
			UserID:           serviceName,
			SourceIP:         r.RemoteAddr,
			StudyInstanceUID: req.StudyID,
		}
		
		// Enviar evento según la definición del contrato
		_ = s.auditClient.SendEvent(contracts.RecordAuditEventRequest{
			Event: event,
		})
	}()

	// 5. Consolidación de resultado asíncrona
	go func() {
		_ = s.resultClient.ConsolidateReading(response)
	}()

	// 6. Escribir respuesta JSON completa para evitar errores EOF
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(&response)
}