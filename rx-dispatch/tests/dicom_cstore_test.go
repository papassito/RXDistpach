package tests

import (
	"testing"
	"time"

	"rx-dispatch/internal/dicom"
)

func TestDICOMCStoreSCPIngestion(t *testing.T) {
	testPort := 11188
	var receivedDataset *dicom.Dataset

	// 1. Initialize C-STORE SCP listener
	server := dicom.NewSCPServer(
		dicom.ServerConfig{
			AETitle: "RX_DISPATCH",
			Port:    testPort,
		},
		func(ds *dicom.Dataset) error {
			receivedDataset = ds
			return nil
		},
	)

	if err := server.Start(); err != nil {
		t.Fatalf("failed to start DICOM SCP on port %d: %v", testPort, err)
	}
	defer server.Stop()

	time.Sleep(50 * time.Millisecond)

	// 2. Generate authentic synthetic DICOM image with standard tags
	expectedPatientName := "GARCIA^MARIA"
	expectedPatientID := "PAC-DCM-7701"
	expectedAccession := "ACC-2026-X99"
	expectedStudyDesc := "RX TORAX AP DIGITAL"

	dicomBytes := dicom.GenerateSyntheticDICOM(
		expectedPatientName,
		expectedPatientID,
		expectedAccession,
		expectedStudyDesc,
	)

	// 3. Send over TCP using DICOM Upper Layer Protocol SCU client
	err := dicom.SendCStore(
		"127.0.0.1",
		testPort,
		"SIEMENS_YSIO_XR", // Calling AE Title
		"RX_DISPATCH",     // Called AE Title
		dicomBytes,
	)
	if err != nil {
		t.Fatalf("C-STORE transmission failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	// 4. Verify received dataset in SCP listener
	if receivedDataset == nil {
		t.Fatalf("expected callback receivedDataset to be non-nil")
	}
	instances, last := server.Stats()
	if instances != 1 {
		t.Fatalf("expected 1 received DICOM instance, got %d", instances)
	}
	if last == nil {
		t.Fatalf("expected non-nil last received dataset")
	}

	if last.PatientName != expectedPatientName {
		t.Fatalf("expected PatientName '%s', got '%s'", expectedPatientName, last.PatientName)
	}
	if last.PatientID != expectedPatientID {
		t.Fatalf("expected PatientID '%s', got '%s'", expectedPatientID, last.PatientID)
	}
	if last.Modality != "DX" {
		t.Fatalf("expected Modality 'DX', got '%s'", last.Modality)
	}
	if last.StudyDescription != expectedStudyDesc {
		t.Fatalf("expected StudyDescription '%s', got '%s'", expectedStudyDesc, last.StudyDescription)
	}
	if last.ChecksumSHA == "" {
		t.Fatalf("expected computed SHA-256 for received DICOM payload")
	}
	if last.SizeBytes != len(dicomBytes) {
		t.Fatalf("expected byte size %d, got %d", len(dicomBytes), last.SizeBytes)
	}
}
