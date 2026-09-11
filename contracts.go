package contracts

// This file implements the Data Transfer Objects (DTOs) as defined in
// the authoritative CONTRACTS.md document. These structs are used for
// serialization and deserialization of JSON payloads for inter-service
// communication.

// RequestGenericReadingRequest is the contract for the POST /reader/analyze endpoint.
// It carries the necessary information for the rx-reader service to perform a generic reading.
type RequestGenericReadingRequest struct {
	StudyID          string `json:"studyId"`
	StudyType        string `json:"studyType"`
	AnatomicalRegion string `json:"anatomicalRegion"`
}

// AuditResource contains correlated identity vectors for an audit event.
// It is a nested structure within RecordAuditEventRequest, as defined in SECURITY.md.
type AuditResource struct {
	CallingAETitle  string `json:"callingAeTitle"`
	AssociationID   string `json:"associationId"`
	SourceIPAddress string `json:"sourceIpAddress"`
}

// RecordAuditEventRequest is the contract for the POST /audit/record endpoint.
// It is used by all services to report critical events to the rx-audit service.
type RecordAuditEventRequest struct {
	Actor     string        `json:"actor"`
	Timestamp string        `json:"timestamp"` // ISO-8601 format
	Event     string        `json:"event"`
	Resource  AuditResource `json:"resource"`
	Result    string        `json:"result"` // e.g., "SUCCESS" / "FAILURE"
}

// HealthResponse is the common contract for the GET /healthz endpoint.
// It provides a standardized health status report for any given service.
type HealthResponse struct {
	Service   string `json:"service"`
	Port      int    `json:"port"`
	Status    string `json:"status"`    // e.g., "UP", "DOWN"
	Timestamp string `json:"timestamp"` // ISO-8601 format
	Version   string `json:"version"`
}