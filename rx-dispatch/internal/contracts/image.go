package contracts

import (
	"rx-dispatch/internal/models"
)

// ProcessImageRequest requests RX Image to locate, validate, and prepare derivatives.
type ProcessImageRequest struct {
	StudyID    string `json:"studyId"`
	ImageID    string `json:"imageId"`
	FileName   string `json:"fileName"`
	RawData    string `json:"rawData"` // Original image data or URI
	Projection string `json:"projection"`
	Format     string `json:"format"`
}

// ProcessImageResponse returns the validated image entity with immutable original and derivative URIs.
type ProcessImageResponse struct {
	Image         models.XRayImage `json:"image"`
	OriginalURI   string           `json:"originalUri"`
	DerivativeURI string           `json:"derivativeUri"`
	ChecksumSHA   string           `json:"checksumSha"`
	IntegrityOK   bool             `json:"integrityOk"`
}

// ValidateIntegrityRequest asks RX Image to verify an existing image against its stored checksum.
type ValidateIntegrityRequest struct {
	ImageID     string `json:"imageId"`
	ExpectedSHA string `json:"expectedSha"`
}

// ValidateIntegrityResponse reports integrity check outcome.
type ValidateIntegrityResponse struct {
	ImageID     string `json:"imageId"`
	Calculated  string `json:"calculated"`
	Expected    string `json:"expected"`
	IntegrityOK bool   `json:"integrityOk"`
}
