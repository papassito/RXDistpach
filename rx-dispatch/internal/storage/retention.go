package storage

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// ArtifactMetadata provides retention evaluation criteria.
type ArtifactMetadata struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"` // "ORIGINAL_IMAGE", "DERIVATIVE_IMAGE", "RESULT_PDF", "DICOM_RAW"
	StudyID     string    `json:"studyId"`
	FileName    string    `json:"fileName"`
	SizeBytes   int       `json:"sizeBytes"`
	StoredAt    time.Time `json:"storedAt"`
	IsImmutable bool      `json:"isImmutable"`
}

// RetentionConfig defines disk management and purge boundaries.
type RetentionConfig struct {
	DerivativeRetentionDays int   `json:"derivativeRetentionDays"` // Default: 30 days
	OriginalsProtected      bool  `json:"originalsProtected"`      // Default: true (Immutable)
	MaxDiskUsageBytes       int64 `json:"maxDiskUsageBytes"`       // Limit in bytes
	WatermarkPercent        int   `json:"watermarkPercent"`        // Trigger purge at e.g. 80%
}

// PurgeResult details the outcome of an automated or manual retention purge.
type PurgeResult struct {
	EvaluatedCount     int       `json:"evaluatedCount"`
	PurgedDerivatives  int       `json:"purgedDerivatives"`
	PreservedOriginals int       `json:"preservedOriginals"`
	BytesFreed         int64     `json:"bytesFreed"`
	RemainingBytes     int64     `json:"remainingBytes"`
	OriginalBytesTotal int64     `json:"originalBytesTotal"`
	ExecutedAt         time.Time `json:"executedAt"`
	TriggerReason      string    `json:"triggerReason"`
}

// RetentionManager controls storage lifecycle and prevents disk overflow.
type RetentionManager struct {
	mu     sync.RWMutex
	config RetentionConfig
}

// NewRetentionManager initializes disk retention rules.
func NewRetentionManager(cfg *RetentionConfig) *RetentionManager {
	if cfg == nil {
		cfg = &RetentionConfig{
			DerivativeRetentionDays: 30,
			OriginalsProtected:      true,
			MaxDiskUsageBytes:       50 * 1024 * 1024 * 1024, // 50 GB default
			WatermarkPercent:        80,
		}
	}
	return &RetentionManager{
		config: *cfg,
	}
}

// GetConfig returns current retention parameters.
func (rm *RetentionManager) GetConfig() RetentionConfig {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.config
}

// UpdateConfig modifies retention parameters.
func (rm *RetentionManager) UpdateConfig(cfg RetentionConfig) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.config = cfg
	log.Printf("[STORAGE-RETENTION] Policy updated: Derivatives retention=%d days, Originals protected=%v",
		cfg.DerivativeRetentionDays, cfg.OriginalsProtected)
}

// EvaluateAndPurge purges eligible derivatives based on age and immutability rules.
// CRITICAL CLINICAL RULE: Originals marked IsImmutable are NEVER purged.
func (rm *RetentionManager) EvaluateAndPurge(artifacts map[string]ArtifactMetadata, deleteFunc func(id string)) PurgeResult {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	now := time.Now()
	cutoffDate := now.AddDate(0, 0, -rm.config.DerivativeRetentionDays)

	res := PurgeResult{
		EvaluatedCount: len(artifacts),
		ExecutedAt:     now,
		TriggerReason:  fmt.Sprintf("Retention policy expiration (%d days cutoff)", rm.config.DerivativeRetentionDays),
	}

	for id, meta := range artifacts {
		// Strict immutability protection for clinical originals
		if meta.IsImmutable || meta.Type == "ORIGINAL_IMAGE" || meta.Type == "DICOM_RAW" {
			res.PreservedOriginals++
			res.OriginalBytesTotal += int64(meta.SizeBytes)
			res.RemainingBytes += int64(meta.SizeBytes)
			continue
		}

		// Candidate for purge: expired non-immutable derivatives or temporary caches
		if meta.StoredAt.Before(cutoffDate) {
			res.PurgedDerivatives++
			res.BytesFreed += int64(meta.SizeBytes)
			if deleteFunc != nil {
				deleteFunc(id)
			}
			log.Printf("[STORAGE-RETENTION] Purged expired derivative %s (%s, %d bytes)", id, meta.FileName, meta.SizeBytes)
		} else {
			res.RemainingBytes += int64(meta.SizeBytes)
		}
	}

	log.Printf("[STORAGE-RETENTION] Purge completed: Freed %d bytes across %d derivatives. Preserved %d clinical originals.",
		res.BytesFreed, res.PurgedDerivatives, res.PreservedOriginals)
	return res
}
