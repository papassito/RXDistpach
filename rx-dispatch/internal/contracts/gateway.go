package contracts

import (
	"rx-dispatch/internal/models"
)

// IngressStudySubmission represents a client payload submitting a full radiological examination.
type IngressStudySubmission struct {
	StudyIdentifier    string         `json:"studyIdentifier"`
	StudyType          string         `json:"studyType"`
	AnatomicalRegion   string         `json:"anatomicalRegion"`
	ReferringPhysician string         `json:"referringPhysician"`
	ClinicalIndication string         `json:"clinicalIndication"`
	Patient            models.Patient `json:"patient"`
	Images             []struct {
		FileName   string `json:"fileName"`
		Data       string `json:"data"`
		Projection string `json:"projection"`
		Format     string `json:"format"`
	} `json:"images"`
}

// OrchestratedDispatchFlowRequest triggers an automated end-to-end execution through all 9 services.
type OrchestratedDispatchFlowRequest struct {
	StudyID         string `json:"studyId"`
	DeliveryChannel string `json:"deliveryChannel"`
	DeliveryTarget  string `json:"deliveryTarget"`
}

// OrchestratedDispatchFlowResponse returns the status of each step across the 9 microservices.
type OrchestratedDispatchFlowResponse struct {
	StudyID         string                `json:"studyId"`
	SecurityChecked bool                  `json:"securityChecked"`
	StudyStatus     models.StudyStatus    `json:"studyStatus"`
	ImageIntegrity  bool                  `json:"imageIntegrity"`
	ReadingType     string                `json:"readingType"`
	PackageChecksum string                `json:"packageChecksum"`
	DeliveryToken   string                `json:"deliveryToken"`
	DeliveryStatus  models.DeliveryStatus `json:"deliveryStatus"`
	AuditCount      int                   `json:"auditCount"`
}

// GatewayTopologyStatus returns the live health of all 9 independent microservices.
type GatewayTopologyStatus struct {
	GatewayService  HealthResponse `json:"gatewayService"`
	SecurityService HealthResponse `json:"securityService"`
	StudyService    HealthResponse `json:"studyService"`
	StorageService  HealthResponse `json:"storageService"`
	ImageService    HealthResponse `json:"imageService"`
	ReaderService   HealthResponse `json:"readerService"`
	ResultService   HealthResponse `json:"resultService"`
	DeliveryService HealthResponse `json:"deliveryService"`
	AuditService    HealthResponse `json:"auditService"`
}
