package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"rx-dispatch/internal/contracts"
	"rx-dispatch/internal/models"
	"rx-dispatch/internal/shared"
)

// mockAuditClient es una implementaciÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â³n falsa del AuditClient para evitar llamadas de red en las pruebas.
type mockAuditClient struct{}

// SendEvent es una implementaciÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â³n no operativa que satisface la interfaz.
func (m *mockAuditClient) SendEvent(event contracts.RecordAuditEventRequest) {
	// En una prueba mÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â¡s avanzada, podrÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â­amos verificar que este mÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â©todo fue llamado.
	// Por ahora, simplemente evitamos la llamada de red.
}

// mockResultClient es una implementaciÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â³n falsa del ResultClient para evitar llamadas de red en las pruebas.
type mockResultClient struct{}

// ConsolidateReading es una implementaciÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â³n no operativa que satisface la interfaz.
func (m *mockResultClient) ConsolidateReading(reading models.GenericReading) {
	// no-op
}

func TestAnalyzeHandler(t *testing.T) {
	// 1. PreparaciÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â³n (Arrange)
	// Crear una instancia del servidor con clientes falsos (mocks).
	srv := &server{
		auditClient:  &mockAuditClient{},
		resultClient: &mockResultClient{},
	}

	// Crear el cuerpo de la solicitud JSON.
	requestBody := contracts.RequestGenericReadingRequest{
		StudyID:          "TEST-123",
		StudyType:        "CR",
		AnatomicalRegion: "CHEST",
	}
	body, _ := json.Marshal(requestBody)

	// Crear una solicitud HTTP falsa.
	req := httptest.NewRequest(http.MethodPost, "/reader/analyze", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	// 2. ActuaciÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â³n (Act)
	// Llamar al manejador directamente.
	srv.analyzeHandler(rr, req)

	// 3. AfirmaciÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â³n (Assert)
	// Verificar que el cÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â³digo de estado es 200 OK.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Decodificar la respuesta JSON para verificar su contenido.
	var response models.GenericReading
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Could not decode response body: %v", err)
	}

	// Verificar que el StudyID coincide.
	if response.StudyID != "TEST-123" {
		t.Errorf("handler returned unexpected studyId: got %v want %v", response.StudyID, "TEST-123")
	}

	// Verificar la INVARIANTE CLÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚ÂNICA: la presencia del disclaimer.
	if response.Disclaimer.ReadingType != shared.ReadingTypeGenericAutomated {
		t.Errorf("handler did not inject correct ReadingType: got %v want %v", response.Disclaimer.ReadingType, shared.ReadingTypeGenericAutomated)
	}
	if response.Disclaimer.SignatureStatus != shared.StatusWithoutMedicalSignature {
		t.Errorf("handler did not inject correct SignatureStatus: got %v want %v", response.Disclaimer.SignatureStatus, shared.StatusWithoutMedicalSignature)
	}
	if !strings.Contains(response.Disclaimer.FullMandatoryLegalWarning, "no constituye un diagnÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â³stico mÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â©dico") {
		t.Fatalf("handler did not include the mandatory legal warning in the disclaimer")
	}

	t.Log("TestAnalyzeHandler passed: Correctly processed request and injected clinical invariant.")
}