package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"rx-dispatch/internal/config"
	"rx-dispatch/internal/contracts"
	"rx-dispatch/internal/models"
	"rx-dispatch/internal/shared"
	"rx-dispatch/internal/transport"
)

type ResultServer struct {
	port     int
	topology config.TopologyConfig
}

func NewResultServer(port int, topology config.TopologyConfig) *ResultServer {
	return &ResultServer{
		port:     port,
		topology: topology,
	}
}

func (s *ResultServer) handleHealthz(w http.ResponseWriter, r *http.Request) {
	resp := contracts.HealthResponse{
		Service:   "rx-result",
		Port:      s.port,
		Status:    "UP",
		Timestamp: time.Now(),
		Version:   "2.0.0-decoupled",
	}
	transport.WriteJSON(w, http.StatusOK, resp)
}

func (s *ResultServer) handleBuild(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		transport.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", "rx-result")
		return
	}

	var req contracts.BuildResultRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		transport.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid payload: %v", err), "rx-result")
		return
	}

	packageID := fmt.Sprintf("PKG-%s-%d", req.StudyID, time.Now().UnixNano()/1e6)

	// Build file manifest for package
	fileManifest := []string{
		"resultado_lectura_generica.pdf",
		"lectura_descriptiva.txt",
		"aviso_normativo_disclaimer.txt",
	}
	for i, img := range req.Study.Images {
		fileManifest = append(fileManifest, fmt.Sprintf("images/rx_%02d_%s", i+1, img.FileName))
	}

	// Compute composite package checksum
	pkgContent := fmt.Sprintf("%s|%s|%s|%d", packageID, req.StudyID, req.Reading.ID, len(fileManifest))
	checksumSHA := shared.CalculateStringSHA256(pkgContent)

	resultPackage := models.ResultPackage{
		ID:                       packageID,
		StudyID:                  req.StudyID,
		GeneratedAt:              time.Now(),
		Study:                    req.Study,
		Reading:                  req.Reading,
		ResultType:               "LECTURA GENÉRICA AUTOMATIZADA",
		StatusNotice:             "SIN FIRMA MÉDICA",
		OfficialReportNotice:     "INFORME MÉDICO OFICIAL: NO INCLUIDO",
		PhysicianSignatureNotice: "FIRMA DEL MÉDICO RESPONSABLE: DEBE SOLICITARSE",
		MandatoryLegend:          shared.MandatoryLegalDisclaimer,
		RequestInstructions:      "Para solicitar el informe radiológico oficial firmado por el médico especialista responsable, contacte al Servicio de Radiodiagnóstico indicando el Identificador del Estudio.",
		PackageFileNames:         fileManifest,
		ChecksumSHA:              checksumSHA,
	}

	// Store package manifest into RX Storage
	storageReq := contracts.StoreArtifactRequest{
		ArtifactID:   packageID,
		ArtifactType: "RESULT_PACKAGE",
		StudyID:      req.StudyID,
		FileName:     fmt.Sprintf("package_%s.json", packageID),
		Data:         checksumSHA,
	}
	if err := transport.PostJSON(fmt.Sprintf("%s/storage/artifacts", s.topology.StorageURL), storageReq, nil); err != nil {
		log.Printf("[RX-RESULT] WARNING: Failed to store result package artifact in rx-storage for study %s: %v", req.StudyID, err)
	}

	// Record in RX Audit
	auditReq := contracts.RecordAuditEventRequest{
		EventType: "result_package_built",
		StudyID:   req.StudyID,
		Service:   "rx-result",
		Actor:     "SYSTEM_RESULT_BUILDER",
		Details:   fmt.Sprintf("Paquete %s ensamblado con %d archivos y leyenda obligatoria.", packageID, len(fileManifest)),
	}
	if err := transport.PostJSON(fmt.Sprintf("%s/audit/events", s.topology.AuditURL), auditReq, nil); err != nil {
		log.Printf("[RX-RESULT] WARNING: Failed to record audit event for study %s: %v", req.StudyID, err)
	}

	resp := contracts.BuildResultResponse{
		ResultPackage: resultPackage,
		PackageID:     packageID,
		ChecksumSHA:   checksumSHA,
		FileManifest:  fileManifest,
	}

	log.Printf("[RX-RESULT] Result package %s built for study %s", packageID, req.StudyID)
	transport.WriteJSON(w, http.StatusOK, resp)
}

func main() {
	port := config.GetEnvInt("RX_RESULT_PORT", 8086)
	topology := config.LoadTopologyConfig()
	server := NewResultServer(port, topology)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", server.handleHealthz)
	mux.HandleFunc("/result/build", server.handleBuild)

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	log.Printf("[RX-RESULT] Starting on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("[RX-RESULT] Server failed: %v", err)
	}
}
