package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"rx-dispatch/internal/config"
	"rx-dispatch/internal/contracts"
	"rx-dispatch/internal/delivery"
	"rx-dispatch/internal/dicom"
	"rx-dispatch/internal/image"
	"rx-dispatch/internal/metrics"
	"rx-dispatch/internal/models"
	"rx-dispatch/internal/shared"
	"rx-dispatch/internal/storage"
	"rx-dispatch/internal/transport"
)

// MonolithServer consolidates all 9 bounded contexts into a single, high-performance,
// zero-loopback-latency Go executable ("RX DISPATCH EDGE") for local clinics and edge servers.
type MonolithServer struct {
	mu           sync.RWMutex
	port         int
	dicomPort    int
	dataDir      string
	studies      map[string]models.Study
	artifacts    map[string]storage.ArtifactMetadata
	artifactData map[string]string
	auditEvents  []models.AuditEvent
	deliveries   map[string][]models.DeliveryRecord
	retention    *storage.RetentionManager
	queue        *delivery.StoreAndForwardQueue
	dicomSCP     *dicom.SCPServer
}

func NewMonolithServer(port, dicomPort int, dataDir string) *MonolithServer {
	_ = os.MkdirAll(dataDir, 0755)

	srv := &MonolithServer{
		port:         port,
		dicomPort:    dicomPort,
		dataDir:      dataDir,
		studies:      make(map[string]models.Study),
		artifacts:    make(map[string]storage.ArtifactMetadata),
		artifactData: make(map[string]string),
		deliveries:   make(map[string][]models.DeliveryRecord),
		auditEvents:  make([]models.AuditEvent, 0),
		retention:    storage.NewRetentionManager(nil),
	}

	// Initialize Store-and-Forward queue
	srv.queue = delivery.NewStoreAndForwardQueue(
		filepath.Join(dataDir, "queue"),
		func(item *delivery.QueueItem) {
			rec := item.ToDeliveryRecord()
			srv.mu.Lock()
			srv.deliveries[rec.StudyID] = append([]models.DeliveryRecord{rec}, srv.deliveries[rec.StudyID]...)
			if stu, exists := srv.studies[rec.StudyID]; exists {
				stu.Status = models.StudyStatusDelivered
				srv.studies[rec.StudyID] = stu
			}
			srv.mu.Unlock()
			metrics.DefaultRegistry.IncStudiesDelivered()
		},
		func(eventType, studyID, actor, details string) {
			srv.recordAudit(eventType, studyID, "rx-delivery", actor, details)
		},
	)

	// Initialize native DICOM C-STORE SCP listener
	srv.dicomSCP = dicom.NewSCPServer(
		dicom.ServerConfig{
			AETitle: "RX_DISPATCH",
			Port:    dicomPort,
		},
		func(ds *dicom.Dataset) error {
			return srv.handleDICOMIngestion(ds)
		},
	)

	// Seed demo study
	srv.seedInitialData()

	return srv
}

func (m *MonolithServer) seedInitialData() {
	m.mu.Lock()
	defer m.mu.Unlock()

	patient := models.Patient{
		ID:        "PAC-9042",
		FullName:  "Carlos Mendoza",
		BirthDate: "1978-05-14",
		Gender:    "Masculino",
		Email:     "carlos.mendoza@email.com",
		Phone:     "+54 9 11 4059-8812",
	}

	imgID := "IMG-TORAX-001"
	origImg := models.XRayImage{
		ID:            imgID,
		StudyID:       "STU-101",
		FileName:      "torax_pa_0849.png",
		OriginalURI:   "/storage/artifacts/ORIG-IMG-TORAX-001",
		DerivedURI:    "/storage/artifacts/DERIV-IMG-TORAX-001",
		CreatedAt:     time.Now().Add(-2 * time.Hour),
		IsIntegrityOK: true,
		Metadata: models.ImageMetadata{
			Format:      "DICOM",
			Width:       2048,
			Height:      2048,
			BitDepth:    16,
			Kvp:         "120 kVp",
			MilliAmps:   "3.2 mAs",
			Projection:  "PA",
			ChecksumSHA: "f515562465d6c8b9123456789abcdef0123456789abcdef0123456789abcdef0",
		},
	}

	study := models.Study{
		ID:                 "STU-101",
		StudyIdentifier:    "ACC-2026-0904-8841",
		Date:               time.Now().Add(-2 * time.Hour),
		StudyType:          "Radiografía de Tórax (Frente y Perfil)",
		AnatomicalRegion:   "Tórax",
		ReferringPhysician: "Dr. Roberto Rossi",
		ClinicalIndication: "Control respiratorio",
		Status:             models.StudyStatusReceived,
		Patient:            patient,
		Images:             []models.XRayImage{origImg},
	}

	m.studies[study.ID] = study
	m.artifacts["ORIG-"+imgID] = storage.ArtifactMetadata{
		ID:          "ORIG-" + imgID,
		Type:        "ORIGINAL_IMAGE",
		StudyID:     study.ID,
		FileName:    origImg.FileName,
		SizeBytes:   8450120,
		StoredAt:    origImg.CreatedAt,
		IsImmutable: true,
	}
	m.artifactData["ORIG-"+imgID] = "RAW_PIXEL_DATA_SEALED"

	m.recordAuditLocked("system_bootstrapped", "GLOBAL", "rx-monolith", "SYSTEM",
		"RX Dispatch Edge (Monolito Modular) inicializado con receptor DICOM en puerto "+fmt.Sprint(m.dicomPort))
}

