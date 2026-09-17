package models

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type Disclaimer struct {
	Text                      string `json:"text"`
	ReadingType               string `json:"reading_type"`
	SignatureStatus           string `json:"signature_status"`
	MedicalReportStatus       string `json:"medical_report_status"`
	FullMandatoryLegalWarning string `json:"full_mandatory_legal_warning"`
}

type GenericReading struct {
	StudyID          string     `json:"study_id"`
	Status           string     `json:"status"`
	Findings         string     `json:"findings"`
	AnatomicalRegion string     `json:"anatomical_region"`
	Timestamp        string     `json:"timestamp"`
	ReadingContent   string     `json:"reading_content"`
	Disclaimer       Disclaimer `json:"disclaimer"`
}

type AuditEventEPHI struct {
	EventTimestamp     time.Time `json:"event_timestamp"`
	EventAction        string    `json:"event_action"`
	EventOutcome       string    `json:"event_outcome"`
	EventErrorCode     string    `json:"event_error_code"`
	UserID             string    `json:"user_id"`
	SourceIP           string    `json:"source_ip"`
	SourceAETitle      string    `json:"source_ae_title"`
	DestinationIP      string    `json:"destination_ip"`
	DestinationAETitle string    `json:"destination_ae_title"`
	PatientID          string    `json:"patient_id"`
	StudyInstanceUID   string    `json:"study_instance_uid"`
	AccessionNumber    string    `json:"accession_number"`
	NumberOfInstances  int       `json:"number_of_instances"`
	SecurityTLSStatus  string    `json:"security_tls_status"`
}

type Resource struct {
	CallingAETitle  string `json:"callingAeTitle"`
	AssociationID   string `json:"associationId"`
	SourceIPAddress string `json:"sourceIpAddress"`
}

// HashPayload computes the SHA-256 hash of a byte slice and returns it as a hex string.
// This is used to ensure ePHI integrity as per SECURITY.md.
func HashPayload(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
