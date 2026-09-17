package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"rx-dispatch/internal/contracts"
	"rx-dispatch/internal/models"
)

func TestHealthzHandler(t *testing.T) {
	t.Run("should return 200 OK for GET request", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/healthz", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(healthCheckHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	})

	t.Run("should return 405 Method Not Allowed for POST request", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/healthz", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(healthCheckHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusMethodNotAllowed)
		}
	})
}

func TestEventHandler(t *testing.T) {
	t.Run("should return 202 for valid POST request", func(t *testing.T) {
		event := models.AuditEvent{
			EventTimestamp: time.Now(),
			EventAction:    "TEST_ACTION",
			EventOutcome:   "SUCCESS",
			UserID:         "test-user",
		}
		reqBody := contracts.RecordAuditEventRequest{Event: event}
		body, _ := json.Marshal(reqBody)

		req, err := http.NewRequest(http.MethodPost, "/audit/event", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(eventHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusAccepted {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusAccepted)
		}
	})

	t.Run("should return 405 for GET request", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/audit/event", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(eventHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
		}
	})

	t.Run("should return 400 for malformed JSON", func(t *testing.T) {
		body := []byte(`{"event": "not a struct"`)
		req, err := http.NewRequest(http.MethodPost, "/audit/event", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(eventHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("should return 400 for incomplete payload", func(t *testing.T) {
		// JSON válido, pero lógicamente inválido (evento vacío)
		reqBody := contracts.RecordAuditEventRequest{Event: models.AuditEvent{}}
		body, _ := json.Marshal(reqBody)

		req, err := http.NewRequest(http.MethodPost, "/audit/event", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(eventHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("should return 413 for body too large", func(t *testing.T) {
		largeBody := strings.Repeat("a", 1<<20+1)
		req, err := http.NewRequest(http.MethodPost, "/audit/event", strings.NewReader(largeBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(eventHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusRequestEntityTooLarge {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusRequestEntityTooLarge)
		}
	})
}
