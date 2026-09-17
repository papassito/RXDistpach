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

// mockAuditClient es una implementación falsa del AuditClient para evitar llamadas de red en las pruebas.
type mockAuditClient struct {
	lastEvent contracts.RecordAuditEventRequest
}

// SendEvent es una implementación que captura el evento para su posterior inspección.
func (m *mockAuditClient) SendEvent(req contracts.RecordAuditEventRequest) error {
	// En una prueba más avanzada, podríamos verificar que este método fue llamado.
	m.lastEvent = req
	return nil
}

// mockResultClient es una implementación falsa del ResultClient para evitar llamadas de red en las pruebas.
type mockResultClient struct{}

// ConsolidateReading es una implementación no operativa que satisface la interfaz.
func (m *mockResultClient) ConsolidateReading(req interface{}) error {
	return nil
}

func TestAnalyzeHandler(t *testing.T) {
	// 1. Preparación (Arrange)
	// Crear una instancia del servidor con clientes falsos (mocks) que puedan capturar datos.
	mockAudit := &mockAuditClient{}
	srv := &server{
		auditClient:  mockAudit,
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

	// 2. Actuación (Act)
	// Llamar al manejador directamente.
	http.HandlerFunc(srv.analyzeHandler).ServeHTTP(rr, req)

	// 3. Afirmación (Assert)
	// Verificar que el código de estado es 200 OK.
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

	// Verificar la INVARIANTE CLÍNICA: la presencia del disclaimer.
	if response.Disclaimer.ReadingType != shared.ReadingTypeGenericAutomated {
		t.Errorf("handler did not inject correct ReadingType: got %v want %v", response.Disclaimer.ReadingType, shared.ReadingTypeGenericAutomated)
	}
	if response.Disclaimer.SignatureStatus != shared.StatusWithoutMedicalSignature {
		t.Errorf("handler did not inject correct SignatureStatus: got %v want %v", response.Disclaimer.SignatureStatus, shared.StatusWithoutMedicalSignature)
	}
	if !strings.Contains(response.Disclaimer.FullMandatoryLegalWarning, "no constituye un diagnóstico médico") {
		t.Fatalf("handler did not include the mandatory legal warning in the disclaimer")
	}

	// Verificar el evento de auditoría enviado al mock
	if mockAudit.lastEvent.Event.StudyInstanceUID != "TEST-123" {
		t.Errorf("audit event has wrong StudyInstanceUID: got %q want %q", mockAudit.lastEvent.Event.StudyInstanceUID, "TEST-123")
	}
	// La solicitud de prueba no es TLS, por lo que el estado debe ser false.
	if mockAudit.lastEvent.Event.SecurityTLSStatus != false {
		t.Errorf("audit event has wrong SecurityTLSStatus: got %v want %v", mockAudit.lastEvent.Event.SecurityTLSStatus, false)
	}

	t.Log("TestAnalyzeHandler passed: Correctly processed request and injected clinical invariant.")
}
