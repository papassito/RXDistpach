### 6. `SECURITY.md`

```markdown
# GOBERNANZA DE SEGURIDAD, INTEGRIDAD Y AUDITORÍA PERICIAL
**DOCUMENT STATUS:** `AUTHORITATIVE`  
**AUDIT BASELINE:** `FROZEN NORMATIVE`  
**REGULATORY FRAMEWORK ALIGNMENT:** LFPDPPP, GDPR (Art. 32), Título 21 CFR Parte 11.10.

## 1. GESTIÓN DE ASOCIACIÓN DICOM (A-ASSOCIATE-RJ)
El sistema ejerce un nivel secundario de autorización en la capa de aplicación DICOM (conforme a la familia de estándares NEMA PS3.8 para comunicaciones de red). Tras la superación del perímetro de red IP, se evalúa el `Calling AE Title`. Ante una identidad no configurada en los registros autorizados, el servicio debe emitir el comando de rechazo protocolario normativo (`A-ASSOCIATE-RJ`), evitando el cierre arbitrario del socket TCP para asegurar la estabilidad de la modalidad de imagen.

## 2. INTEGRIDAD CRIPTOGRÁFICA Y PRESUNCIÓN DE ALTERACIÓN
La responsabilidad de preservación se divide obligatoriamente en dos dominios:
*   **Integridad Lógica Interna (SHA-256):** El aplicativo calcula y vincula un hash criptográfico `SHA-256` en el momento de la ingesta para todo archivo depositado en `ORIG-*`. Este control no previene el borrado, pero garantiza la capacidad forense de **detectar modificaciones posteriores** respecto de la carga útil original de imagen médica.
*   **Límites de Destrucción Transaccional:** Las subrutinas de purga del sistema operativo (Disk Watermark Purge) tienen un cerco de autorización estricto sobre datos efímeros y derivados (`DERIV-*`). El aplicativo no emite comandos destructivos sobre el subdirectorio de orígenes (`ORIG-*`).

## 3. TRAZABILIDAD Y CADENA DE CUSTODIA LÓGICA (AUDIT TRAIL)
El componente `rx-audit` establece una bitácora transaccional para eventos críticos, resguardando la **integridad lógica del registro**. Consolida una matriz de verificación que captura:
*   **Actor:** Entidad u origen que desencadena el evento (e.g., `System:rx-study`, UUID de credencial de consumo de API).
*   **Timestamp:** Marca de tiempo estandarizada y referenciada strictly al Tiempo Universal Coordinado (UTC ISO-8601).
*   **Event:** Taxonomía predefinida de la acción ejecutada (`ASSOCIATION_ACCEPTED`, `DISPATCH_QUEUED`).
*   **Resource:** Vectores de identidad correlacionada pertinentes a la transacción (`Calling AE Title`, `Association ID`, `Source IP Address`).
*   **Result:** Disposición final determinística del evento (`SUCCESS` / `FAILURE`).