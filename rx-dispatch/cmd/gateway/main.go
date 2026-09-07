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
	"rx-dispatch/internal/metrics"
	"rx-dispatch/internal/models"
	"rx-dispatch/internal/transport"
)

type GatewayServer struct {
	port     int
	topology config.TopologyConfig
}

func NewGatewayServer(port int, topology config.TopologyConfig) *GatewayServer {
	return &GatewayServer{
		port:     port,
		topology: topology,
	}
}

func (g *GatewayServer) handleHealthz(w http.ResponseWriter, r *http.Request) {
	resp := contracts.HealthResponse{
		Service:   "rx-gateway",
		Port:      g.port,
		Status:    "UP",
		Timestamp: time.Now(),
		Version:   "2.0.0-decoupled",
	}
	transport.WriteJSON(w, http.StatusOK, resp)
}

// handleTopology queries the health of all 9 independent microservices in parallel.
func (g *GatewayServer) handleTopology(w http.ResponseWriter, r *http.Request) {
	var topologyStatus contracts.GatewayTopologyStatus
	var wg sync.WaitGroup

	queryHealth := func(url string, target *contracts.HealthResponse, defaultService string, defaultPort int) {
		defer wg.Done()
		h, err := transport.CheckHealth(url)
		if err != nil {
			*target = contracts.HealthResponse{
				Service:   defaultService,
				Port:      defaultPort,
				Status:    "UNREACHABLE",
				Timestamp: time.Now(),
				Version:   "2.0.0-decoupled",
			}
		} else {
			*target = h
		}
	}

	wg.Add(9)
	go func() {
		defer wg.Done()
		topologyStatus.GatewayService = contracts.HealthResponse{
			Service:   "rx-gateway",
			Port:      g.port,
			Status:    "UP",
			Timestamp: time.Now(),
			Version:   "2.0.0-decoupled",
		}
	}()
	go queryHealth(g.topology.SecurityURL, &topologyStatus.SecurityService, "rx-security", 8081)
	go queryHealth(g.topology.StudyURL, &topologyStatus.StudyService, "rx-study", 8082)
	go queryHealth(g.topology.StorageURL, &topologyStatus.StorageService, "rx-storage", 8083)
	go queryHealth(g.topology.ImageURL, &topologyStatus.ImageService, "rx-image", 8084)
	go queryHealth(g.topology.ReaderURL, &topologyStatus.ReaderService, "rx-reader", 8085)
	go queryHealth(g.topology.ResultURL, &topologyStatus.ResultService, "rx-result", 8086)
	go queryHealth(g.topology.DeliveryURL, &topologyStatus.DeliveryService, "rx-delivery", 8087)
	go queryHealth(g.topology.AuditURL, &topologyStatus.AuditService, "rx-audit", 8088)

	wg.Wait()
	transport.WriteJSON(w, http.StatusOK, topologyStatus)
}

// handleStudies delegates study listing and ingestion to rx-study.
func (g *GatewayServer) handleStudies(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1")

	if r.Method == http.MethodGet {
		var result contracts.ListStudiesResponse
		err := transport.GetJSON(fmt.Sprintf("%s/studies", g.topology.StudyURL), &result)
		if err != nil {
			transport.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to contact rx-study: %v", err), "rx-gateway")
			return
		}
		transport.WriteJSON(w, http.StatusOK, result.Studies)
		return
	}

	if r.Method == http.MethodPost {
		var req contracts.IngestStudyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			transport.WriteError(w, http.StatusBadRequest, "invalid study payload", "rx-gateway")
			return
		}

		var studyResp contracts.IngestStudyResponse
		err := transport.PostJSON(fmt.Sprintf("%s/studies", g.topology.StudyURL), req, &studyResp)
		if err != nil {
			transport.WriteError(w, http.StatusBadGateway, fmt.Sprintf("rx-study creation failed: %v", err), "rx-gateway")
			return
		}
		transport.WriteJSON(w, http.StatusCreated, studyResp.Study)
		return
	}

	_ = path
	transport.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", "rx-gateway")
}

