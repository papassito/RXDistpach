package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"rx-dispatch/internal/models"
	"rx-dispatch/internal/shared"
)

func TestEnqueueHandler(t *testing.T) {
	// --- Caso de éxito ---
	testReading := models.GenericReading{
		StudyID:          "DELIVERY-TEST-001",
		AnatomicalRegion: "KNEE",
		Timestamp:        time.Now().UTC().Format(time.RFC3339),
		ReadingContent:   "Test content for delivery queue.",
		Disclaimer: models.Disclaimer{
			ReadingType:               shared.ReadingTypeGenericAutomated,
			SignatureStatus:           shared.StatusWithoutMedicalSignature,
			MedicalReportStatus:       shared.MedicalReportNotIncluded,
			FullMandatoryLegalWarning: shared.MandatoryLegalDisclaimer,
		},
	}
	body, _ := json.Marshal(testReading)

	reqPost := httptest.NewRequest(http.MethodPost, "/delivery/enqueue", bytes.NewReader(body))
	rrPost := httptest.NewRecorder()

	srv := &server{}
	srv.enqueueHandler(rrPost, reqPost)

	// Verificar que el código de estado es 202 Accepted.
	if status := rrPost.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusAccepted)
	}

	// --- Caso de Error: Método Incorrecto ---
	reqGet := httptest.NewRequest(http.MethodGet, "/delivery/enqueue", nil)
	rrGet := httptest.NewRecorder()
	srv.enqueueHandler(rrGet, reqGet)
	if status := rrGet.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code for GET request: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}

	t.Log("TestEnqueueHandler passed: Correctly handled valid POST and rejected invalid GET.")
}