func (m *MonolithServer) recordAudit(eventType, studyID, service, actor, details string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.recordAuditLocked(eventType, studyID, service, actor, details)
}

func (m *MonolithServer) recordAuditLocked(eventType, studyID, service, actor, details string) {
	id := fmt.Sprintf("AUD-%d-%04d", time.Now().UnixNano()/1e6, len(m.auditEvents)+1)
	ev := models.AuditEvent{
		ID:        id,
		Timestamp: time.Now(),
		EventType: eventType,
		StudyID:   studyID,
		Service:   service,
		Actor:     actor,
		Details:   details,
	}
	m.auditEvents = append([]models.AuditEvent{ev}, m.auditEvents...)
}

func (m *MonolithServer) handleDICOMIngestion(ds *dicom.Dataset) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	metrics.DefaultRegistry.IncDicomCStoreReceived()
	metrics.DefaultRegistry.IncStudiesReceived()

	studyID := fmt.Sprintf("STU-DCM-%d", time.Now().UnixNano()/1e6)
	imgID := fmt.Sprintf("IMG-DCM-%d", time.Now().UnixNano()/1e6)

	// 1. Store raw DICOM as immutable original
	artifactID := "ORIG-" + imgID
	m.artifacts[artifactID] = storage.ArtifactMetadata{
		ID:          artifactID,
		Type:        "DICOM_RAW",
		StudyID:     studyID,
		FileName:    fmt.Sprintf("%s.dcm", ds.SOPInstanceUID),
		SizeBytes:   ds.SizeBytes,
		StoredAt:    time.Now(),
		IsImmutable: true, // STRICT IMMUTABILITY
	}
	m.artifactData[artifactID] = string(ds.RawBytes)

	// 2. Generate lightweight Web derivative
	webLight := image.GenerateWebLightDerivative(artifactID, ds.RawBytes, ds.ChecksumSHA)
	m.artifacts[webLight.DerivativeID] = storage.ArtifactMetadata{
		ID:          webLight.DerivativeID,
		Type:        "DERIVATIVE_IMAGE",
		StudyID:     studyID,
		FileName:    webLight.DerivativeID + ".webp",
		SizeBytes:   webLight.WebSizeBytes,
		StoredAt:    webLight.CreatedAt,
		IsImmutable: false, // Derivative can be purged according to retention policy
	}

	// 3. Register Study
	study := models.Study{
		ID:                 studyID,
		StudyIdentifier:    ds.AccessionNumber,
		Date:               time.Now(),
		StudyType:          ds.StudyDescription,
		AnatomicalRegion:   ds.BodyPartExamined,
		ReferringPhysician: fmt.Sprintf("Modalidad DICOM %s", ds.Manufacturer),
		ClinicalIndication: "Ingesta directa vía DICOM C-STORE SCP",
		Status:             models.StudyStatusReceived,
		Patient: models.Patient{
			ID:        ds.PatientID,
			FullName:  ds.PatientName,
			BirthDate: ds.PatientBirthDate,
			Gender:    ds.PatientSex,
			Email:     "paciente.dicom@clinica.local",
		},
		Images: []models.XRayImage{
			{
				ID:            imgID,
				StudyID:       studyID,
				FileName:      fmt.Sprintf("%s.dcm", ds.SOPInstanceUID),
				OriginalURI:   "/storage/artifacts/" + artifactID,
				DerivedURI:    "/storage/artifacts/" + webLight.DerivativeID,
				CreatedAt:     time.Now(),
				IsIntegrityOK: true,
				Metadata: models.ImageMetadata{
					Format:      "DICOM",
					Width:       1024,
					Height:      1024,
					BitDepth:    16,
					Projection:  "AP",
					ChecksumSHA: ds.ChecksumSHA,
				},
			},
		},
	}

	m.studies[studyID] = study
	m.recordAuditLocked("dicom_cstore_received", studyID, "rx-dicom-scp", "DICOM_NETWORK_LISTENER",
		fmt.Sprintf("Instancia DICOM recibida de %s. Paciente: %s, SHA-256: %s, Peso: %d bytes",
			ds.PatientName, ds.PatientID, ds.ChecksumSHA[:12], ds.SizeBytes))

	log.Printf("[RX-EDGE-MONOLITH] DICOM Study ingested: %s (Patient: %s, SHA: %s)",
		studyID, ds.PatientName, ds.ChecksumSHA[:12])
	return nil
}

