package contracts

import (
	"rx-dispatch/internal/models"
)

// RecordAuditEventRequest records an operational or security event.
type RecordAuditEventRequest struct {
	EventType string `json:"eventType"`
	StudyID   string `json:"studyId"`
	Service   string `json:"service"`
	Actor     string `json:"actor"`
	Details   string `json:"details"`
}

// RecordAuditEventResponse acknowledges audit persistence.
type RecordAuditEventResponse struct {
	EventID string `json:"eventId"`
	Success bool   `json:"success"`
}

// QueryAuditEventsResponse returns matching audit trail records.
type QueryAuditEventsResponse struct {
	Events []models.AuditEvent `json:"events"`
	Count  int                 `json:"count"`
}
