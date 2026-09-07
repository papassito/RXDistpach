package dicom

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// Dataset holds clinically relevant tags parsed from an incoming DICOM instance.
type Dataset struct {
	PatientName       string    `json:"patientName"`
	PatientID         string    `json:"patientId"`
	PatientBirthDate  string    `json:"patientBirthDate"`
	PatientSex        string    `json:"patientSex"`
	Modality          string    `json:"modality"`
	StudyDescription  string    `json:"studyDescription"`
	StudyDate         string    `json:"studyDate"`
	AccessionNumber   string    `json:"accessionNumber"`
	StudyInstanceUID  string    `json:"studyInstanceUid"`
	SeriesInstanceUID string    `json:"seriesInstanceUid"`
	SOPInstanceUID    string    `json:"sopInstanceUid"`
	BodyPartExamined  string    `json:"bodyPartExamined"`
	Manufacturer      string    `json:"manufacturer"`
	InstitutionName   string    `json:"institutionName"`
	Rows              int       `json:"rows"`
	Columns           int       `json:"columns"`
	ChecksumSHA       string    `json:"checksumSha"`
	SizeBytes         int       `json:"sizeBytes"`
	ReceivedAt        time.Time `json:"receivedAt"`
	RawBytes          []byte    `json:"-"`
}

// ParseDICOMBytes inspects a raw byte array and extracts core DICOM tags.
func ParseDICOMBytes(data []byte) (*Dataset, error) {
	if len(data) < 132 {
		return nil, fmt.Errorf("dicom payload too short (%d bytes)", len(data))
	}

	h := sha256.Sum256(data)
	checksum := hex.EncodeToString(h[:])

	ds := &Dataset{
		ReceivedAt:  time.Now(),
		SizeBytes:   len(data),
		ChecksumSHA: checksum,
		RawBytes:    data,
		Modality:    "DX",
	}

	// Check standard 128-byte preamble + "DICM"
	hasMagic := string(data[128:132]) == "DICM"
	offset := 0
	if hasMagic {
		offset = 132
	}

	// Simple parser for standard Little Endian Explicit/Implicit elements
	for offset+8 <= len(data) {
		group := binary.LittleEndian.Uint16(data[offset : offset+2])
		element := binary.LittleEndian.Uint16(data[offset+2 : offset+4])

		var length int
		var valueOffset int

		// Check if explicit VR (two uppercase ASCII letters)
		vr := string(data[offset+4 : offset+6])
		if isStandardVR(vr) {
			if vr == "OB" || vr == "OW" || vr == "OF" || vr == "SQ" || vr == "UT" || vr == "UN" {
				if offset+12 > len(data) {
					break
				}
				length = int(binary.LittleEndian.Uint32(data[offset+8 : offset+12]))
				valueOffset = offset + 12
			} else {
				length = int(binary.LittleEndian.Uint16(data[offset+6 : offset+8]))
				valueOffset = offset + 8
			}
		} else {
			// Implicit VR: 32-bit length directly after tag
			length = int(binary.LittleEndian.Uint32(data[offset+4 : offset+8]))
			valueOffset = offset + 8
		}

		if length < 0 || valueOffset+length > len(data) {
			break
		}

		val := cleanString(data[valueOffset : valueOffset+length])

		// Map key DICOM clinical tags
		switch {
		case group == 0x0010 && element == 0x0010: // PatientName
			ds.PatientName = val
		case group == 0x0010 && element == 0x0020: // PatientID
			ds.PatientID = val
		case group == 0x0010 && element == 0x0030: // PatientBirthDate
			ds.PatientBirthDate = val
		case group == 0x0010 && element == 0x0040: // PatientSex
			ds.PatientSex = val
		case group == 0x0008 && element == 0x0060: // Modality
			ds.Modality = val
		case group == 0x0008 && element == 0x1030: // StudyDescription
			ds.StudyDescription = val
		case group == 0x0008 && element == 0x0020: // StudyDate
			ds.StudyDate = val
		case group == 0x0008 && element == 0x0050: // AccessionNumber
			ds.AccessionNumber = val
		case group == 0x0020 && element == 0x000D: // StudyInstanceUID
			ds.StudyInstanceUID = val
		case group == 0x0020 && element == 0x000E: // SeriesInstanceUID
			ds.SeriesInstanceUID = val
		case group == 0x0008 && element == 0x0018: // SOPInstanceUID
			ds.SOPInstanceUID = val
		case group == 0x0018 && element == 0x0015: // BodyPartExamined
			ds.BodyPartExamined = val
		case group == 0x0008 && element == 0x0070: // Manufacturer
			ds.Manufacturer = val
		case group == 0x0008 && element == 0x0080: // InstitutionName
			ds.InstitutionName = val
		case group == 0x0028 && element == 0x0010: // Rows
			if len(data[valueOffset:valueOffset+length]) >= 2 {
				ds.Rows = int(binary.LittleEndian.Uint16(data[valueOffset : valueOffset+2]))
			}
		case group == 0x0028 && element == 0x0011: // Columns
			if len(data[valueOffset:valueOffset+length]) >= 2 {
				ds.Columns = int(binary.LittleEndian.Uint16(data[valueOffset : valueOffset+2]))
			}
		}

		offset = valueOffset + length
		// DICOM tags are always aligned to even byte boundaries
		if offset%2 != 0 {
			offset++
		}
	}

	// Fallback sensible defaults if metadata is minimal
	if ds.PatientName == "" {
		ds.PatientName = "PACIENTE ANONIMIZADO"
	}
	if ds.PatientID == "" {
		ds.PatientID = fmt.Sprintf("PAC-%s", checksum[:8])
	}
	if ds.StudyDescription == "" {
		ds.StudyDescription = "ESTUDIO RADIOGRAFICO DIGITAL"
	}
	if ds.StudyInstanceUID == "" {
		ds.StudyInstanceUID = fmt.Sprintf("1.2.840.10008.%s", checksum[:16])
	}

	return ds, nil
}

func isStandardVR(s string) bool {
	vrs := map[string]bool{
		"AE": true, "AS": true, "AT": true, "CS": true, "DA": true, "DS": true,
		"DT": true, "FL": true, "FD": true, "IS": true, "LO": true, "LT": true,
		"OB": true, "OD": true, "OF": true, "OL": true, "OW": true, "PN": true,
		"SH": true, "SL": true, "SQ": true, "SS": true, "ST": true, "TM": true,
		"UI": true, "UL": true, "UN": true, "US": true, "UT": true,
	}
	return vrs[s]
}

func cleanString(b []byte) string {
	s := strings.Trim(string(b), "\x00 \t\r\n")
	return s
}
