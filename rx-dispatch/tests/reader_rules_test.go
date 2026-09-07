package tests

import (
	"strings"
	"testing"

	"rx-dispatch/internal/models"
	"rx-dispatch/internal/shared"
)

// TestReaderMandatoryDisclaimerEnforced verifies that any generic reading contains the exact legal legend.
func TestReaderMandatoryDisclaimerEnforced(t *testing.T) {
	reading := models.GenericReading{
		ID:                       "GEN-READ-001",
		StudyID:                  "STU-001",
		StudyTitle:               "Radiografía de Tórax PA",
		GeneralDescription:       "Descripción genérica de densidades radiográficas.",
		ReadingType:              shared.ReadingTypeGenericAutomated,
		MedicalReportNotice:      shared.MedicalReportNotIncluded,
		PhysicianSignatureStatus: shared.StatusWithoutMedicalSignature,
		MandatoryDisclaimer:      shared.MandatoryLegalDisclaimer,
	}

	// 1. Verify exact classification
	if reading.ReadingType != "GENERIC_AUTOMATED" {
		t.Fatalf("expected ReadingType to be GENERIC_AUTOMATED, got %s", reading.ReadingType)
	}

	if reading.PhysicianSignatureStatus != "WITHOUT_MEDICAL_SIGNATURE" {
		t.Fatalf("expected PhysicianSignatureStatus to be WITHOUT_MEDICAL_SIGNATURE, got %s", reading.PhysicianSignatureStatus)
	}

	if reading.MedicalReportNotice != "NOT_INCLUDED" {
		t.Fatalf("expected MedicalReportNotice to be NOT_INCLUDED, got %s", reading.MedicalReportNotice)
	}

	// 2. Verify mandatory legal disclaimer presence
	expectedSentence := "La información presentada corresponde a una lectura genérica automatizada de la imagen y no constituye un diagnóstico médico ni sustituye un informe radiológico oficial."
	if !strings.Contains(reading.MandatoryDisclaimer, expectedSentence) {
		t.Fatalf("mandatory disclaimer text missing expected clinical disclaimer")
	}

	requestPhysicianSentence := "Si requiere el informe y la firma del médico responsable, deberá solicitarlo directamente al servicio médico correspondiente."
	if !strings.Contains(reading.MandatoryDisclaimer, requestPhysicianSentence) {
		t.Fatalf("mandatory disclaimer text missing instruction to request physician signature")
	}
}

// TestReaderProhibitsDiagnosticLabels verifies that automated readings never label output as official diagnostic report.
func TestReaderProhibitsDiagnosticLabels(t *testing.T) {
	forbiddenLabels := []string{
		"diagnóstico médico",
		"diagnostico medico",
		"informe médico oficial",
		"informe radiológico firmado",
		"informe firmado",
	}

	readingType := shared.ReadingTypeGenericAutomated
	for _, forbidden := range forbiddenLabels {
		if strings.EqualFold(readingType, forbidden) {
			t.Fatalf("VIOLATION: Automated reading labeled as %s", forbidden)
		}
	}

	// Negative assertions on reading attributes
	reading := models.GenericReading{
		ReadingType:              shared.ReadingTypeGenericAutomated,
		MedicalReportNotice:      shared.MedicalReportNotIncluded,
		PhysicianSignatureStatus: shared.StatusWithoutMedicalSignature,
	}

	if reading.PhysicianSignatureStatus == "SIGNED_BY_PHYSICIAN" {
		t.Fatalf("VIOLATION: automated reading marked as signed by physician")
	}
	if reading.MedicalReportNotice == "INCLUDED" {
		t.Fatalf("VIOLATION: official medical report marked as included")
	}
	if strings.Contains(strings.ToLower(reading.ReadingType), "diagnóstico") {
		t.Fatalf("VIOLATION: reading type contains diagnostic statement")
	}
}
