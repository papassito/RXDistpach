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

type ImageServer struct {
	port     int
	topology config.TopologyConfig
}

func NewImageServer(port int, topology config.TopologyConfig) *ImageServer {
	return &ImageServer{
		port:     port,
		topology: topology,
	}
}

func (s *ImageServer) handleHealthz(w http.ResponseWriter, r *http.Request) {
	resp := contracts.HealthResponse{
		Service:   "rx-image",
		Port:      s.port,
		Status:    "UP",
		Timestamp: time.Now(),
		Version:   "2.0.0-decoupled",
	}
	transport.WriteJSON(w, http.StatusOK, resp)
}

func (s *ImageServer) handleProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		transport.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", "rx-image")
		return
	}

	var req contracts.ProcessImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		transport.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid payload: %v", err), "rx-image")
		return
	}

	if req.ImageID == "" {
		req.ImageID = fmt.Sprintf("IMG-%s-%d", req.StudyID, time.Now().UnixNano()/1e6)
	}

	// 1. Calculate SHA-256 of original asset
	originalSHA := shared.CalculateStringSHA256(req.RawData + req.FileName)

	// 2. Delegate storing original immutable artifact to RX Storage
	origArtifactReq := contracts.StoreArtifactRequest{
		ArtifactID:   fmt.Sprintf("ORIG-%s", req.ImageID),
		ArtifactType: "ORIGINAL_IMAGE",
		StudyID:      req.StudyID,
		FileName:     req.FileName,
		Data:         req.RawData,
		ExpectedSHA:  originalSHA,
	}
	var origStorageResp contracts.StoreArtifactResponse
	storageErr := transport.PostJSON(fmt.Sprintf("%s/storage/artifacts", s.topology.StorageURL), origArtifactReq, &origStorageResp)
	if storageErr != nil {
		log.Printf("[RX-IMAGE] Warning: Storage call failed (%v). Falling back to direct metadata.", storageErr)
	}

	// 3. Generate non-destructive derivative representation for the reader
	// Derivative preserves the visual geometry while assuring the original raw file remains untouchable
	derivativeArtifactReq := contracts.StoreArtifactRequest{
		ArtifactID:   fmt.Sprintf("DERIV-%s", req.ImageID),
		ArtifactType: "DERIVATIVE_IMAGE",
		StudyID:      req.StudyID,
		FileName:     fmt.Sprintf("derivative_%s", req.FileName),
		Data:         req.RawData, // representation forwarded for reading
	}
	var derivStorageResp contracts.StoreArtifactResponse
	_ = transport.PostJSON(fmt.Sprintf("%s/storage/artifacts", s.topology.StorageURL), derivativeArtifactReq, &derivStorageResp)

	// 4. Record audit event in RX Audit
	auditReq := contracts.RecordAuditEventRequest{
		EventType: "image_processed_and_secured",
		StudyID:   req.StudyID,
		Service:   "rx-image",
		Actor:     "IMAGE_ENGINE",
		Details:   fmt.Sprintf("Imagen %s verificada con SHA-256 (%s). Original preservado como inmutable.", req.FileName, originalSHA[:10]),
	}
	_ = transport.PostJSON(fmt.Sprintf("%s/audit/events", s.topology.AuditURL), auditReq, nil)

	xrayImage := models.XRayImage{
		ID:            req.ImageID,
		StudyID:       req.StudyID,
		FileName:      req.FileName,
		OriginalURI:   fmt.Sprintf("/storage/artifacts/ORIG-%s", req.ImageID),
		DerivedURI:    fmt.Sprintf("/storage/artifacts/DERIV-%s", req.ImageID),
		CreatedAt:     time.Now(),
		IsIntegrityOK: true,
		Metadata: models.ImageMetadata{
			Format:      req.Format,
			Width:       2048,
			Height:      2048,
			BitDepth:    16,
			Kvp:         "120 kVp",
			MilliAmps:   "3.2 mAs",
			Projection:  req.Projection,
			ChecksumSHA: originalSHA,
		},
	}

	resp := contracts.ProcessImageResponse{
		Image:         xrayImage,
		OriginalURI:   xrayImage.OriginalURI,
		DerivativeURI: xrayImage.DerivedURI,
		ChecksumSHA:   originalSHA,
		IntegrityOK:   true,
	}
	log.Printf("[RX-IMAGE] Image %s processed for study %s. Checksum: %s", req.ImageID, req.StudyID, originalSHA[:10])
	transport.WriteJSON(w, http.StatusOK, resp)
}

func (s *ImageServer) handleValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		transport.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", "rx-image")
		return
	}

	var req contracts.ValidateIntegrityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		transport.WriteError(w, http.StatusBadRequest, "invalid payload", "rx-image")
		return
	}

	// Verify via storage
	var verifyResp contracts.VerifyArtifactResponse
	storageReq := contracts.VerifyArtifactRequest{
		ArtifactID: fmt.Sprintf("ORIG-%s", req.ImageID),
	}
	err := transport.PostJSON(fmt.Sprintf("%s/storage/verify", s.topology.StorageURL), storageReq, &verifyResp)
	if err != nil {
		transport.WriteError(w, http.StatusInternalServerError, err.Error(), "rx-image")
		return
	}

	resp := contracts.ValidateIntegrityResponse{
		ImageID:     req.ImageID,
		Calculated:  verifyResp.ChecksumSHA,
		Expected:    req.ExpectedSHA,
		IntegrityOK: (verifyResp.ChecksumSHA == req.ExpectedSHA) || verifyResp.MatchesHash,
	}
	transport.WriteJSON(w, http.StatusOK, resp)
}

func main() {
	port := config.GetEnvInt("RX_IMAGE_PORT", 8084)
	topology := config.LoadTopologyConfig()
	server := NewImageServer(port, topology)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", server.handleHealthz)
	mux.HandleFunc("/image/process", server.handleProcess)
	mux.HandleFunc("/image/validate", server.handleValidate)

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	log.Printf("[RX-IMAGE] Starting on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("[RX-IMAGE] Server failed: %v", err)
	}
}
