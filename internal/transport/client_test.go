package transport

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"rx-dispatch/internal/contracts"
)

func TestAuditClient_SendEvent(t *testing.T) {
	t.Run("should send event successfully on 2xx response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/audit/event" {
				t.Errorf("Expected to request '/audit/event', got %s", r.URL.Path)
			}
			w.WriteHeader(http.StatusAccepted)
		}))
		defer server.Close()

		os.Setenv("RX_AUDIT_ENDPOINT", server.URL)
		defer os.Unsetenv("RX_AUDIT_ENDPOINT")

		client := NewAuditClient()
		err := client.SendEvent(contracts.RecordAuditEventRequest{})

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("should return error on non-2xx response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		os.Setenv("RX_AUDIT_ENDPOINT", server.URL)
		defer os.Unsetenv("RX_AUDIT_ENDPOINT")

		client := NewAuditClient()
		err := client.SendEvent(contracts.RecordAuditEventRequest{})

		if err == nil {
			t.Error("Expected an error, got nil")
		}
	})

	t.Run("should return error on timeout", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond) // Sleep longer than client timeout
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		os.Setenv("RX_AUDIT_ENDPOINT", server.URL)
		defer os.Unsetenv("RX_AUDIT_ENDPOINT")

		client := NewAuditClient()
		client.client.Timeout = 50 * time.Millisecond // Set a short timeout

		err := client.SendEvent(contracts.RecordAuditEventRequest{})

		if err == nil {
			t.Error("Expected a timeout error, got nil")
		}
	})
}
