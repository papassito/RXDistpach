# CONTRATOS DE API REST E IDEMPOTENCIA
**DOCUMENT STATUS:** `AUTHORITATIVE`  
**AUDIT BASELINE:** `FROZEN NORMATIVE`  
**REGULATORY FRAMEWORK ALIGNMENT:** Previsión técnica contra duplicidad de registros clínicos y transacciones redundantes.

## 1. ORQUESTACIÓN DE DESPACHO ASÍNCRONO
**Endpoint:** `POST /api/v1/dispatch-flow`  
**Mandato de Cabecera:** `Idempotency-Key: <UUID>`

**Payload Canónico Estructural:**
```json
{
  "studyId": "STU-8829-X",
  "deliveryTarget": "+521234567890",
  "channel": "WEBHOOK"
}