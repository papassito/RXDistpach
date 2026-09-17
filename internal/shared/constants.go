package shared

// This file contains normative constants that enforce the clinical and legal
// invariants of the RX DISPATCH system, as defined in MANIFEST.md.

const (
	// ReadingTypeGenericAutomated is the mandatory type for all non-diagnostic readings.
	ReadingTypeGenericAutomated = "GENERIC_AUTOMATED"

	// StatusWithoutMedicalSignature indicates the reading has not been reviewed by a physician.
	StatusWithoutMedicalSignature = "WITHOUT_MEDICAL_SIGNATURE"

	// MedicalReportNotIncluded clarifies that this is not an official report.
	MedicalReportNotIncluded = "NOT_INCLUDED"

	// MandatoryLegalDisclaimer is the full legal text required by REQ-RDR-002.
	MandatoryLegalDisclaimer = "La información presentada corresponde a una lectura genérica automatizada de la imagen y no constituye un diagnóstico médico ni sustituye un informe radiológico oficial. Si requiere el informe y la firma del médico responsable, deberá solicitarlo directamente al servicio médico correspondiente."

	// DicomAssociateRJ is the normative protocol command for rejecting a DICOM association, as per SECURITY.md.
	DicomAssociateRJ = "A-ASSOCIATE-RJ"
)
