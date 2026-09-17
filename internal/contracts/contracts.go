package contracts

import "rx-dispatch/internal/models"

// RequestGenericReadingRequest is the contract for the POST /reader/analyze endpoint.
// It carries the necessary information for the rx-reader service to perform a generic reading.
type RequestGenericReadingRequest struct {
	StudyID          string `json:"study_id"`
	StudyType        string `json:"study_type"`
	AnatomicalRegion string `json:"anatomical_region"`
}

// RecordAuditEventRequest is the contract for the POST /audit/record endpoint.
// It wraps the canonical AuditEvent model defined in the models package.
type RecordAuditEventRequest struct {
	Event models.AuditEvent `json:"event"`
}

type HealthResponse struct {
	Service   string `json:"service"`
	Port      int    `json:"port"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

type ConsolidateReadingRequest struct {
	StudyID string `json:"study_id"`
	Details string `json:"details"`
}