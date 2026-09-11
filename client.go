package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"rx-dispatch/internal/contracts"
)

const (
	auditServiceURL = "http://127.0.0.1:8088/audit/record"
	clientTimeout   = 5 * time.Second
)

// AuditClient es un cliente dedicado para enviar eventos al servicio rx-audit.
type AuditClient struct {
	client  *http.Client
	baseURL string
}

// NewAuditClient crea un nuevo cliente para el servicio de auditoría.
func NewAuditClient() *AuditClient {
	return &AuditClient{
		client: &http.Client{
			Timeout: clientTimeout,
		},
		baseURL: auditServiceURL,
	}
}

// SendEvent envía un evento de auditoría estructurado.
// Está diseñado para ser "fire-and-forget"; registra errores internamente
// pero no los devuelve, para no bloquear la lógica de negocio principal del servicio que lo llama.
func (c *AuditClient) SendEvent(event contracts.RecordAuditEventRequest) {
	payload, err := json.Marshal(event)
	if err != nil {
		log.Printf("[AUDIT-CLIENT] ERROR: Failed to marshal audit event: %v", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL, bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("[AUDIT-CLIENT] ERROR: Failed to create audit request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		log.Printf("[AUDIT-CLIENT] ERROR: Failed to send audit event to %s: %v", c.baseURL, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		log.Printf("[AUDIT-CLIENT] WARN: Received non-success status %d from audit service", resp.StatusCode)
	} else {
		log.Printf("[AUDIT-CLIENT] Successfully sent event '%s' for actor '%s'", event.Event, event.Actor)
	}
}