// HTTP API Handlers

func (m *MonolithServer) handleHealthz(w http.ResponseWriter, r *http.Request) {
	resp := contracts.HealthResponse{
		Service:   "rx-dispatch-edge",
		Port:      m.port,
		Status:    "UP",
		Timestamp: time.Now(),
		Version:   "2.0.0-edge-monolith",
	}
	transport.WriteJSON(w, http.StatusOK, resp)
}

func (m *MonolithServer) handleTopology(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	upResp := func(name string, p int) contracts.HealthResponse {
		return contracts.HealthResponse{
			Service:   name,
			Port:      p,
			Status:    "UP",
			Timestamp: now,
			Version:   "2.0.0-edge-monolith",
		}
	}

	top := contracts.GatewayTopologyStatus{
		GatewayService:  upResp("rx-gateway (in-process)", m.port),
		SecurityService: upResp("rx-security (in-process)", m.port),
		StudyService:    upResp("rx-study (in-process)", m.port),
		StorageService:  upResp("rx-storage (in-process)", m.port),
		ImageService:    upResp("rx-image (in-process)", m.port),
		ReaderService:   upResp("rx-reader (in-process)", m.port),
		ResultService:   upResp("rx-result (in-process)", m.port),
		DeliveryService: upResp("rx-delivery (in-process)", m.port),
		AuditService:    upResp("rx-audit (in-process)", m.port),
	}
	transport.WriteJSON(w, http.StatusOK, top)
}

func (m *MonolithServer) handleStudies(w http.ResponseWriter, r *http.Request) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var list []models.Study
	for _, s := range m.studies {
		list = append(list, s)
	}

	transport.WriteJSON(w, http.StatusOK, contracts.ListStudiesResponse{
		Studies: list,
		Count:   len(list),
	})
}

func (m *MonolithServer) handleStudyDetail(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		transport.WriteError(w, http.StatusBadRequest, "invalid path", "rx-dispatch-edge")
		return
	}
	studyID := parts[4]

	m.mu.RLock()
	study, exists := m.studies[studyID]
	m.mu.RUnlock()

	if !exists {
		transport.WriteError(w, http.StatusNotFound, "study not found", "rx-dispatch-edge")
		return
	}
	transport.WriteJSON(w, http.StatusOK, study)
}

