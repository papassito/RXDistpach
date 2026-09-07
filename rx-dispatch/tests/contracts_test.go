package tests

import (
	"encoding/json"
	"testing"
	"time"

	"rx-dispatch/internal/contracts"
	"rx-dispatch/internal/models"
)

// TestContractSerialization verifies JSON contract stability across microservices.
func TestContractSerialization(t *testing.T) {
	studyReq := contracts.IngestStudyRequest{
		StudyIdentifier:    "RX-2026-TEST",
		StudyType:          "Tórax AP",
		AnatomicalRegion:   "Tórax",
		ReferringPhysician: "Dr. Test",
		ClinicalIndication: "Control",
		Patient: models.Patient{
			ID:       "PAC-01",
			FullName: "Paciente Prueba",
			Email:    "paciente@test.com",
		},
	}

	data, err := json.Marshal(studyReq)
	if err != nil {
		t.Fatalf("failed to marshal IngestStudyRequest: %v", err)
	}

	var decoded contracts.IngestStudyRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal IngestStudyRequest: %v", err)
	}

	if decoded.StudyIdentifier != studyReq.StudyIdentifier {
		t.Fatalf("mismatch in decoded identifier: %s != %s", decoded.StudyIdentifier, studyReq.StudyIdentifier)
	}
}

// TestHealthResponseContract verifies health contract consistency.
func TestHealthResponseContract(t *testing.T) {
	health := contracts.HealthResponse{
		Service:   "rx-reader",
		Port:      8085,
		Status:    "UP",
		Timestamp: time.Now(),
		Version:   "2.0.0-decoupled",
	}

	data, err := json.Marshal(health)
	if err != nil {
		t.Fatalf("failed to marshal health response: %v", err)
	}

	var decoded contracts.HealthResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal health response: %v", err)
	}

	if decoded.Service != "rx-reader" || decoded.Status != "UP" {
		t.Fatalf("invalid decoded health values")
	}
}
