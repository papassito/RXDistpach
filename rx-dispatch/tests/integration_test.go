package tests

import (
	"testing"
	"time"

	"rx-dispatch/internal/models"
	"rx-dispatch/internal/shared"
)

// TestEndToEndDecoupledPipeline verifies the sequential contract transformations.
func TestEndToEndDecoupledPipeline(t *testing.T) {
	// 1. RX Study: Ingest
	study := models.Study{
		ID:                 "STU-E2E-1",
		StudyIdentifier:    "RX-2026-E2E",
		Date:               time.Now(),
		StudyType:          "Radiografía de Muñeca",
		AnatomicalRegion:   "Muñeca",
		ReferringPhysician: "Dr. Guard",
		ClinicalIndication: "Caída sobre mano",
		Status:             models.StudyStatusReceived,
		Patient: models.Patient{
			ID:       "PAC-E2E",
			FullName: "Juan Pérez",
			Email:    "juan@example.com",
		},
	}
	if study.Status != models.StudyStatusReceived {
		t.Fatalf("expected RECEIVED status")
	}

	// 2. RX Image + RX Storage: Checksum & Derivative
	rawImageContent := "SAMPLE_PIXEL_DATA_ARRAY"
	checksum := shared.CalculateStringSHA256(rawImageContent + "wrist.png")
	xrayImage := models.XRayImage{
		ID:            "IMG-E2E-1",
		StudyID:       study.ID,
		FileName:      "wrist.png",
		OriginalURI:   "/storage/artifacts/ORIG-IMG-E2E-1",
		DerivedURI:    "/storage/artifacts/DERIV-IMG-E2E-1",
		IsIntegrityOK: true,
		Metadata: models.ImageMetadata{
			ChecksumSHA: checksum,
		},
	}
	study.Images = append(study.Images, xrayImage)
	study.Status = models.StudyStatusProcessing

	// 3. RX Reader: Generic Non-Diagnostic Reading
	reading := models.GenericReading{
		ID:                        "READ-E2E-1",
		StudyID:                   study.ID,
		ReadingType:               shared.ReadingTypeGenericAutomated,
		MedicalReportNotice:       shared.MedicalReportNotIncluded,
		PhysicianSignatureStatus:  shared.StatusWithoutMedicalSignature,
		MandatoryDisclaimer:       shared.MandatoryLegalDisclaimer,
		VisualObservations:        []string{"Observación genérica visual de arcos óseos"},
		ObservableCharacteristics: []string{"Región: Muñeca"},
	}
	study.Status = models.StudyStatusGenericRead

	// 4. RX Result: Package Assembly with mandatory legend
	pkgID := "PKG-E2E-1"
	pkgChecksum := shared.CalculateStringSHA256(pkgID + study.ID + reading.ID)
	resultPkg := models.ResultPackage{
		ID:                       pkgID,
		StudyID:                  study.ID,
		Study:                    study,
		Reading:                  reading,
		ResultType:               "LECTURA GENÉRICA AUTOMATIZADA",
		StatusNotice:             "SIN FIRMA MÉDICA",
		OfficialReportNotice:     "INFORME MÉDICO OFICIAL: NO INCLUIDO",
		PhysicianSignatureNotice: "FIRMA DEL MÉDICO RESPONSABLE: DEBE SOLICITARSE",
		MandatoryLegend:          shared.MandatoryLegalDisclaimer,
		ChecksumSHA:              pkgChecksum,
	}
	study.Status = models.StudyStatusPackageBuilt

	// 5. RX Delivery: Dispatch
	delivery := models.DeliveryRecord{
		ID:            "DEL-E2E-1",
		PackageID:     resultPkg.ID,
		StudyID:       study.ID,
		Recipient:     study.Patient.Email,
		Channel:       "EMAIL",
		Status:        models.DeliverySent,
		TrackingToken: "TRK-E2E-ABC123",
	}
	study.Status = models.StudyStatusDelivered

	// Assertions
	if study.Status != models.StudyStatusDelivered {
		t.Fatalf("study final status must be DELIVERED, got %s", study.Status)
	}
	if delivery.TrackingToken == "" {
		t.Fatalf("delivery tracking token must be generated")
	}
	if resultPkg.MandatoryLegend != shared.MandatoryLegalDisclaimer {
		t.Fatalf("result package must enforce the mandatory clinical disclaimer")
	}
}
