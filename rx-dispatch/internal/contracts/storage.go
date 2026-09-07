package contracts

import "time"

// StoreArtifactRequest saves a raw byte payload or encoded data into RX Storage.
type StoreArtifactRequest struct {
	ArtifactID   string `json:"artifactId"`
	ArtifactType string `json:"artifactType"` // "ORIGINAL_IMAGE", "DERIVATIVE_IMAGE", "RESULT_PDF"
	StudyID      string `json:"studyId"`
	FileName     string `json:"fileName"`
	Data         string `json:"data"` // Base64 encoded or content string
	ExpectedSHA  string `json:"expectedSha,omitempty"`
}

// StoreArtifactResponse confirms storage and returns storage URI and verified hash.
type StoreArtifactResponse struct {
	ArtifactID  string    `json:"artifactId"`
	StorageURI  string    `json:"storageUri"`
	ChecksumSHA string    `json:"checksumSha"`
	SizeBytes   int       `json:"sizeBytes"`
	StoredAt    time.Time `json:"storedAt"`
	IsImmutable bool      `json:"isImmutable"`
}

// VerifyArtifactRequest verifies the checksum and immutability of a stored artifact.
type VerifyArtifactRequest struct {
	ArtifactID string `json:"artifactId"`
}

// VerifyArtifactResponse reports artifact status and integrity.
type VerifyArtifactResponse struct {
	ArtifactID  string `json:"artifactId"`
	Exists      bool   `json:"exists"`
	ChecksumSHA string `json:"checksumSha"`
	MatchesHash bool   `json:"matchesHash"`
	IsImmutable bool   `json:"isImmutable"`
}

// GetArtifactResponse returns artifact content and metadata.
type GetArtifactResponse struct {
	ArtifactID   string `json:"artifactId"`
	ArtifactType string `json:"artifactType"`
	StudyID      string `json:"studyId"`
	FileName     string `json:"fileName"`
	Data         string `json:"data"`
	ChecksumSHA  string `json:"checksumSha"`
}
