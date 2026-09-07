package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"rx-dispatch/internal/config"
	"rx-dispatch/internal/contracts"
	"rx-dispatch/internal/models"
	"rx-dispatch/internal/shared"
	"rx-dispatch/internal/transport"

	"github.com/google/uuid" // Import for robust unique identifiers
)

type ReaderServer struct {
	port     int
	topology config.TopologyConfig
}

func NewReaderServer(port int, topology config.TopologyConfig) *ReaderServer {
	return &ReaderServer{
		port:     port,
		topology: topology,
	}
}

func (s *ReaderServer) handleHealthz(w http.ResponseWriter, r *http.Request) {
	resp := contracts.HealthResponse{
		Service:   "rx-reader",
		Port:      s.port,
		Status:    "UP",
		Timestamp: time.Now(),
		Version:   "2.0.0-decoupled",
	}
	transport.WriteJSON(w, http.StatusOK, resp)
}

// generateObservationsByRegion encapsulates the business logic for creating generic text based on anatomy.
// This makes the main handler cleaner and the rules easier to manage.
func generateObservationsByRegion(anatomicalRegion string) (visualObservations, observableCharacteristics []string) {
	regionLower := strings.ToLower(anatomicalRegion)

	switch {
	case strings.Contains(regionLower, "tórax") || strings.Contains(regionLower, "torax") || strings.Contains(regionLower, "chest"):
		visualObservations = []string{
			"Campos pleuropulmonares con aireación y trama vascular de distribución simétrica observable.",
			"Silueta cardiomediastínica en límites visuales de amplitud transversal.",
			"Ángulos costofrénicos y cardiofrénicos visualizados libres de velamiento.",
			"Estructura ósea de la jaula torácica sin solución de continuidad patente en esta proyección.",
		}
		observableCharacteristics = []string{
			"Técnica de exposición: Adecuada penetración con visualización de cuerpos vertebrales dorsales.",
			"Inspiración radiográfica: Recuento aproximado de 9 a 10 arcos costales posteriores.",
			"Centrado clavicular simétrico respecto a apófisis espinosas dorsales.",
		}
	case strings.Contains(regionLower, "muñeca") || strings.Contains(regionLower, "mano") || strings.Contains(regionLower, "wrist"):
		visualObservations = []string{
			"Cortesía de planos corticales óseos en radio distal y estiloides cubital.",
			"Alineación visual de huesos del carpo con mantenimiento de arcos de Gilula en esta proyección.",
			"Líneas articulares radiocarpiana e intercarpiana con preservación de amplitud radiográfica.",
			"Espacio de partes blandas adyacente sin aumento difuso de densidad.",
		}
		observableCharacteristics = []string{
			"Densidad ósea trabecular homogénea sin imágenes líticas ni blásticas focales evidentes.",
			"Ausencia de cuerpos extraños radiopacos de densidad cálcica o metálica anormal.",
		}
	default:
		visualObservations = []string{
			fmt.Sprintf("Estructuras anatómicas reconocibles en proyección estándar de %s.", anatomicalRegion),
			"Diferenciación de contraste de densidades físicas: aire, grasa, partes blandas y hueso.",
			"Contornos óseos y axiales observables según técnica radiológica efectuada.",
		}
		observableCharacteristics = []string{
			fmt.Sprintf("Región anatómica observada: %s.", anatomicalRegion),
			"Ausencia de artefactos de movimiento significativos durante el disparo.",
		}
	}

	return visualObservations, observableCharacteristics
}

func (s *ReaderServer) handleAnalyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		transport.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", "rx-reader")
		return
	}

	var req contracts.RequestGenericReadingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		transport.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid payload: %v", err), "rx-reader")
		return
	}

	visualObservations, observableCharacteristics := generateObservationsByRegion(req.AnatomicalRegion) // Use the refactored function

	// Use UUID for a robust, collision-free unique identifier.
	readingID := fmt.Sprintf("READ-%s", uuid.NewString()) // Switched from time-based to UUID
	reading := models.GenericReading{
		ID:                        readingID,
		StudyID:                   req.StudyID,
		CreatedAt:                 time.Now(),
		StudyTitle:                req.StudyType,
		GeneralDescription:        fmt.Sprintf("Se observa una imagen radiográfica correspondiente a la región anatómica indicada en el estudio (%s).", req.AnatomicalRegion),
		VisualObservations:        visualObservations,
		ObservableCharacteristics: observableCharacteristics,
		TechnicalCaveats:          "Esta información corresponde exclusivamente a una lectura genérica/automatizada de la imagen y no constituye un diagnóstico médico ni sustituye un informe radiológico emitido y firmado por un médico responsable.",
		// All relevant fields are now part of the model itself for consistency
		ReadingType:              shared.ReadingTypeGenericAutomated,
		MedicalReportNotice:      shared.MedicalReportNotIncluded,
		PhysicianSignatureStatus: shared.StatusWithoutMedicalSignature,
		MandatoryDisclaimer:      shared.MandatoryLegalDisclaimer,
	}

	// Audit record to RX Audit
	auditReq := contracts.RecordAuditEventRequest{
		EventType: "generic_reading_generated",
		StudyID:   req.StudyID,
		Service:   "rx-reader",
		Actor:     "SYSTEM_GENERIC_READER",
		Details:   "Lectura genérica de imagen generada. Clasificación: SIN FIRMA MÉDICA. Informe oficial: NO INCLUIDO.",
	}
	if err := transport.PostJSON(fmt.Sprintf("%s/audit/events", s.topology.AuditURL), auditReq, nil); err != nil {
		log.Printf("[RX-READER] WARNING: Failed to record audit event for study %s: %v", req.StudyID, err)
	}

	// The response contract is now simplified, containing only the complete Reading object.
	resp := contracts.RequestGenericReadingResponse{
		Reading:                  reading,
		IsNonDiagnosticCertified: true,
	}

	log.Printf("[RX-READER] Generic reading %s created for study %s (%s)", readingID, req.StudyID, reading.ReadingType)
	transport.WriteJSON(w, http.StatusOK, resp)
}

func main() {
	port := config.GetEnvInt("RX_READER_PORT", 8085)
	topology := config.LoadTopologyConfig()
	server := NewReaderServer(port, topology)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", server.handleHealthz)
	mux.HandleFunc("/reader/analyze", server.handleAnalyze)

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	log.Printf("[RX-READER] Starting on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("[RX-READER] Server failed: %v", err)
	}
}
