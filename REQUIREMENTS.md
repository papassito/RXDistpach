# RX DISPATCH — System Requirements

**DOCUMENT STATUS:** `CONSOLIDATED BASELINE`
**PURPOSE:** Catalog of normative and functional system requirements.

---

## 1. `rx-reader` Requirements

*   **REQ-RDR-001 (Non-Diagnostic Output):** The system must provide a component (`rx-reader`) to generate a generic, automated reading that is explicitly not a medical diagnosis.
    *   **Evidence:** `IMPLEMENTED / STATICALLY VERIFIED`

*   **REQ-RDR-002 (Mandatory Disclaimer):** Every reading from `rx-reader` must programmatically inject the constants `ReadingTypeGenericAutomated`, `StatusWithoutMedicalSignature`, `MedicalReportNotIncluded`, and the full `MandatoryLegalDisclaimer` text.
    *   **Evidence:** `IMPLEMENTED / STATICALLY VERIFIED`

## 2. Audit & Traceability Requirements

*   **REQ-AUD-001 (Centralized Audit Service):** The system must have a central service (`rx-audit`) to log critical events.
    *   **Evidence:** `DEFINED`

*   **REQ-AUD-002 (Component Event Reporting):** Key components like `rx-reader` must send an audit event (`RecordAuditEventRequest`) to `rx-audit` after their main operation.
    *   **Evidence:** `IMPLEMENTED / STATICALLY VERIFIED`

## 3. Security & Access Requirements

*   **REQ-SEC-001 (Internal Service Isolation):** Communication between internal microservices must be restricted to the local loopback interface (`127.0.0.1`).
    *   **Evidence:** `DEFINED`

*   **REQ-SEC-002 (PxLab Read-Only Access):** Any system interaction with the PxLab database must be strictly `READ-ONLY`.
    *   **Evidence:** `DEFINED`

## 4. Dispatch & Delivery Requirements

*   **REQ-DEL-001 (Store-and-Forward Mechanism):** The delivery service (`rx-delivery`) must implement a resilient queue to handle and retry failed deliveries to external systems.
    *   **Evidence:** `DEFINED`

## 5. Infrastructure Requirements

*   **REQ-INF-001 (Network Perimeter):** Access controls must be applied to exposed ports, including `11112/TCP` (DICOM) and `8080/TCP` (Gateway).
    *   **Evidence:** `DEFINED`

*   **REQ-INF-002 (Persistent Storage):** A persistent storage volume is required for study data, queues, and audit logs.
    *   **Evidence:** `DEFINED`

*   **REQ-INF-003 (WORM Responsibility):** The system provides logical immutability (SHA-256). Physical WORM infrastructure is the responsibility of the system administrator.
    *   **Evidence:** `DEFINED`