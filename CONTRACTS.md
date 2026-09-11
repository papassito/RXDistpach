# ESPECIFICACIÓN DE CONTRATOS DE API
**DOCUMENT STATUS:** `AUTHORITATIVE`  
**AUDIT BASELINE:** `FROZEN NORMATIVE`  
**ALINEACIÓN:** Contratos de comunicación entre servicios.

Todos los servicios se comunican mediante HTTP REST con esquemas JSON explícitos definidos en `internal/contracts/`.

---

## 1. Contratos de `rx-reader` (:8085)

### `POST /reader/analyze`
Ejecuta la lectura genérica automatizada no diagnóstica.
*   **Estado:** `DEFINED`

**Request (`contracts.RequestGenericReadingRequest`):**
```json
{
  "studyId": "string",
  "studyType": "string",
  "anatomicalRegion": "string"
}