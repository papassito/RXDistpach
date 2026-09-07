package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"rx-dispatch/internal/config"
	"rx-dispatch/internal/contracts"
	"rx-dispatch/internal/models"
	"rx-dispatch/internal/transport"
)

type StudyServer struct {
	mu       sync.RWMutex
	studies  map[string]models.Study
	port     int
	topology config.TopologyConfig
}

func NewStudyServer(port int, topology config.TopologyConfig) *StudyServer {
	s := &StudyServer{
		studies:  make(map[string]models.Study),
		port:     port,
		topology: topology,
	}
	s.seedDefaultStudies()
	return s
}

func (s *StudyServer) seedDefaultStudies() {
	sampleStudy := models.Study{
		ID:                 "STU-101",
		StudyIdentifier:    "RX-2026-0849",
		Date:               time.Now(),
		StudyType:          "Radiografía de Tórax PA y Lateral",
		AnatomicalRegion:   "Tórax",
		ReferringPhysician: "Dra. Elena Valenzuela (Medicina Interna)",
		ClinicalIndication: "Evaluación respiratoria de rutina, control de patrón broncopulmonar.",
		Status:             models.StudyStatusReceived,
		Patient: models.Patient{
			ID:        "PAC-7741",
			FullName:  "Carlos Mendoza Rivas",
			BirthDate: "1982-04-15",
			Gender:    "Masculino",
			Email:     "carlos.mendoza@ejemplo.com",
			Phone:     "+34 612 345 678",
		},
		Images: []models.XRayImage{
			{
				ID:            "IMG-101-1",
				StudyID:       "STU-101",
				FileName:      "torax_pa_0849.png",
				OriginalURI:   "/storage/artifacts/ORIG-IMG-101-1",
				DerivedURI:    "/storage/artifacts/DERIV-IMG-101-1",
				CreatedAt:     time.Now(),
				IsIntegrityOK: true,
				Metadata: models.ImageMetadata{
					Format:      "PNG/DICOM",
					Width:       2048,
					Height:      2048,
					BitDepth:    16,
					Kvp:         "120 kVp",
					MilliAmps:   "3.2 mAs",
					Projection:  "PA",
					ChecksumSHA: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				},
			},
		},
	}
	s.studies[sampleStudy.ID] = sampleStudy
}

func (s *StudyServer) handleHealthz(w http.ResponseWriter, r *http.Request) {
	resp := contracts.HealthResponse{
		Service:   "rx-study",
		Port:      s.port,
		Status:    "UP",
		Timestamp: time.Now(),
		Version:   "2.0.0-decoupled",
	}
	transport.WriteJSON(w, http.StatusOK, resp)
}

func (s *StudyServer) handleStudies(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	// GET /studies or POST /studies
	if len(parts) == 1 && parts[0] == "studies" {
		if r.Method == http.MethodGet {
			s.mu.RLock()
			list := make([]models.Study, 0, len(s.studies))
			for _, st := range s.studies {
				list = append(list, st)
			}
			s.mu.RUnlock()
			transport.WriteJSON(w, http.StatusOK, contracts.ListStudiesResponse{
				Studies: list,
				Count:   len(list),
			})
			return
		}

		if r.Method == http.MethodPost {
			var req contracts.IngestStudyRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				transport.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid payload: %v", err), "rx-study")
				return
			}

			newID := fmt.Sprintf("STU-%d", time.Now().UnixNano()/1e6)
			identifier := req.StudyIdentifier
			if identifier == "" {
				identifier = fmt.Sprintf("RX-%d-%04d", time.Now().Year(), time.Now().Unix()%10000)
			}

			study := models.Study{
				ID:                 newID,
				StudyIdentifier:    identifier,
				Date:               time.Now(),
				StudyType:          req.StudyType,
				AnatomicalRegion:   req.AnatomicalRegion,
				ReferringPhysician: req.ReferringPhysician,
				ClinicalIndication: req.ClinicalIndication,
				Status:             models.StudyStatusReceived,
				Patient:            req.Patient,
				Images:             make([]models.XRayImage, 0),
			}

			s.mu.Lock()
			s.studies[newID] = study
			s.mu.Unlock()

			// Emit audit event
			auditReq := contracts.RecordAuditEventRequest{
				EventType: "study_received",
				StudyID:   study.ID,
				Service:   "rx-study",
				Actor:     "STUDY_INGEST",
				Details:   fmt.Sprintf("Estudio %s recibido para paciente %s.", study.StudyIdentifier, study.Patient.FullName),
			}
			_ = transport.PostJSON(fmt.Sprintf("%s/audit/events", s.topology.AuditURL), auditReq, nil)

			log.Printf("[RX-STUDY] Ingested study %s (%s)", study.ID, study.StudyIdentifier)
			transport.WriteJSON(w, http.StatusCreated, contracts.IngestStudyResponse{Study: study})
			return
		}

		transport.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", "rx-study")
		return
	}

	// Routes with study ID: /studies/{id}, /studies/{id}/status, /studies/{id}/images
	studyID := parts[1]
	s.mu.RLock()
	study, exists := s.studies[studyID]
	s.mu.RUnlock()

	if !exists {
		transport.WriteError(w, http.StatusNotFound, "study not found", "rx-study")
		return
	}

	if len(parts) == 2 {
		if r.Method == http.MethodGet {
			transport.WriteJSON(w, http.StatusOK, study)
			return
		}
	} else if len(parts) == 3 {
		action := parts[2]
		if action == "status" && (r.Method == http.MethodPatch || r.Method == http.MethodPost) {
			var updateReq contracts.UpdateStudyStatusRequest
			if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
				transport.WriteError(w, http.StatusBadRequest, "invalid status payload", "rx-study")
				return
			}
			s.mu.Lock()
			study.Status = updateReq.Status
			s.studies[studyID] = study
			s.mu.Unlock()

			// Audit transition
			auditReq := contracts.RecordAuditEventRequest{
				EventType: "study_status_updated",
				StudyID:   study.ID,
				Service:   "rx-study",
				Actor:     "STUDY_LIFECYCLE",
				Details:   fmt.Sprintf("Estado actualizado a %s. Motivo: %s", updateReq.Status, updateReq.Reason),
			}
			_ = transport.PostJSON(fmt.Sprintf("%s/audit/events", s.topology.AuditURL), auditReq, nil)

			transport.WriteJSON(w, http.StatusOK, study)
			return
		}

		if action == "images" && r.Method == http.MethodPost {
			var attachReq contracts.AttachImageToStudyRequest
			if err := json.NewDecoder(r.Body).Decode(&attachReq); err != nil {
				transport.WriteError(w, http.StatusBadRequest, "invalid image payload", "rx-study")
				return
			}
			s.mu.Lock()
			study.Images = append(study.Images, attachReq.Image)
			s.studies[studyID] = study
			s.mu.Unlock()

			transport.WriteJSON(w, http.StatusOK, study)
			return
		}
	}

	transport.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", "rx-study")
}

func main() {
	port := config.GetEnvInt("RX_STUDY_PORT", 8082)
	topology := config.LoadTopologyConfig()
	server := NewStudyServer(port, topology)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", server.handleHealthz)
	mux.HandleFunc("/studies", server.handleStudies)
	mux.HandleFunc("/studies/", server.handleStudies)

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	log.Printf("[RX-STUDY] Starting on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("[RX-STUDY] Server failed: %v", err)
	}
}
