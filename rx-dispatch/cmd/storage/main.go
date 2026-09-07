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
	"rx-dispatch/internal/shared"
	"rx-dispatch/internal/storage"
	"rx-dispatch/internal/transport"
)

type StoredArtifact struct {
	ID          string
	Type        string
	StudyID     string
	FileName    string
	Data        string
	ChecksumSHA string
	SizeBytes   int
	StoredAt    time.Time
	IsImmutable bool
}

type StorageServer struct {
	mu        sync.RWMutex
	artifacts map[string]StoredArtifact
	port      int
	retention *storage.RetentionManager
}

func NewStorageServer(port int) *StorageServer {
	return &StorageServer{
		artifacts: make(map[string]StoredArtifact),
		port:      port,
		retention: storage.NewRetentionManager(nil),
	}
}

func (s *StorageServer) handleHealthz(w http.ResponseWriter, r *http.Request) {
	resp := contracts.HealthResponse{
		Service:   "rx-storage",
		Port:      s.port,
		Status:    "UP",
		Timestamp: time.Now(),
		Version:   "2.0.0-decoupled",
	}
	transport.WriteJSON(w, http.StatusOK, resp)
}

func (s *StorageServer) handleArtifacts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var req contracts.StoreArtifactRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			transport.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid payload: %v", err), "rx-storage")
			return
		}

		if req.ArtifactID == "" {
			req.ArtifactID = fmt.Sprintf("ART-%d", time.Now().UnixNano()/1e6)
		}

		calculatedSHA := shared.CalculateStringSHA256(req.Data + req.FileName)
		if req.ExpectedSHA != "" && req.ExpectedSHA != calculatedSHA {
			transport.WriteError(w, http.StatusConflict, "integrity hash mismatch", "rx-storage")
			return
		}

		isImmutable := strings.HasPrefix(req.ArtifactType, "ORIGINAL")

		s.mu.Lock()
		// If already exists and is immutable, prevent overwrite
		if existing, exists := s.artifacts[req.ArtifactID]; exists && existing.IsImmutable {
			s.mu.Unlock()
			transport.WriteError(w, http.StatusForbidden, "cannot overwrite immutable original artifact", "rx-storage")
			return
		}

		item := StoredArtifact{
			ID:          req.ArtifactID,
			Type:        req.ArtifactType,
			StudyID:     req.StudyID,
			FileName:    req.FileName,
			Data:        req.Data,
			ChecksumSHA: calculatedSHA,
			SizeBytes:   len(req.Data),
			StoredAt:    time.Now(),
			IsImmutable: isImmutable,
		}
		s.artifacts[req.ArtifactID] = item
		s.mu.Unlock()

		log.Printf("[RX-STORAGE] Stored artifact %s (Type: %s, SHA: %s, Immutable: %v)", item.ID, item.Type, calculatedSHA[:10], isImmutable)

		resp := contracts.StoreArtifactResponse{
			ArtifactID:  item.ID,
			StorageURI:  fmt.Sprintf("/storage/artifacts/%s", item.ID),
			ChecksumSHA: calculatedSHA,
			SizeBytes:   item.SizeBytes,
			StoredAt:    item.StoredAt,
			IsImmutable: isImmutable,
		}
		transport.WriteJSON(w, http.StatusCreated, resp)

	case http.MethodGet:
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(pathParts) < 3 {
			transport.WriteError(w, http.StatusBadRequest, "missing artifact id", "rx-storage")
			return
		}
		artID := pathParts[2]

		s.mu.RLock()
		item, exists := s.artifacts[artID]
		s.mu.RUnlock()

		if !exists {
			transport.WriteError(w, http.StatusNotFound, "artifact not found", "rx-storage")
			return
		}

		resp := contracts.GetArtifactResponse{
			ArtifactID:   item.ID,
			ArtifactType: item.Type,
			StudyID:      item.StudyID,
			FileName:     item.FileName,
			Data:         item.Data,
			ChecksumSHA:  item.ChecksumSHA,
		}
		transport.WriteJSON(w, http.StatusOK, resp)

	default:
		transport.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", "rx-storage")
	}
}

func (s *StorageServer) handleVerify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		transport.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", "rx-storage")
		return
	}

	var req contracts.VerifyArtifactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		transport.WriteError(w, http.StatusBadRequest, "invalid verify request", "rx-storage")
		return
	}

	s.mu.RLock()
	item, exists := s.artifacts[req.ArtifactID]
	s.mu.RUnlock()

	if !exists {
		transport.WriteJSON(w, http.StatusOK, contracts.VerifyArtifactResponse{
			ArtifactID:  req.ArtifactID,
			Exists:      false,
			MatchesHash: false,
		})
		return
	}

	calculated := shared.CalculateStringSHA256(item.Data + item.FileName)
	matches := (calculated == item.ChecksumSHA)

	transport.WriteJSON(w, http.StatusOK, contracts.VerifyArtifactResponse{
		ArtifactID:  item.ID,
		Exists:      true,
		ChecksumSHA: item.ChecksumSHA,
		MatchesHash: matches,
		IsImmutable: item.IsImmutable,
	})
}

func (s *StorageServer) handleRetention(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		transport.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", "rx-storage")
		return
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	var totalBytes int64
	var originalBytes int64
	var derivativeBytes int64
	var origCount int
	var derivCount int

	for _, a := range s.artifacts {
		totalBytes += int64(a.SizeBytes)
		if a.IsImmutable || a.Type == "ORIGINAL_IMAGE" {
			originalBytes += int64(a.SizeBytes)
			origCount++
		} else {
			derivativeBytes += int64(a.SizeBytes)
			derivCount++
		}
	}

	transport.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"policy":          s.retention.GetConfig(),
		"totalArtifacts":  len(s.artifacts),
		"originalCount":   origCount,
		"derivativeCount": derivCount,
		"totalBytes":      totalBytes,
		"originalBytes":   originalBytes,
		"derivativeBytes": derivativeBytes,
	})
}

func (s *StorageServer) handlePurge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		transport.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", "rx-storage")
		return
	}

	s.mu.Lock()
	metaMap := make(map[string]storage.ArtifactMetadata)
	for id, a := range s.artifacts {
		metaMap[id] = storage.ArtifactMetadata{
			ID:          a.ID,
			Type:        a.Type,
			StudyID:     a.StudyID,
			FileName:    a.FileName,
			SizeBytes:   a.SizeBytes,
			StoredAt:    a.StoredAt,
			IsImmutable: a.IsImmutable,
		}
	}

	purgeRes := s.retention.EvaluateAndPurge(metaMap, func(id string) {
		delete(s.artifacts, id)
	})
	s.mu.Unlock()

	transport.WriteJSON(w, http.StatusOK, purgeRes)
}

func main() {
	port := config.GetEnvInt("RX_STORAGE_PORT", 8083)
	server := NewStorageServer(port)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", server.handleHealthz)
	mux.HandleFunc("/storage/artifacts", server.handleArtifacts)
	mux.HandleFunc("/storage/artifacts/", server.handleArtifacts)
	mux.HandleFunc("/storage/verify", server.handleVerify)
	mux.HandleFunc("/storage/retention", server.handleRetention)
	mux.HandleFunc("/storage/retention/purge", server.handlePurge)

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	log.Printf("[RX-STORAGE] Starting on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("[RX-STORAGE] Server failed: %v", err)
	}
}