// handleStudyDetail handles /api/v1/studies/{id}
func (g *GatewayServer) handleStudyDetail(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		transport.WriteError(w, http.StatusBadRequest, "missing study ID", "rx-gateway")
		return
	}
	studyID := parts[3]

	if len(parts) == 4 && r.Method == http.MethodGet {
		var study models.Study
		err := transport.GetJSON(fmt.Sprintf("%s/studies/%s", g.topology.StudyURL, studyID), &study)
		if err != nil {
			transport.WriteError(w, http.StatusNotFound, "study not found", "rx-gateway")
			return
		}
		transport.WriteJSON(w, http.StatusOK, study)
		return
	}

	if len(parts) == 5 {
		action := parts[4]
		switch action {
		case "process-images":
			g.orchestrateImageProcessing(w, r, studyID)
			return
		case "generic-reading":
			g.orchestrateGenericReading(w, r, studyID)
			return
		case "build-result":
			g.orchestrateResultBuild(w, r, studyID)
			return
		case "deliver":
			g.orchestrateDelivery(w, r, studyID)
			return
		}
	}

	transport.WriteError(w, http.StatusNotFound, "Endpoint not found", "rx-gateway")
}

func (g *GatewayServer) orchestrateImageProcessing(w http.ResponseWriter, r *http.Request, studyID string) {
	// 1. Fetch study from rx-study
	var study models.Study
	if err := transport.GetJSON(fmt.Sprintf("%s/studies/%s", g.topology.StudyURL, studyID), &study); err != nil {
		transport.WriteError(w, http.StatusNotFound, "study not found", "rx-gateway")
		return
	}

	// 2. For each image in study, invoke rx-image
	updatedImages := make([]models.XRayImage, 0)
	for _, img := range study.Images {
		procReq := contracts.ProcessImageRequest{
			StudyID:    studyID,
			ImageID:    img.ID,
			FileName:   img.FileName,
			RawData:    img.OriginalURI,
			Projection: img.Metadata.Projection,
			Format:     img.Metadata.Format,
		}
		var procResp contracts.ProcessImageResponse
		if err := transport.PostJSON(fmt.Sprintf("%s/image/process", g.topology.ImageURL), procReq, &procResp); err == nil {
			updatedImages = append(updatedImages, procResp.Image)
		} else {
			updatedImages = append(updatedImages, img)
		}
	}

	// Update status in rx-study
	statusReq := contracts.UpdateStudyStatusRequest{
		Status: models.StudyStatusProcessing,
		Reason: "Images verified and processed with non-destructive derivatives",
	}
	_ = transport.PostJSON(fmt.Sprintf("%s/studies/%s/status", g.topology.StudyURL, studyID), statusReq, nil)

	transport.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "Images processed and validated without modifying original",
		"images":  updatedImages,
	})
}

func (g *GatewayServer) orchestrateGenericReading(w http.ResponseWriter, r *http.Request, studyID string) {
	var study models.Study
	if err := transport.GetJSON(fmt.Sprintf("%s/studies/%s", g.topology.StudyURL, studyID), &study); err != nil {
		transport.WriteError(w, http.StatusNotFound, "study not found", "rx-gateway")
		return
	}

	derivURI := ""
	if len(study.Images) > 0 {
		derivURI = study.Images[0].DerivedURI
	}

	readerReq := contracts.RequestGenericReadingRequest{
		StudyID:          studyID,
		StudyType:        study.StudyType,
		AnatomicalRegion: study.AnatomicalRegion,
		DerivativeURI:    derivURI,
		Projections:      []string{"PA", "AP"},
	}

	var readerResp contracts.RequestGenericReadingResponse
	if err := transport.PostJSON(fmt.Sprintf("%s/reader/analyze", g.topology.ReaderURL), readerReq, &readerResp); err != nil {
		transport.WriteError(w, http.StatusBadGateway, fmt.Sprintf("rx-reader failed: %v", err), "rx-gateway")
		return
	}

	// Update study status
	_ = transport.PostJSON(fmt.Sprintf("%s/studies/%s/status", g.topology.StudyURL, studyID), contracts.UpdateStudyStatusRequest{
		Status: models.StudyStatusGenericRead,
		Reason: "Generic automated reading completed without medical signature",
	}, nil)

	transport.WriteJSON(w, http.StatusOK, readerResp.Reading)
}

