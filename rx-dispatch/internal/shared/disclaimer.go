package shared

import (
	"crypto/sha256"
	"encoding/hex"
)

// Mandatory Legal & Clinical Constants as mandated by RX Dispatch specification
const (
	ReadingTypeGenericAutomated      = "GENERIC_AUTOMATED"
	StatusWithoutMedicalSignature    = "WITHOUT_MEDICAL_SIGNATURE"
	MedicalReportNotIncluded         = "NOT_INCLUDED"
	PhysicianSignatureMustBeRequired = "MUST_BE_REQUESTED"

	// Mandatory Legal Disclaimer Text (Verbatim as required)
	MandatoryLegalDisclaimer = `La información presentada corresponde a una lectura genérica automatizada de la imagen y no constituye un diagnóstico médico ni sustituye un informe radiológico oficial. Si requiere el informe y la firma del médico responsable, deberá solicitarlo directamente al servicio médico correspondiente.`

	// Detailed Notice for Generated Results
	MandatoryResultNotice = `AVISO IMPORTANTE:
La lectura incluida en este resultado es de carácter GENÉRICO y tiene como finalidad proporcionar una descripción general de la imagen.
No constituye un diagnóstico médico ni sustituye el informe radiológico oficial.
Si requiere el informe y la firma del médico responsable, deberá solicitarlo directamente al servicio correspondiente.`
)

// CalculateSHA256 returns the hexadecimal SHA-256 hash of a byte slice.
func CalculateSHA256(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// CalculateStringSHA256 returns the hexadecimal SHA-256 hash of a string.
func CalculateStringSHA256(s string) string {
	return CalculateSHA256([]byte(s))
}
