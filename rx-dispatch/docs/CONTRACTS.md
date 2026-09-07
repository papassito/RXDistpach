# ESPECIFICACIÓN DE CONTRATOS - RX DISPATCH

Todos los microservicios se comunican mediante HTTP REST con esquemas JSON explícitos ubicados en `internal/contracts/`.

---

## 1. Health Checks (Todos los Servicios)

Cada uno de los 9 servicios expone un endpoint `/healthz`:

```http
GET /healthz
```

**Respuesta (200 OK):**
```json
{
  "service": "rx-reader",
  "port": 8085,
  "status": "UP",
  "timestamp": "2026-09-04T16:50:00Z",
  "version": "2.0.0-decoupled"
}
```

---

## 2. Contratos por Servicio

### RX Security (`:8081`)
- `POST /security/authenticate`: Solicita token de acceso (`AuthenticateRequest` -> `AuthenticateResponse`).
- `POST /security/authorize`: Valida token y scope (`AuthorizeRequest` -> `AuthorizeResponse`).

### RX Study (`:8082`)
- `POST /studies`: Registra nuevo estudio radiológico (`IngestStudyRequest` -> `IngestStudyResponse`).
- `GET /studies`: Lista los estudios registrados.
- `GET /studies/{id}`: Obtiene estudio por ID.
- `PATCH /studies/{id}/status`: Actualiza etapa del ciclo de vida (`UpdateStudyStatusRequest`).

### RX Storage (`:8083`)
- `POST /storage/artifacts`: Almacena un artefacto (`StoreArtifactRequest` -> `StoreArtifactResponse`).
- `GET /storage/artifacts/{id}`: Recupera artefacto (`GetArtifactResponse`).
- `POST /storage/verify`: Verifica integridad SHA-256 (`VerifyArtifactRequest` -> `VerifyArtifactResponse`).

### RX Image (`:8084`)
- `POST /image/process`: Localiza, valida integridad y genera derivado (`ProcessImageRequest` -> `ProcessImageResponse`).
- `POST /image/validate`: Comprueba hash contra el registro (`ValidateIntegrityRequest` -> `ValidateIntegrityResponse`).

### RX Reader (`:8085`)
- `POST /reader/analyze`: Ejecuta lectura genérica automatizada no diagnóstica (`RequestGenericReadingRequest` -> `RequestGenericReadingResponse`).
  - Salida garantizada: `TYPE: GENERIC_AUTOMATED`, `STATUS: WITHOUT_MEDICAL_SIGNATURE`, `MEDICAL_REPORT: NOT_INCLUDED`, `mandatoryDisclaimer`.

### RX Result (`:8086`)
- `POST /result/build`: Ensambla el paquete final de resultado (`BuildResultRequest` -> `BuildResultResponse`).

### RX Delivery (`:8087`)
- `POST /delivery/dispatch`: Despacha el resultado al destinatario (`DispatchDeliveryRequest` -> `DispatchDeliveryResponse`).
- `GET /delivery/records`: Consulta el historial de envíos.

### RX Audit (`:8088`)
- `POST /audit/events`: Registra un evento operacional (`RecordAuditEventRequest` -> `RecordAuditEventResponse`).
- `GET /audit/events`: Consulta eventos con filtro opcional `?studyId=...`.

### RX Gateway (`:8080`)
- `GET /gateway/topology`: Inspecciona el estado de salud simultáneo de los 9 microservicios.
- `POST /api/v1/dispatch-flow`: Orquesta la canalización de extremo a extremo a través de los servicios desacoplados.
