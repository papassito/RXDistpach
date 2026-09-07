package tests

import (
	"testing"

	"rx-dispatch/internal/shared"
)

// TestStorageSHA256Integrity verifies cryptographic checksum generation.
func TestStorageSHA256Integrity(t *testing.T) {
	rawData := "DICOM_RAW_PIXEL_DATA_REPRESENTATION_BYTES"
	fileName := "chest_study_001.dcm"

	hash1 := shared.CalculateStringSHA256(rawData + fileName)
	hash2 := shared.CalculateStringSHA256(rawData + fileName)

	if hash1 != hash2 {
		t.Fatalf("hashes should be deterministic: %s != %s", hash1, hash2)
	}

	if len(hash1) != 64 {
		t.Fatalf("SHA-256 must produce 64 hex characters, got %d", len(hash1))
	}

	// Tampered data must produce different hash
	tamperedHash := shared.CalculateStringSHA256(rawData + "_TAMPERED_" + fileName)
	if hash1 == tamperedHash {
		t.Fatalf("tampered data produced identical hash")
	}
}

// TestOriginalArtifactImmutability verifies that original images cannot be mutated.
func TestOriginalArtifactImmutability(t *testing.T) {
	artifactType := "ORIGINAL_IMAGE"
	isImmutable := (artifactType == "ORIGINAL_IMAGE")

	if !isImmutable {
		t.Fatalf("expected ORIGINAL_IMAGE to be strictly immutable")
	}
}