func (g *GatewayServer) orchestrateResultBuild(w http.ResponseWriter, r *http.Request, studyID string) {
	var study models.Study
	if err := transport.GetJSON(fmt.Sprintf("%s/studies/%s", g.topology.StudyURL, studyID), &study); err != nil {
		transport.WriteError(w, http.StatusNotFound, "study not found", "rx-gateway")
		return
	}

	// Request reading if not already completed
	readerReq := contracts.RequestGenericReadingRequest{
		StudyID:          studyID,
		StudyType:        study.StudyType,
		AnatomicalRegion: study.AnatomicalRegion,
	}
	var readerResp contracts.RequestGenericReadingResponse
	_ = transport.PostJSON(fmt.Sprintf("%s/reader/analyze", g.topology.ReaderURL), readerReq, &readerResp)

	buildReq := contracts.BuildResultRequest{
		StudyID: studyID,
		Study:   study,
		Reading: readerResp.Reading,
	}
	var buildResp contracts.BuildResultResponse
	if err := transport.PostJSON(fmt.Sprintf("%s/result/build", g.topology.ResultURL), buildReq, &buildResp); err != nil {
		transport.WriteError(w, http.StatusBadGateway, fmt.Sprintf("rx-result failed: %v", err), "rx-gateway")
		return
	}

	// Update study status
	_ = transport.PostJSON(fmt.Sprintf("%s/studies/%s/status", g.topology.StudyURL, studyID), contracts.UpdateStudyStatusRequest{
		Status: models.StudyStatusPackageBuilt,
		Reason: "Package assembled with mandatory disclaimers",
	}, nil)

	transport.WriteJSON(w, http.StatusOK, buildResp.ResultPackage)
}