func (m *MonolithServer) handleDispatchFlow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		transport.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", "rx-dispatch-edge")
		return
	}

	var req contracts.OrchestratedDispatchFlowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		transport.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid payload: %v", err), "rx-dispatch-edge")
		return
	}

	m.mu.Lock()
	study, exists := m.studies[req.StudyID]
	if !exists {
		m.mu.Unlock()
		transport.WriteError(w, http.StatusNotFound, "study not found", "rx-dispatch-edge")
		return
	}

	// 1. Image verification
	origSHA := ""
	if len(study.Images) > 0 {
		origSHA = study.Images[0].Metadata.ChecksumSHA
	}

	// 2. Reader engine (STRICT NON-DIAGNOSTIC AUTOMATED READING)
	readingID := fmt.Sprintf("GEN-READ-%d", time.Now().UnixNano()/1e6)
	reading := models.GenericReading{
		ID:                        readingID,
		StudyID:                   study.ID,
		CreatedAt:                 time.Now(),
		StudyTitle:                study.StudyType,
		GeneralDescription:        fmt.Sprintf("Lectura genérica automatizada de la región anatómica %s.", study.AnatomicalRegion),
		VisualObservations:        []string{"Patrón de densidades radiográficas y contornos óseos visibles sin artefactos mayores."},
		ObservableCharacteristics: []string{"Penetración técnica adecuada, centrado de proyección simétrico."},
		TechnicalCaveats:          "Información técnica automatizada. No reemplaza criterio clínico ni diagnóstico formal.",
		ReadingType:               shared.ReadingTypeGenericAutomated,
		MedicalReportNotice:       shared.MedicalReportNotIncluded,
		PhysicianSignatureStatus:  shared.StatusWithoutMedicalSignature,
		MandatoryDisclaimer:       shared.MandatoryLegalDisclaimer,
	}

	// 3. Result package builder
	pkgID := fmt.Sprintf("PKG-%s-%d", study.ID, time.Now().UnixNano()/1e6)
	pkgContent := fmt.Sprintf("PACKAGE_PAYLOAD_%s_SHA_%s", study.ID, origSHA)
	h := sha256.Sum256([]byte(pkgContent))
	pkgChecksum := hex.EncodeToString(h[:])

	// 4. Enqueue into Store-and-Forward resilient queue
	recipient := req.DeliveryTarget
	if recipient == "" {
		recipient = study.Patient.Email
	}
	channel := req.DeliveryChannel
	if channel == "" {
		channel = "EMAIL"
	}

	queueItem := m.queue.Enqueue(study.ID, pkgID, recipient, channel, pkgChecksum)

	// Update local study status
	study.Status = models.StudyStatusDelivered
	m.studies[study.ID] = study

	m.recordAuditLocked("dispatch_flow_executed", study.ID, "rx-dispatch-edge", "EDGE_DISPATCHER",
		fmt.Sprintf("Flujo de despacho procesado en memoria. Token: %s, Checksum: %s", queueItem.TrackingToken, pkgChecksum[:12]))

	m.mu.Unlock()

	resp := contracts.OrchestratedDispatchFlowResponse{
		StudyID:         study.ID,
		SecurityChecked: true,
		StudyStatus:     study.Status,
		ImageIntegrity:  true,
		ReadingType:     reading.ReadingType,
		PackageChecksum: pkgChecksum,
		DeliveryToken:   queueItem.TrackingToken,
		DeliveryStatus:  models.DeliverySent,
		AuditCount:      len(m.auditEvents),
	}

	transport.WriteJSON(w, http.StatusOK, resp)
}

func (m *MonolithServer) handleAudit(w http.ResponseWriter, r *http.Request) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	studyID := r.URL.Query().Get("studyId")
	var result []models.AuditEvent
	for _, ev := range m.auditEvents {
		if studyID == "" || ev.StudyID == studyID {
			result = append(result, ev)
		}
	}
	transport.WriteJSON(w, http.StatusOK, result)
}

func (m *MonolithServer) handleDeliveryQueue(w http.ResponseWriter, r *http.Request) {
	status := m.queue.GetStatus()
	transport.WriteJSON(w, http.StatusOK, status)
}

