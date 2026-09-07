package tests

import (
	"testing"

	"rx-dispatch/internal/image"
)

func TestWebLightDerivativeGeneration(t *testing.T) {
	origArtifactID := "ORIG-TEST-001"
	rawBytes := []byte("SIMULATED_DICOM_HIGH_RESOLUTION_PIXEL_DATA_16BIT")
	rawSHA := "d41d8cd98f00b204e9800998ecf8427e0123456789abcdef0123456789abcdef"

	deriv := image.GenerateWebLightDerivative(origArtifactID, rawBytes, rawSHA)

	if deriv.OriginalArtifactID != origArtifactID {
		t.Fatalf("expected OriginalArtifactID %s, got %s", origArtifactID, deriv.OriginalArtifactID)
	}
	if deriv.OriginalSHA256 != rawSHA {
		t.Fatalf("expected OriginalSHA256 %s, got %s", rawSHA, deriv.OriginalSHA256)
	}
	if deriv.DerivativeSHA256 == "" {
		t.Fatalf("expected non-empty DerivativeSHA256")
	}
	if deriv.DerivativeSHA256 == deriv.OriginalSHA256 {
		t.Fatalf("derivative SHA must be distinct from original SHA")
	}
	if deriv.Format != "WEBP_MEDICAL_OPTIMIZED" {
		t.Fatalf("expected format WEBP_MEDICAL_OPTIMIZED, got %s", deriv.Format)
	}
	if deriv.Dimensions != "1024x1024" {
		t.Fatalf("expected dimensions 1024x1024, got %s", deriv.Dimensions)
	}
}
