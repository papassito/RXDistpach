# RX DISPATCH — Component Catalog

**DOCUMENT STATUS:** `CONSOLIDATED BASELINE`
**PURPOSE:** Catalog of declared components, their responsibilities, and evidence status.

---

## 1. DOMAIN SERVICES

RX DISPATCH is composed of 9 declared domain services. The evidence level for each is based on the current audit scope.

### `rx-gateway`
*   **Responsibility:** API Gateway, ingress, and internal routing.
*   **Evidence Level:** `DEFINED`

### `rx-security`
*   **Responsibility:** Authorization boundary and access control for REST APIs.
*   **Evidence Level:** `DEFINED`

### `rx-study`
*   **Responsibility:** DICOM C-STORE SCP for study ingestion and metadata extraction.
*   **Declared Port:** `11112/TCP`
*   **Evidence Level:** `DEFINED`

### `rx-storage`
*   **Responsibility:** Manages local persistence of original (`ORIG-*`) and derived files.
*   **Evidence Level:** `DEFINED`

### `rx-image`
*   **Responsibility:** Generates non-destructive derived image representations (`DERIV-*`).
*   **Evidence Level:** `DEFINED`

### `rx-reader`
*   **Responsibility:** Generates a **generic, automated, non-diagnostic reading** and ensures the mandatory injection of the clinical disclaimer.
*   **Declared Endpoint:** `POST /reader/analyze`
*   **Evidence Level:** `DEFINED`

### `rx-result`
*   **Responsibility:** Consolidates metadata and generated artifacts into a final result payload.
*   **Evidence Level:** `DEFINED`

### `rx-delivery`
*   **Responsibility:** Manages outbound dispatch to external systems using a `Store-and-Forward` queue.
*   **Evidence Level:** `DEFINED`

### `rx-audit`
*   **Responsibility:** Centralized transactional audit log for traceability.
*   **Evidence Level:** `DEFINED`

---

## 2. SHARED INTERNAL PACKAGES

The `internal/` directory contains shared logic used across services.

| Package              | Declared Responsibility                      | Evidence Level                  |
| -------------------- | -------------------------------------------- | ------------------------------- |
| `internal/config`    | Configuration and service topology loading.  | `DECLARED`                      |
| `internal/contracts` | Data Transfer Objects (DTOs) for APIs.       | `DECLARED`                      |
| `internal/models`    | Core domain models (e.g., `GenericReading`). | `DECLARED`                      |
| `internal/shared`    | Normative constants (e.g., disclaimers).     | `DECLARED`                      |
| `internal/transport` | HTTP client helpers for inter-service calls. | `DECLARED`                      |
| `internal/webassets` | Static asset embedding.                      | `DECLARED`                      |

---

## 3. EXTERNAL SYSTEMS

### PxLab
*   **Integration Policy:** Any component interacting with PxLab must adhere to a strict `READ ONLY` policy.
*   **Adapter Component:** The specific component responsible for this integration is currently `UNVERIFIED`.

---

## 4. COMPONENT VERIFICATION PRINCIPLE

The status of each component is subject to the phased audit process defined in `MANIFEST.md`. The existence of a component in this document does not imply it is fully functional or certified.

The rule `DOCUMENTED != IMPLEMENTED` is strictly enforced.
