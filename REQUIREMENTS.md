# ESPECIFICACIÓN DE INFRAESTRUCTURA Y ENTORNOS
**DOCUMENT STATUS:** `AUTHORITATIVE`  
**AUDIT BASELINE:** `FROZEN NORMATIVE`  
**REGULATORY FRAMEWORK ALIGNMENT:** NOM-024-SSA3-2012 (MEX), Título 21 CFR Parte 11 (Controles Técnicos), LFPDPPP (MEX).

## 1. PERÍMETRO DE RED Y CONTROL DE ACCESO (ACL ENFORCEMENT)
Para facilitar la observancia técnica de los marcos de seguridad sanitaria, se imponen las siguientes fronteras topológicas:

**Controles de Red Entrantes (Inbound):**
*   `11112 / TCP`: Listener DICOM C-STORE SCP. **Requiere segmentación de red estricta** (VLAN clínica, mitigación por Firewall a nivel IP). *Nota jurídica:* La evaluación del identificador `Calling AE Title` es exclusivamente un mecanismo de capa de aplicación DICOM y no constituye, por sí mismo, una identidad criptográfica o un control de acceso perimetral.
*   `8089 / 8090 / TCP`: Listener de API Gateway según modo de despliegue. Exposición sujeta a la política de control de acceso de la jurisdicción del usuario final.

**Aislamiento de Lazo Local (Loopback Enforcement):**
*   `8081` a `8088`: Vinculación exclusiva e irrenunciable a `127.0.0.1` (loopback).

**Controles de Red Salientes (Outbound):**
*   `443 / TCP`: Salida autorizada para transmisiones externas utilizando obligatoriamente TLS 1.2 o superior.

## 2. GESTIÓN DE ALMACENAMIENTO LOCAL
Se requiere la asignación de un volumen persistente dedicado (`./storage_data`).
`RequiredStorage = OriginalStudyStorage + DerivedStorage + LocalQueueDB + AuditDB + SafetyMargin`

## 3. RESPONSABILIDAD TRASLADADA SOBRE INMUTABILIDAD FÍSICA (WORM)
RX DISPATCH **NO** provee inmutabilidad física ni criptográfica a nivel de bloque de almacenamiento en la capa de aplicación. Para cumplimentar la retención obligatoria de expedientes clínicos (e.g., NOM-004-SSA3-2012):
*   La inmutabilidad de infraestructura WORM (*Write Once, Read Many*), la gestión de instantáneas (*snapshots*) y las políticas de retención física de la ruta `ORIG-*` residen enteramente bajo la jurisdicción administrativa e infraestructura del Responsable del Tratamiento de Datos (Centro de Salud / Hospital).