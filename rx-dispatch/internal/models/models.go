package models

import "time"

// StudyStatus represents the lifecycle state of a study in RX Dispatch.
type StudyStatus string

const (
	StudyStatusReceived     StudyStatus = "RECEIVED"
	StudyStatusImagesLoaded StudyStatus = "IMAGES_LOADED"
	StudyStatusProcessing   StudyStatus = "PROCESSING"
	StudyStatusGenericRead  StudyStatus = "GENERIC_READ_COMPLETED"
	StudyStatusPackageBuilt StudyStatus = "PACKAGE_BUILT"
	StudyStatusDelivered    StudyStatus = "DELIVERED"
	StudyStatusFailed       StudyStatus = "FAILED"
)

// DeliveryStatus represents the state of a dispatch delivery.
type DeliveryStatus string

const (
	DeliveryPending    DeliveryStatus = "PENDING"
	DeliveryProcessing DeliveryStatus = "PROCESSING"
	DeliverySent       DeliveryStatus = "SENT"
	DeliveryFailed     DeliveryStatus = "FAILED"
)

// Patient represents patient demographic data associated with the study.
type Patient struct {
	ID        string `json:"id"`
	FullName  string `json:"fullName"`
	BirthDate string `json:"birthDate"`
	Gender    string `json:"gender"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

// ImageMetadata stores acquisition and technical parameters of the radiograph.
type ImageMetadata struct {
	Format      string `json:"format"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	BitDepth    int    `json:"bitDepth"`
	Kvp         string `json:"kvp"`
	MilliAmps   string `json:"milliAmps"`
	Projection  string `json:"projection"`
	ChecksumSHA string `json:"checksumSHA"`
}

// XRayImage represents an immutable radiographic asset and its derivative.
type XRayImage struct {
	ID            string        `json:"id"`
	StudyID       string        `json:"studyId"`
	FileName      string        `json:"fileName"`
	OriginalURI   string        `json:"originalUri"`
	DerivedURI    string        `json:"derivedUri"`
	CreatedAt     time.Time     `json:"createdAt"`
	Metadata      ImageMetadata `json:"metadata"`
	IsIntegrityOK bool          `json:"isIntegrityOk"`
}

// Study represents the core radiological study container.
type Study struct {
	ID                 string      `json:"id"`
	StudyIdentifier    string      `json:"studyIdentifier"`
	Date               time.Time   `json:"date"`
	StudyType          string      `json:"studyType"`
	AnatomicalRegion   string      `json:"anatomicalRegion"`
	ReferringPhysician string      `json:"referringPhysician"`
	ClinicalIndication string      `json:"clinicalIndication"`
	Status             StudyStatus `json:"status"`
	Patient            Patient     `json:"patient"`
	Images             []XRayImage `json:"images"`
}

// GenericReading represents the strictly generic, non-diagnostic observation output.
type GenericReading struct {
	ID                        string    `json:"id"`
	StudyID                   string    `json:"studyId"`
	CreatedAt                 time.Time `json:"createdAt"`
	StudyTitle                string    `json:"studyTitle"`
	GeneralDescription        string    `json:"generalDescription"`
	VisualObservations        []string  `json:"visualObservations"`
	ObservableCharacteristics []string  `json:"observableCharacteristics"`
	TechnicalCaveats          string    `json:"technicalCaveats"`
	ReadingType               string    `json:"readingType"`              // GENERIC_AUTOMATED
	MedicalReportNotice       string    `json:"medicalReportNotice"`      // NOT_INCLUDED
	PhysicianSignatureStatus  string    `json:"physicianSignatureStatus"` // WITHOUT_MEDICAL_SIGNATURE
	MandatoryDisclaimer       string    `json:"mandatoryDisclaimer"`
}

// ResultPackage represents the assembled dispatch packet delivered to the recipient.
type ResultPackage struct {
	ID                       string         `json:"id"`
	StudyID                  string         `json:"studyId"`
	GeneratedAt              time.Time      `json:"generatedAt"`
	Study                    Study          `json:"study"`
	Reading                  GenericReading `json:"reading"`
	ResultType               string         `json:"resultType"`
	StatusNotice             string         `json:"statusNotice"`
	OfficialReportNotice     string         `json:"officialReportNotice"`
	PhysicianSignatureNotice string         `json:"physicianSignatureNotice"`
	MandatoryLegend          string         `json:"mandatoryLegend"`
	RequestInstructions      string         `json:"requestInstructions"`
	PackageFileNames         []string       `json:"packageFileNames"`
	ChecksumSHA              string         `json:"checksumSha"`
}

// DeliveryRecord tracks an external dispatch transaction.
type DeliveryRecord struct {
	ID            string         `json:"id"`
	PackageID     string         `json:"packageId"`
	StudyID       string         `json:"studyId"`
	Recipient     string         `json:"recipient"`
	Channel       string         `json:"channel"`
	Status        DeliveryStatus `json:"status"`
	SentAt        time.Time      `json:"sentAt"`
	TrackingToken string         `json:"trackingToken"`
	FailureReason string         `json:"failureReason,omitempty"`
}

// AuditEvent represents an immutable operational log entry.
type AuditEvent struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	EventType string    `json:"eventType"`
	StudyID   string    `json:"studyId"`
	Service   string    `json:"service"`
	Actor     string    `json:"actor"`
	Details   string    `json:"details"`
}
