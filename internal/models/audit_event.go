package models

import "time"

type AuditEvent struct {
    ID                 string    `json:"id" db:"id"`
    EventTimestamp     time.Time `json:"event_timestamp" db:"event_timestamp"`
    EventAction        string    `json:"event_action" db:"event_action"`
    EventOutcome       string    `json:"event_outcome" db:"event_outcome"`
    EventErrorCode     string    `json:"event_error_code,omitempty" db:"event_error_code"`
    UserID             string    `json:"user_id" db:"user_id"`
    SourceIP           string    `json:"source_ip" db:"source_ip"`
    SourceAETitle      string    `json:"source_ae_title,omitempty" db:"source_ae_title"`
    DestinationIP      string    `json:"destination_ip" db:"destination_ip"`
    DestinationAETitle string    `json:"destination_ae_title,omitempty" db:"destination_ae_title"`
    PatientID          string    `json:"patient_id,omitempty" db:"patient_id"`
    StudyInstanceUID   string    `json:"study_instance_uid,omitempty" db:"study_instance_uid"`
    AccessionNumber    string    `json:"accession_number,omitempty" db:"accession_number"`
    NumberOfInstances  int       `json:"number_of_instances,omitempty" db:"number_of_instances"`
    SecurityTLSStatus  bool      `json:"security_tls_status" db:"security_tls_status"`
    CreatedAt          time.Time `json:"created_at" db:"created_at"`
}