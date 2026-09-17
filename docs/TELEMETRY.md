# TELEMETRÍA, MÉTRICAS Y ESTADO DEL SISTEMA
**DOCUMENT STATUS:** `AUTHORITATIVE`  
**AUDIT BASELINE:** `FROZEN NORMATIVE`  

Este documento define los mecanismos de observabilidad y telemetría del sistema RX DISPATCH.

## 1. Verificación de Estado (Health Checks)

*   **Endpoint:** `/healthz`
*   **Disponibilidad:** Todos los 9 microservicios deben exponer este endpoint.
*   **Responsabilidad:** Proporcionar una señal simple y rápida del estado de salud del servicio.
*   **Estado:** `DEFINED` como un requisito para todos los servicios.

**Respuesta (`contracts.HealthResponse`):**
```json
{
  "service": "rx-reader",
  "port": 8085,
  "status": "UP",
  "timestamp": "2026-09-04T16:50:00Z",
  "version": "2.0.0-decoupled"
}
```

## 2. Registros (Logs)

Cada servicio genera registros en la salida estándar (`stdout`) para informar sobre su estado operativo, advertencias y errores. Los registros están prefijados con el nombre del componente (ej. `[RX-READER]`) para facilitar la correlación.
*   **Estado:** `DEFINED` como un requisito para todos los servicios.

## 3. Bitácora de Auditoría (Audit Trail)

El servicio `rx-audit` actúa como un colector centralizado de telemetría de negocio y seguridad. Los eventos enviados por otros servicios (como `generic_reading_generated` de `rx-reader`) constituyen una forma de telemetría de alto nivel que permite trazar el flujo de una transacción a través del sistema.
*   **Estado:** `DEFINED`.