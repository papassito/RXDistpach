package contracts

import (
	"rx-dispatch/internal/models"
)

// BuildResultRequest requests assembly of the dispatch packet.
type BuildResultRequest struct {
	StudyID string                `json:"studyId"`
	Study   models.Study          `json:"study"`
	Reading models.GenericReading `json:"reading"`
}

// BuildResultResponse returns the assembled ResultPackage.
type BuildResultResponse struct {
	ResultPackage models.ResultPackage `json:"resultPackage"`
	PackageID     string               `json:"packageId"`
	ChecksumSHA   string               `json:"checksumSha"`
	FileManifest  []string             `json:"fileManifest"`
}
