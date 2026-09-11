# CONTRATO Y COMPORTAMIENTO DICOM
**DOCUMENT STATUS:** `AUTHORITATIVE`  
**AUDIT BASELINE:** `FROZEN NORMATIVE`  

Este documento especifica las capacidades DICOM del sistema RX DISPATCH.

## 1. Capacidades del Sistema

### C-STORE SCP (Service Class Provider)
*   **Componente:** `rx-study`
*   **Puerto:** `11112/TCP`
*   **Responsabilidad:** Actúa como un servidor que escucha y acepta asociaciones DICOM para recibir y almacenar instancias de imágenes médicas (operación C-STORE).
*   **Estado:** `DEFINED` en la arquitectura del sistema como la principal puerta de entrada para las modalidades de imagen.

## 2. Capacidades No Implementadas (o No Verificadas)

*   **C-ECHO:** No hay evidencia de implementación de un servicio de verificación de conectividad DICOM.
*   **C-STORE SCU (Service Class User):** No hay evidencia de que el sistema actúe como un cliente para enviar imágenes a otros PACS.

## 3. Uso del `Calling AE Title`

El `Calling AE Title` de una modalidad que inicia una asociación puede ser utilizado por `rx-study` como un identificador para aplicar lógicas de negocio o listas de admisión (`allowlist`). Sin embargo, este es un control a nivel de aplicación y no debe ser considerado un mecanismo de autenticación criptográfica o de seguridad perimetral.