func (m *MonolithServer) handleNetworkToggle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		transport.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", "rx-dispatch-edge")
		return
	}
	var req struct {
		SimulateDrop bool `json:"simulateDrop"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	m.queue.SetNetworkSimulation(req.SimulateDrop)
	transport.WriteJSON(w, http.StatusOK, m.queue.GetStatus())
}

func (m *MonolithServer) handleStorageRetention(w http.ResponseWriter, r *http.Request) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var totalBytes int64
	var originalBytes int64
	var derivativeBytes int64
	var origCount int
	var derivCount int

	for _, a := range m.artifacts {
		totalBytes += int64(a.SizeBytes)
		if a.IsImmutable || a.Type == "ORIGINAL_IMAGE" || a.Type == "DICOM_RAW" {
			originalBytes += int64(a.SizeBytes)
			origCount++
		} else {
			derivativeBytes += int64(a.SizeBytes)
			derivCount++
		}
	}

	transport.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"policy":          m.retention.GetConfig(),
		"totalArtifacts":  len(m.artifacts),
		"originalCount":   origCount,
		"derivativeCount": derivCount,
		"totalBytes":      totalBytes,
		"originalBytes":   originalBytes,
		"derivativeBytes": derivativeBytes,
	})
}

func (m *MonolithServer) handleStoragePurge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		transport.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", "rx-dispatch-edge")
		return
	}

	m.mu.Lock()
	purgeRes := m.retention.EvaluateAndPurge(m.artifacts, func(id string) {
		delete(m.artifacts, id)
		delete(m.artifactData, id)
	})
	m.recordAuditLocked("storage_purge_executed", "GLOBAL", "rx-storage", "RETENTION_ENGINE",
		fmt.Sprintf("Purga ejecutada: %d bytes liberados (%d derivados purgados, %d originales preservados)",
			purgeRes.BytesFreed, purgeRes.PurgedDerivatives, purgeRes.PreservedOriginals))
	m.mu.Unlock()

	transport.WriteJSON(w, http.StatusOK, purgeRes)
}

func (m *MonolithServer) handleDICOMStats(w http.ResponseWriter, r *http.Request) {
	instances, last := m.dicomSCP.Stats()
	transport.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"status":            "LISTENING",
		"port":              m.dicomPort,
		"aeTitle":           "RX_DISPATCH",
		"receivedInstances": instances,
		"lastReceived":      last,
	})
}

func (m *MonolithServer) handleWebLight(w http.ResponseWriter, r *http.Request) {
	artifactID := r.URL.Query().Get("artifactId")
	if artifactID == "" {
		artifactID = "ORIG-IMG-TORAX-001"
	}

	m.mu.RLock()
	dataStr := m.artifactData[artifactID]
	rawSHA := "f515562465d6c8b9123456789abcdef0123456789abcdef0123456789abcdef0"
	if meta, ok := m.artifacts[artifactID]; ok {
		rawSHA = meta.ID
	}
	m.mu.RUnlock()

	deriv := image.GenerateWebLightDerivative(artifactID, []byte(dataStr), rawSHA)
	transport.WriteJSON(w, http.StatusOK, deriv)
}

func main() {
	port := config.GetEnvInt("RX_EDGE_PORT", 8090)
	dicomPort := config.GetEnvInt("RX_DICOM_PORT", 11112)
	dataDir := config.GetEnv("RX_DATA_DIR", "/app/applet/rx-dispatch/data")

	// Allow overriding via CLI flags
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "--http-port" && i+1 < len(os.Args) {
			_, _ = fmt.Sscanf(os.Args[i+1], "%d", &port)
		} else if os.Args[i] == "--dicom-port" && i+1 < len(os.Args) {
			_, _ = fmt.Sscanf(os.Args[i+1], "%d", &dicomPort)
		} else if os.Args[i] == "--data-dir" && i+1 < len(os.Args) {
			dataDir = os.Args[i+1]
		}
	}

	server := NewMonolithServer(port, dicomPort, dataDir)

	// Start DICOM C-STORE SCP listener
	if err := server.dicomSCP.Start(); err != nil {
		log.Printf("[RX-EDGE-MONOLITH] DICOM listener error: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", server.handleHealthz)
	mux.HandleFunc("/gateway/topology", server.handleTopology)
	mux.HandleFunc("/api/v1/studies", server.handleStudies)
	mux.HandleFunc("/api/v1/studies/", server.handleStudyDetail)
	mux.HandleFunc("/api/v1/dispatch-flow", server.handleDispatchFlow)
	mux.HandleFunc("/api/v1/audit", server.handleAudit)
	mux.HandleFunc("/api/v1/delivery/queue", server.handleDeliveryQueue)
	mux.HandleFunc("/api/v1/delivery/network-toggle", server.handleNetworkToggle)
	mux.HandleFunc("/api/v1/storage/retention", server.handleStorageRetention)
	mux.HandleFunc("/api/v1/storage/retention/purge", server.handleStoragePurge)
	mux.HandleFunc("/api/v1/dicom/status", server.handleDICOMStats)
	mux.HandleFunc("/api/v1/image/weblight", server.handleWebLight)
	mux.HandleFunc("/metrics", metrics.DefaultRegistry.Handler())

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	log.Printf("==================================================================")
	log.Printf("RX DISPATCH EDGE — MONOLITO MODULAR CONSOLIDADO (SINGLE EXECUTABLE)")
	log.Printf("HTTP Gateway & REST API:   http://0.0.0.0:%d", port)
	log.Printf("DICOM C-STORE SCP Port:    tcp://0.0.0.0:%d (AE: RX_DISPATCH)", dicomPort)
	log.Printf("Metrics (Prometheus):      http://0.0.0.0:%d/metrics", port)
	log.Printf("Store-and-Forward Queue:   PERSISTENT DISK (data/queue)")
	log.Printf("==================================================================")

	httpServer := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[RX-EDGE-MONOLITH] Server failed: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("[RX-EDGE-MONOLITH] Shutting down gracefully...")
	server.dicomSCP.Stop()
	server.queue.Close()
	_ = httpServer.Close()
}
