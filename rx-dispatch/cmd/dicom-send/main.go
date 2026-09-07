package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"rx-dispatch/internal/dicom"
)

type SendResult struct {
	Success          bool   `json:"success"`
	CallingAET       string `json:"callingAet"`
	CalledAET        string `json:"calledAet"`
	Host             string `json:"host"`
	Port             int    `json:"port"`
	PatientName      string `json:"patientName"`
	PatientID        string `json:"patientId"`
	StudyDescription string `json:"studyDescription"`
	AccessionNumber  string `json:"accessionNumber"`
	BytesSent        int    `json:"bytesSent"`
	Error            string `json:"error,omitempty"`
}

func main() {
	host := flag.String("host", "127.0.0.1", "DICOM SCP Host")
	port := flag.Int("port", 11112, "DICOM SCP Port")
	callingAET := flag.String("calling-aet", "SIEMENS_YSIO_DIGITAL", "Calling AE Title (Modality)")
	calledAET := flag.String("called-aet", "RX_DISPATCH", "Called AE Title (Server)")
	patientName := flag.String("patient-name", "PEREZ^ALEJANDRO", "Patient Name (DICOM PN format)")
	patientID := flag.String("patient-id", "PAC-DCM-9941", "Patient ID")
	studyDesc := flag.String("study-desc", "RX TORAX AP DIGITAL", "Study Description")
	accession := flag.String("accession", "ACC-2026-DX01", "Accession Number")
	asJSON := flag.Bool("json", true, "Output result as JSON")
	flag.Parse()

	dicomBytes := dicom.GenerateSyntheticDICOM(
		*patientName,
		*patientID,
		*accession,
		*studyDesc,
	)

	err := dicom.SendCStore(*host, *port, *callingAET, *calledAET, dicomBytes)

	res := SendResult{
		Success:          err == nil,
		CallingAET:       *callingAET,
		CalledAET:        *calledAET,
		Host:             *host,
		Port:             *port,
		PatientName:      *patientName,
		PatientID:        *patientID,
		StudyDescription: *studyDesc,
		AccessionNumber:  *accession,
		BytesSent:        len(dicomBytes),
	}

	if err != nil {
		res.Error = err.Error()
	}

	if *asJSON {
		_ = json.NewEncoder(os.Stdout).Encode(res)
	} else {
		if err != nil {
			fmt.Fprintf(os.Stderr, "DICOM C-STORE failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("DICOM C-STORE OK: %d bytes sent from %s to %s@%s:%d\n", len(dicomBytes), *callingAET, *calledAET, *host, *port)
	}

	if err != nil {
		os.Exit(1)
	}
}
