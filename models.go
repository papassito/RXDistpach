package models

// This file defines the core domain models for the RX DISPATCH system.
// Unlike DTOs in the 'contracts' package, these models represent the
// internal business logic entities.

// Disclaimer represents the set of mandatory legal and clinical disclaimers
// that must accompany any automated reading. This structure enforces the
// clinical invariant defined in MANIFEST.md.
type Disclaimer struct {
	ReadingType               string `json:"readingType"`               // e.g., "GENERIC_AUTOMATED"
	SignatureStatus           string `json:"signatureStatus"`           // e.g., "WITHOUT_MEDICAL_SIGNATURE"
	MedicalReportStatus       string `json:"medicalReportStatus"`       // e.g., "NOT_INCLUDED"
	FullMandatoryLegalWarning string `json:"fullMandatoryLegalWarning"` // The full legal text.
}

// GenericReading is the core domain model representing the output of the
// rx-reader service. It encapsulates the automated analysis and the
// mandatory, non-negotiable disclaimers.
type GenericReading struct {
	StudyID          string     `json:"studyId"`
	AnatomicalRegion string     `json:"anatomicalRegion"`
	Timestamp        string     `json:"timestamp"` // ISO-8601 format
	ReadingContent   string     `json:"readingContent"`
	Disclaimer       Disclaimer `json:"disclaimer"`
}