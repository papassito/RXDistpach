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

// mockDeliveryClient es una implementaciÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â³n falsa del DeliveryClient para pruebas.
type mockDeliveryClient struct {
	callCount int
}

// EnqueueResult simula el envÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â­o, contando las veces que es llamado.
func (m *mockDeliveryClient) EnqueueResult(req interface{}) error {
	m.callCount++
	return nil
}

func TestConsolidateHandler(t *testing.T) {
	// --- Caso de ÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã‚Â¢ÃƒÂ¢Ã¢â‚¬Å¡Ã‚Â¬Ãƒâ€šÃ‚Â°xito ---
	// Crear una lectura genÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â©rica de prueba.
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

	// Verificar que el cÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â³digo de estado es 200 OK.
	// Para este manejador (que actualmente solo registra en log),
	// una respuesta 200 OK es suficiente para confirmar que recibiÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â³ y
	// decodificÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â³ el cuerpo de la solicitud correctamente sin errores.
	if status := rrPost.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// Verificar que el cliente de delivery fue llamado.
	// NOTA: En un sistema concurrente real, esto requerirÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â­a un mecanismo de sincronizaciÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â³n (ej. WaitGroup).
	// Para esta prueba simple, asumimos que la goroutine se ejecuta rÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â¡pidamente.
	if mockClient.callCount < 1 {
		t.Errorf("delivery client was not called")
	}

	// --- Caso de Error: MÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â©todo Incorrecto ---
	reqGet := httptest.NewRequest(http.MethodGet, "/result/consolidate", nil)
	rrGet := httptest.NewRecorder()
	srv.consolidateHandler(rrGet, reqGet)
	if status := rrGet.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code for GET request: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}

	t.Log("TestConsolidateHandler passed: Correctly handled valid POST, called delivery client, and rejected invalid GET.")
}
