# ROADMAP DE RECONSTRUCCIÓN DE RX DISPATCH
**DOCUMENT STATUS:** `BASELINE`
**PURPOSE:** Guiar la reconstrucción del sistema de forma ordenada y verificable.

Este roadmap sigue los principios del `MANIFEST.md`. Cada fase debe completarse y sus hallazgos deben ser resueltos antes de proceder a la siguiente. El objetivo es construir sobre evidencia verificada, no sobre suposiciones.

---

## FASE 1: GOBERNANZA DOCUMENTAL (BASELINE)

*   **Objetivo:** Establecer un conjunto de documentos normativos limpios, coherentes y sin contradicciones.
*   **Entregable:** Los 13 documentos canónicos (`README.md`, `ARCHITECTURE.md`, `GITHUB.md`, etc.) actualizados al estado `DEFINED`, eliminando toda evidencia heredada no verificada.
*   **Estado Actual:** `COMPLETED`. Este es el baseline sobre el cual se trabajará.

---

## FASE 2: DEFINICIÓN Y VERIFICACIÓN DE CONTRATOS

*   **Objetivo:** Definir y codificar las estructuras de datos (DTOs) y las firmas de API para la comunicación entre servicios.
*   **Entregable:** Paquetes `internal/contracts` y `internal/models` implementados y estáticamente verificables.
*   **Estado Actual:** `PENDING`.

---

## FASE 3: RECONSTRUCCIÓN DE COMPONENTES (SCAFFOLDING)

*   **Objetivo:** Crear la estructura de directorios y los puntos de entrada (`main.go`) para cada uno de los 9 servicios, sin lógica de negocio completa, pero con la configuración y los endpoints HTTP básicos.
*   **Entregable:** Esqueleto funcional de los 9 microservicios que pueden compilar e iniciarse.
*   **Estado Actual:** `PENDING`.

---

## FASE 4: IMPLEMENTACIÓN Y PRUEBAS UNITARIAS

*   **Objetivo:** Implementar la lógica de negocio de cada servicio, siguiendo los contratos definidos. Cada pieza de funcionalidad debe ir acompañada de pruebas unitarias.
*   **Entregable:** Código fuente de los servicios con una cobertura de pruebas unitarias aceptable.
*   **Estado Actual:** `PENDING`.

---

## FASE 5: PRUEBAS DE INTEGRACIÓN Y VERIFICACIÓN E2E

*   **Objetivo:** Verificar que los servicios se comunican correctamente entre sí y que el flujo de datos completo (DICOM -> Delivery) funciona como se espera en un entorno controlado.
*   **Entregable:** Scripts de prueba de integración y evidencia de un flujo E2E exitoso.
*   **Estado Actual:** `PENDING`.

---

## FASE 6: AUDITORÍA DE SEGURIDAD Y CIERRE

*   **Objetivo:** Realizar una auditoría de seguridad sobre la nueva implementación, verificando el cumplimiento de los requerimientos de `SECURITY.md`.
*   **Entregable:** Reporte de auditoría de seguridad y plan de mitigación de hallazgos.
*   **Estado Actual:** `PENDING`.