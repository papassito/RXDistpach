package tests

import (
	"testing"
	"time"

	"rx-dispatch/internal/storage"
)

func TestStorageRetentionPurgePolicy(t *testing.T) {
	mgr := storage.NewRetentionManager(&storage.RetentionConfig{
		DerivativeRetentionDays: 30, // 30 days for web thumbnails/previews
		OriginalsProtected:      true,
		MaxDiskUsageBytes:       100 * 1024 * 1024 * 1024,
		WatermarkPercent:        80,
	})

	now := time.Now()
	expiredDate := now.Add(-60 * 24 * time.Hour) // 60 days ago
	freshDate := now.Add(-5 * 24 * time.Hour)    // 5 days ago

	artifacts := map[string]storage.ArtifactMetadata{
		// 1. Clinical Original (Old, but MUST BE PRESERVED ACCORDING TO CLINICAL MANDATE)
		"ORIG-TORAX-001": {
			ID:          "ORIG-TORAX-001",
			Type:        "ORIGINAL_IMAGE",
			StudyID:     "STU-001",
			FileName:    "torax_raw.dcm",
			SizeBytes:   15 * 1024 * 1024, // 15MB
			StoredAt:    expiredDate,
			IsImmutable: true, // IMMUTABLE
		},
		// 2. DICOM Raw (Old, but MUST BE PRESERVED)
		"DICOM-TORAX-002": {
			ID:          "DICOM-TORAX-002",
			Type:        "DICOM_RAW",
			StudyID:     "STU-002",
			FileName:    "study.dcm",
			SizeBytes:   20 * 1024 * 1024, // 20MB
			StoredAt:    expiredDate,
			IsImmutable: true, // IMMUTABLE
		},
		// 3. Expired derivative (60 days old, Max is 30 -> ELIGIBLE FOR PURGE)
		"DERIV-EXPIRED-001": {
			ID:          "DERIV-EXPIRED-001",
			Type:        "DERIVATIVE_IMAGE",
			StudyID:     "STU-001",
			FileName:    "preview_expired.webp",
			SizeBytes:   500 * 1024, // 500KB
			StoredAt:    expiredDate,
			IsImmutable: false, // Non-immutable
		},
		// 4. Fresh derivative (5 days old -> MUST BE PRESERVED)
		"DERIV-FRESH-002": {
			ID:          "DERIV-FRESH-002",
			Type:        "DERIVATIVE_IMAGE",
			StudyID:     "STU-002",
			FileName:    "preview_fresh.webp",
			SizeBytes:   400 * 1024, // 400KB
			StoredAt:    freshDate,
			IsImmutable: false,
		},
	}

	purgedIDs := make(map[string]bool)
	purgeFunc := func(id string) {
		purgedIDs[id] = true
		delete(artifacts, id)
	}

	res := mgr.EvaluateAndPurge(artifacts, purgeFunc)

	// Verify purge results
	if res.PurgedDerivatives != 1 {
		t.Fatalf("expected exactly 1 derivative purged, got %d", res.PurgedDerivatives)
	}
	if !purgedIDs["DERIV-EXPIRED-001"] {
		t.Fatalf("expected DERIV-EXPIRED-001 to be purged")
	}

	// Verify strict preservation of immutables
	if purgedIDs["ORIG-TORAX-001"] || purgedIDs["DICOM-TORAX-002"] {
		t.Fatalf("CRITICAL ERROR: Immutable clinical original was purged!")
	}
	if res.PreservedOriginals < 2 {
		t.Fatalf("expected at least 2 originals preserved, got %d", res.PreservedOriginals)
	}

	// Verify fresh derivative was not purged
	if purgedIDs["DERIV-FRESH-002"] {
		t.Fatalf("fresh derivative should not have been purged")
	}
}
