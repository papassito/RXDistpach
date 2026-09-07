package contracts

import (
	"rx-dispatch/internal/models"
)

// RequestGenericReadingRequest is sent to RX Reader to analyze a derivative image.
type RequestGenericReadingRequest struct {
	StudyID          string   `json:"studyId"`
	StudyType        string   `json:"studyType"`
	AnatomicalRegion string   `json:"anatomicalRegion"`
	DerivativeURI    string   `json:"derivativeUri"`
	ImageData        string   `json:"imageData,omitempty"` // Derivative preview/base64
	Projections      []string `json:"projections"`
}

// RequestGenericReadingResponse returns the generic, non-diagnostic reading with mandatory notices.
type RequestGenericReadingResponse struct {
	Reading                  models.GenericReading `json:"reading"`
	ReadingType              string                `json:"readingType"`              // "GENERIC_AUTOMATED"
	PhysicianSignatureStatus string                `json:"physicianSignatureStatus"` // "WITHOUT_MEDICAL_SIGNATURE"
	MedicalReportNotice      string                `json:"medicalReportNotice"`      // "NOT_INCLUDED"
	MandatoryDisclaimer      string                `json:"mandatoryDisclaimer"`
	IsNonDiagnosticCertified bool                  `json:"isNonDiagnosticCertified"`
}