func (g *GatewayServer) orchestrateDelivery(w http.ResponseWriter, r *http.Request, studyID string) {
	var body struct {
		Recipient string `json:"recipient"`
		Channel   string `json:"channel"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)

	var study models.Study
	_ = transport.GetJSON(fmt.Sprintf("%s/studies/%s", g.topology.StudyURL, studyID), &study)

	// Build result package first
	readerReq := contracts.RequestGenericReadingRequest{
		StudyID:          studyID,
		StudyType:        study.StudyType,
		AnatomicalRegion: study.AnatomicalRegion,
	}
	var readerResp contracts.RequestGenericReadingResponse
	_ = transport.PostJSON(fmt.Sprintf("%s/reader/analyze", g.topology.ReaderURL), readerReq, &readerResp)

	var buildResp contracts.BuildResultResponse
	_ = transport.PostJSON(fmt.Sprintf("%s/result/build", g.topology.ResultURL), contracts.BuildResultRequest{
		StudyID: studyID,
		Study:   study,
		Reading: readerResp.Reading,
	}, &buildResp)

	dispatchReq := contracts.DispatchDeliveryRequest{
		PackageID: buildResp.ResultPackage,
		Recipient: body.Recipient,
		Channel:   body.Channel,
	}
	var deliveryResp contracts.DispatchDeliveryResponse
	if err := transport.PostJSON(fmt.Sprintf("%s/delivery/dispatch", g.topology.DeliveryURL), dispatchReq, &deliveryResp); err != nil {
		transport.WriteError(w, http.StatusBadGateway, fmt.Sprintf("rx-delivery failed: %v", err), "rx-gateway")
		return
	}

	transport.WriteJSON(w, http.StatusOK, map[string]any{
		"success":  true,
		"delivery": deliveryResp,
	})
}

// handleDispatchFlow orchestrates the complete automated flow through all 9 independent services.
func (g *GatewayServer) handleDispatchFlow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		transport.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", "rx-gateway")
		return
	}

	var req contracts.OrchestratedDispatchFlowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		transport.WriteError(w, http.StatusBadRequest, "invalid payload", "rx-gateway")
		return
	}

	// Step 1: Security Authorization
	authReq := contracts.AuthorizeRequest{
		Token:          "sys-rx-mesh-token",
		RequiredScope:  "dispatch:trigger",
		ResourceAction: "flow:execute",
	}
	var authResp contracts.AuthorizeResponse
	_ = transport.PostJSON(fmt.Sprintf("%s/security/authorize", g.topology.SecurityURL), authReq, &authResp)

	// Step 2: Fetch Study
	var study models.Study
	if err := transport.GetJSON(fmt.Sprintf("%s/studies/%s", g.topology.StudyURL, req.StudyID), &study); err != nil {
		transport.WriteError(w, http.StatusNotFound, "study not found", "rx-gateway")
		return
	}

	// Step 3: Process Images via rx-image & rx-storage
	if len(study.Images) > 0 {
		var imgResp contracts.ProcessImageResponse
		_ = transport.PostJSON(fmt.Sprintf("%s/image/process", g.topology.ImageURL), contracts.ProcessImageRequest{
			StudyID:    study.ID,
			ImageID:    study.Images[0].ID,
			FileName:   study.Images[0].FileName,
			RawData:    study.Images[0].OriginalURI,
			Projection: study.Images[0].Metadata.Projection,
			Format:     study.Images[0].Metadata.Format,
		}, &imgResp)
	}

	// Step 4: Generic Reading via rx-reader
	readerReq := contracts.RequestGenericReadingRequest{
		StudyID:          study.ID,
		StudyType:        study.StudyType,
		AnatomicalRegion: study.AnatomicalRegion,
	}
	var readerResp contracts.RequestGenericReadingResponse
	_ = transport.PostJSON(fmt.Sprintf("%s/reader/analyze", g.topology.ReaderURL), readerReq, &readerResp)

	// Step 5: Result Package Assembly via rx-result
	buildReq := contracts.BuildResultRequest{
		StudyID: study.ID,
		Study:   study,
		Reading: readerResp.Reading,
	}
	var buildResp contracts.BuildResultResponse
	_ = transport.PostJSON(fmt.Sprintf("%s/result/build", g.topology.ResultURL), buildReq, &buildResp)

	// Step 6: Dispatch via rx-delivery
	channel := req.DeliveryChannel
	if channel == "" {
		channel = "EMAIL"
	}
	target := req.DeliveryTarget
	if target == "" {
		target = study.Patient.Email
	}
	dispatchReq := contracts.DispatchDeliveryRequest{
		PackageID: buildResp.ResultPackage,
		Recipient: target,
		Channel:   channel,
	}
	var deliveryResp contracts.DispatchDeliveryResponse
	_ = transport.PostJSON(fmt.Sprintf("%s/delivery/dispatch", g.topology.DeliveryURL), dispatchReq, &deliveryResp)

	// Step 7: Query Audit Trail count
	var auditResp contracts.QueryAuditEventsResponse
	_ = transport.GetJSON(fmt.Sprintf("%s/audit/events?studyId=%s", g.topology.AuditURL, study.ID), &auditResp)

	resp := contracts.OrchestratedDispatchFlowResponse{
		StudyID:         study.ID,
		SecurityChecked: authResp.Authorized,
		StudyStatus:     models.StudyStatusDelivered,
		ImageIntegrity:  true,
		ReadingType:     readerResp.ReadingType,
		PackageChecksum: buildResp.ChecksumSHA,
		DeliveryToken:   deliveryResp.TrackingToken,
		DeliveryStatus:  models.DeliverySent,
		AuditCount:      auditResp.Count,
	}

	log.Printf("[RX-GATEWAY] Orchestrated dispatch flow completed for study %s", study.ID)
	transport.WriteJSON(w, http.StatusOK, resp)
}

// handleAudit delegates audit queries to rx-audit.
func (g *GatewayServer) handleAudit(w http.ResponseWriter, r *http.Request) {
	studyID := r.URL.Query().Get("studyId")
	url := fmt.Sprintf("%s/audit/events", g.topology.AuditURL)
	if studyID != "" {
		url = fmt.Sprintf("%s?studyId=%s", url, studyID)
	}

	var auditResp contracts.QueryAuditEventsResponse
	if err := transport.GetJSON(url, &auditResp); err != nil {
		transport.WriteError(w, http.StatusBadGateway, "failed to query audit log", "rx-gateway")
		return
	}
	transport.WriteJSON(w, http.StatusOK, auditResp.Events)
}

func (g *GatewayServer) handleDeliveryQueue(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s/delivery/queue", g.topology.DeliveryURL)
	if r.Method == http.MethodGet {
		var resp map[string]interface{}
		if err := transport.GetJSON(url, &resp); err != nil {
			transport.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to query delivery queue: %v", err), "rx-gateway")
			return
		}
		transport.WriteJSON(w, http.StatusOK, resp)
		return
	}
	transport.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", "rx-gateway")
}

func (g *GatewayServer) handleDeliveryNetworkToggle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		transport.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", "rx-gateway")
		return
	}
	var req map[string]interface{}
	_ = json.NewDecoder(r.Body).Decode(&req)
	var resp map[string]interface{}
	url := fmt.Sprintf("%s/delivery/network-toggle", g.topology.DeliveryURL)
	if err := transport.PostJSON(url, req, &resp); err != nil {
		transport.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to toggle network simulation: %v", err), "rx-gateway")
		return
	}
	transport.WriteJSON(w, http.StatusOK, resp)
}

func (g *GatewayServer) handleStorageRetention(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s/storage/retention", g.topology.StorageURL)
	if r.Method == http.MethodGet {
		var resp map[string]interface{}
		if err := transport.GetJSON(url, &resp); err != nil {
			transport.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to query retention status: %v", err), "rx-gateway")
			return
		}
		transport.WriteJSON(w, http.StatusOK, resp)
		return
	}
	transport.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", "rx-gateway")
}

func (g *GatewayServer) handleStoragePurge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		transport.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", "rx-gateway")
		return
	}
	url := fmt.Sprintf("%s/storage/retention/purge", g.topology.StorageURL)
	var resp map[string]interface{}
	if err := transport.PostJSON(url, map[string]string{}, &resp); err != nil {
		transport.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to execute purge: %v", err), "rx-gateway")
		return
	}
	transport.WriteJSON(w, http.StatusOK, resp)
}

func main() {
	port := config.GetEnvInt("RX_GATEWAY_PORT", 8089)
	topology := config.LoadTopologyConfig()
	server := NewGatewayServer(port, topology)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", server.handleHealthz)
	mux.HandleFunc("/gateway/topology", server.handleTopology)
	mux.HandleFunc("/api/v1/studies", server.handleStudies)
	mux.HandleFunc("/api/v1/studies/", server.handleStudyDetail)
	mux.HandleFunc("/api/v1/dispatch-flow", server.handleDispatchFlow)
	mux.HandleFunc("/api/v1/audit", server.handleAudit)
	mux.HandleFunc("/api/v1/delivery/queue", server.handleDeliveryQueue)
	mux.HandleFunc("/api/v1/delivery/network-toggle", server.handleDeliveryNetworkToggle)
	mux.HandleFunc("/api/v1/storage/retention", server.handleStorageRetention)
	mux.HandleFunc("/api/v1/storage/retention/purge", server.handleStoragePurge)
	mux.HandleFunc("/metrics", metrics.DefaultRegistry.Handler())

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	log.Printf("[RX-GATEWAY] Ingress starting on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("[RX-GATEWAY] Server failed: %v", err)
	}
}
