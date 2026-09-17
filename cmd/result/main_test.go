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

// mockDeliveryClient es una implementación falsa del DeliveryClient para pruebas.
type mockDeliveryClient struct {
	callCount int
}

// EnqueueResult simula el envío, contando las veces que es llamado.
func (m *mockDeliveryClient) EnqueueResult(req interface{}) error {
	m.callCount++
	return nil
}

func TestConsolidateHandler(t *testing.T) {
	// --- Caso de éxito ---
	// Crear una lectura genérica de prueba.
	testReading := models.GenericReading{
		StudyID:          "RESULT-TEST-001",
		AnatomicalRegion: "SKULL",
		Timestamp:        time.Now().UTC().Format(time.RFC3339),
		ReadingContent:   "Test content for result consolidation.",
		Disclaimer: models.Disclaimer{
			ReadingType:               shared.ReadingTypeGenericAutomated,
			SignatureStatus:           shared.StatusWithoutMedicalSignature,
			MedicalReportStatus:       shared.MedicalReportNotIncluded,
			FullMandatoryLegalWarning: shared.MandatoryLegalDisclaimer,
		},
	}
	body, _ := json.Marshal(testReading)

	// Crear una solicitud HTTP falsa.
	reqPost := httptest.NewRequest(http.MethodPost, "/result/consolidate", bytes.NewReader(body))
	rrPost := httptest.NewRecorder()

	// Crear el servidor con el mock client.
	mockClient := &mockDeliveryClient{}
	srv := &server{
		deliveryClient: mockClient,
	}

	// Llamar al manejador directamente.
	srv.consolidateHandler(rrPost, reqPost)

	// Verificar que el código de estado es 200 OK.
	// Para este manejador (que actualmente solo registra en log),
	// una respuesta 200 OK es suficiente para confirmar que recibió y
	// decodificó el cuerpo de la solicitud correctamente sin errores.
	if status := rrPost.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// Verificar que el cliente de delivery fue llamado.
	if mockClient.callCount < 1 {
		t.Errorf("delivery client was not called")
	}

	// --- Caso de Error: Método Incorrecto ---
	reqGet := httptest.NewRequest(http.MethodGet, "/result/consolidate", nil)
	rrGet := httptest.NewRecorder()
	srv.consolidateHandler(rrGet, reqGet)
	if status := rrGet.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code for GET request: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}

	t.Log("TestConsolidateHandler passed: Correctly handled valid POST, called delivery client, and rejected invalid GET.")
}
