package contracts

import (
	"rx-dispatch/internal/models"
)

// IngestStudyRequest is sent to RX Study to register a new examination.
type IngestStudyRequest struct {
	StudyIdentifier    string         `json:"studyIdentifier"`
	StudyType          string         `json:"studyType"`
	AnatomicalRegion   string         `json:"anatomicalRegion"`
	ReferringPhysician string         `json:"referringPhysician"`
	ClinicalIndication string         `json:"clinicalIndication"`
	Patient            models.Patient `json:"patient"`
}

// IngestStudyResponse returns the registered study entity.
type IngestStudyResponse struct {
	Study models.Study `json:"study"`
}

// UpdateStudyStatusRequest updates the lifecycle stage of a study.
type UpdateStudyStatusRequest struct {
	Status models.StudyStatus `json:"status"`
	Reason string             `json:"reason,omitempty"`
}

// AttachImageToStudyRequest adds an image reference to a study.
type AttachImageToStudyRequest struct {
	Image models.XRayImage `json:"image"`
}

// ListStudiesResponse returns all active or filtered studies.
type ListStudiesResponse struct {
	Studies []models.Study `json:"studies"`
	Count   int            `json:"count"`
}
