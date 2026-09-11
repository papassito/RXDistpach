# RX DISPATCH — AUDIT & VERIFICATION MANIFEST

**DOCUMENT:** `MANIFEST.md`
**PURPOSE:** Gobierno de Auditoría, Verificación y Corrección por Fases
**MODE:** `FORENSIC / EVIDENCE-DRIVEN`
**PROJECT:** `RX DISPATCH`

---

# 1. PROPÓSITO

Este manifiesto define el procedimiento oficial para auditar, verificar y posteriormente corregir RX DISPATCH de forma progresiva.

El sistema será inspeccionado **PHASE por PHASE**.

Ninguna fase posterior podrá utilizar como hecho una afirmación procedente de una fase anterior que permanezca clasificada como `UNVERIFIED`.

La auditoría debe determinar qué existe realmente, qué está conectado realmente, qué funciona realmente, qué solamente está definido y qué pudiera corresponder a simulaciones, placeholders, mocks, stubs o implementaciones incompletas.

---

# 2. REGLA MAESTRA

```text
EVIDENCE > CLAIM
CODE > DESCRIPTION
EXECUTION > CODE PRESENCE
REAL INTEGRATION > SIMULATION
VERIFICATION > ASSUMPTION
```

Los documentos normativos establecen lo que RX DISPATCH debe ser.

El código establece lo que fue implementado.

La ejecución controlada establece lo que realmente funciona.

La evidencia establece lo que puede declararse verificado.

---

# 3. NIVELES DE EVIDENCIA

Toda afirmación deberá clasificarse exclusivamente como:

```text
DEFINED
IMPLEMENTED
STATICALLY VERIFIED
TESTED
E2E VERIFIED
CERTIFIED
UNVERIFIED
```

Regla obligatoria:

```text
DEFINED
!= IMPLEMENTED

IMPLEMENTED
!= TESTED

TESTED
!= E2E VERIFIED

E2E VERIFIED
!= CERTIFIED
```

No se permite elevar automáticamente el nivel de evidencia.

---

# 4. REGLAS DE AUDITORÍA

Durante una fase declarada `AUDIT ONLY`:

```text
ZERO NEW CODE
ZERO PATCHES
ZERO MUTATION
ZERO INVENTION
ZERO INNOVATION
ZERO ADD CODE
ZERO FILE DELETION
ZERO FILE RENAME
ZERO FILE MOVE
ZERO STRUCTURAL MODIFICATION
ONLY READ
ONLY AUDIT
```

Si existe un problema:

```text
REPORT IT
```

Si falta algo:

```text
REPORT IT
```

Si existe pero no puede demostrarse:

```text
UNVERIFIED
```

Si contradice el baseline:

```text
MISMATCH
```

Si aparenta funcionar mediante simulación:

```text
SIMULATION SUSPECTED
```

La auditoría no deberá corregir el hallazgo dentro de la misma operación de auditoría.

---

# 5. INVARIANTE CLÍNICA

RX DISPATCH puede contener generación automatizada de lectura genérica.

Debe mantenerse permanentemente:

```text
AUTOMATED GENERIC READING
!= MEDICAL INTERPRETATION
!= MEDICAL DIAGNOSIS
!= SIGNED RADIOLOGY REPORT
```

La existencia de:

```text
ReadingTypeGenericAutomated
StatusWithoutMedicalSignature
MedicalReportNotIncluded
MandatoryLegalDisclaimer
```

debe verificarse independientemente de cualquier afirmación documental.

---

# 6. INVARIANTE PXLAB

Toda integración con PxLab será:

```text
READ ONLY
```

Se consideran operaciones prohibidas:

```text
INSERT
UPDATE
DELETE
ALTER
CREATE
DROP
TRUNCATE
MERGE
```

También deberán detectarse mecanismos indirectos capaces de producir mutación.

La presencia documental de la regla no demuestra su cumplimiento en ejecución.

---

# 7. PHASE 0 — BASELINE & INVENTORY

## Objetivo

Determinar exactamente qué existe antes de analizar qué hace.

## Auditar

```text
ROOT
DIRECTORIES
FILES
DOCS
PHASES
SCRIPTS
ASSETS
BINARIES
CONFIG
GO.MOD
GO.SUM
REPOSITORY
GITHUB
```

Buscar:

```text
DUPLICATE FILES
GHOST FILES
GHOST FOLDERS
CORRUPTED FILES
EMPTY FILES
ORPHAN FILES
GENERATED FILES
UNTRACKED STRUCTURES
UNEXPECTED BINARIES
```

Comparar estructura física contra `MAP.md`.

## Resultado

Generar inventario clasificado:

```text
EXPECTED + FOUND
EXPECTED + MISSING
UNEXPECTED + FOUND
DUPLICATED
ORPHAN
UNVERIFIED
```

## Gate

No continuar hasta conocer la estructura real del repositorio.

---

# 8. PHASE 1 — DOCUMENTARY GOVERNANCE

## Objetivo

Determinar la coherencia del baseline normativo.

## Auditar

```text
README.md
REQUIREMENTS.md
ARCHITECTURE.md
CONTRACTS.md
MAP.md
SECURITY.md
DICOM.md
ENVIOS.md
TELEMETRY.md
ROADMAP.md
COMPONENTS.md
PHASES
DOCS
```

Buscar:

```text
CONTRADICTIONS
DUPLICATED AUTHORITY
OBSOLETE REQUIREMENTS
BROKEN REFERENCES
MISSING CONTRACTS
INVALID STATUS CLAIMS
UNSUPPORTED CERTIFICATION CLAIMS
DOCUMENTATION DRIFT
```

## Regla

Un encabezado:

```text
AUTHORITATIVE
FROZEN
SEALED
```

no demuestra por sí mismo que el documento coincida con la implementación.

## Gate

Baseline documental clasificado y contradicciones identificadas.

---

# 9. PHASE 2 — SOURCE TREE & CODE INTEGRITY

## Objetivo

Radiografiar la implementación física.

## Auditar

```text
GO FILES
TS FILES
JSON FILES
PS1 FILES
PACKAGES
MODULES
FUNCTIONS
IMPORTS
GENERATED CODE
EMBEDDED ASSETS
SYNTAX
CHARACTERS
ENCODING
```

Revisar especialmente:

```text
cmd/
internal/
scripts/
assets/
static/
config.json
go.mod
go.sum
```

Buscar:

```text
DEAD CODE
ORPHAN FUNCTIONS
DUPLICATED FUNCTIONS
BROKEN IMPORTS
PLACEHOLDERS
TODO
FIXME
HARDCODED RESPONSES
HARDCODED ENDPOINTS
FALLBACK SUCCESS
```

## Gate

Código inventariado y relaciones estructurales identificadas.

---

# 10. PHASE 3 — COMPONENT VERIFICATION

## Objetivo

Auditar individualmente los dominios declarados.

Baseline esperado:

```text
rx-gateway
rx-security
rx-study
rx-storage
rx-image
rx-reader
rx-result
rx-delivery
rx-audit
```

Para cada componente determinar:

```text
DEFINED?
SOURCE EXISTS?
ENTRYPOINT EXISTS?
CONFIG EXISTS?
CONTRACT EXISTS?
IMPLEMENTATION EXISTS?
DEPENDENCIES EXIST?
REAL FUNCTION OR PLACEHOLDER?
```

## Regla

La existencia de:

```text
cmd/<service>/main.go
```

no demuestra automáticamente que el servicio funcione.

## Gate

Los nueve dominios quedan individualmente clasificados.

---

# 11. PHASE 4 — CONTRACT TRACEABILITY

## Objetivo

Trazar:

```text
REQUIREMENT
     ↓
CONTRACT
     ↓
MODEL
     ↓
FUNCTION
     ↓
CALLER
     ↓
RECEIVER
```

Auditar:

```text
REST ENDPOINTS
REQUEST STRUCTURES
RESPONSE STRUCTURES
JSON
STATUS CODES
ERROR CONTRACTS
AUDIT EVENTS
```

Incluir específicamente:

```text
POST /reader/analyze
RequestGenericReadingRequest
RequestGenericReadingResponse
GenericReading
RecordAuditEventRequest
```

## Gate

Cada contrato queda clasificado como:

```text
DEFINED ONLY
IMPLEMENTED
IMPLEMENTED + CONSUMED
BROKEN
ORPHAN
MISMATCH
UNVERIFIED
```

---

# 12. PHASE 5 — INTERCONNECTIONS

## Objetivo

Reconstruir el grafo real del sistema.

No utilizar únicamente el diagrama documental.

Trazar físicamente:

```text
SOURCE SERVICE
      ↓
CONFIG
      ↓
DESTINATION URL
      ↓
TRANSPORT
      ↓
CONTRACT
      ↓
RECEIVER
      ↓
RESPONSE
```

Comparar contra el flujo esperado:

```text
DICOM
  ↓
rx-study
  ↓
rx-storage / rx-image
  ↓
rx-reader
  ↓
rx-result
  ↓
rx-delivery
  ↓
OUTBOUND
```

y sus conexiones hacia:

```text
rx-audit
rx-security
rx-gateway
```

Buscar:

```text
BROKEN CONNECTION
GHOST ENDPOINT
ORPHAN ENDPOINT
WRONG PORT
HARDCODED URL
CONFIG MISMATCH
CONTRACT MISMATCH
MISSING RECEIVER
CIRCULAR DEPENDENCY
```

## Gate

Debe existir un mapa `DECLARED vs ACTUAL`.

---

# 13. PHASE 6 — SECURITY

## Objetivo

Auditar las fronteras reales de seguridad.

Revisar:

```text
NETWORK BINDINGS
127.0.0.1
EXPOSED PORTS
AUTHENTICATION
AUTHORIZATION
NODE IDENTITY
USERS
TENANT ID
SECRETS
CONFIG
INPUT VALIDATION
PATH HANDLING
FILESYSTEM PERMISSIONS
TLS
TIMEOUTS
PAYLOAD LIMITS
DEPENDENCIES
AUDIT INTEGRITY
```

Verificar especialmente:

```text
11112/TCP
8080/TCP
8081-8088/TCP
443/TCP
```

según configuración real.

## Gate

Hallazgos clasificados por evidencia y severidad.

No declarar seguridad verificada globalmente a partir de un único componente.

---

# 14. PHASE 7 — SIMULATION HUNT

## Objetivo

Detectar funcionalidad aparente.

Buscar explícitamente:

```text
MOCK
FAKE
STUB
DEMO
SAMPLE
SIMULATION
SYNTHETIC
PLACEHOLDER
HARDCODED SUCCESS
HARDCODED RESPONSE
RANDOM EVIDENCE
FAKE ACK
FAKE HEALTH
FAKE TELEMETRY
FAKE DICOM
NO-OP
```

También buscar comportamiento equivalente aunque no utilice esas palabras.

Clasificar:

```text
REAL IMPLEMENTATION
TEST-ONLY SIMULATION
LEGITIMATE MOCK
PRODUCTION-REACHABLE SIMULATION
PLACEHOLDER
UNVERIFIED
```

## Gate

Ninguna simulación alcanzable desde producción puede confundirse con evidencia funcional real.

---

# 15. PHASE 8 — DOMAIN RX VERIFICATION

## Objetivo

Auditar las funciones específicas del producto.

### DICOM

```text
C-STORE SCP
ASSOCIATION
AE TITLE
METADATA EXTRACTION
INGEST
```

### STORAGE

```text
ORIG-*
SHA-256
PERSISTENCE
PURGE
DISK WATERMARK
```

### IMAGE

```text
DERIV-*
TRANSCODING
ORIGINAL PRESERVATION
```

### READER

```text
GENERIC AUTOMATED READING
ANATOMICAL REGION
MANDATORY DISCLAIMER
NON-DIAGNOSTIC CLASSIFICATION
AUDIT EVENT
```

### RESULT

```text
AGGREGATION
HASH
PAYLOAD
```

### DELIVERY

```text
QUEUE
STORE-AND-FORWARD
RETRY
BACKOFF
ACK
DELIVERY
VERIFICATION
```

### AUDIT

```text
EVENT
PERSISTENCE
INTEGRITY
CORRELATION
TRACEABILITY
```

## Gate

Cada dominio queda separado entre:

```text
DEFINED
IMPLEMENTED
STATICALLY VERIFIED
TESTABLE
UNVERIFIED
```

---

# 16. PHASE 9 — DATABASE & PERSISTENCE

## Objetivo

Determinar qué estado persiste realmente.

Auditar:

```text
DATABASE
FILESYSTEM STORAGE
QUEUE STORAGE
AUDIT STORAGE
CONFIGURATION
TEMPORARY STORAGE
DERIVED STORAGE
```

Verificar:

```text
PERSISTENCE
TRANSACTIONS
CONCURRENCY
RECOVERY
RETENTION
PURGE
INTEGRITY
```

PxLab deberá auditarse específicamente para demostrar ausencia de rutas de escritura.

## Gate

Mapa de persistencia real terminado.

---

# 17. PHASE 10 — TEST & SMOKE AUDIT

## Objetivo

Auditar primero los tests existentes antes de ejecutarlos.

Buscar:

```text
*_test.go
TEST SCRIPTS
SMOKE TEST
FIXTURES
MOCKS
TEST DATABASES
TEMP FILES
NETWORK CALLS
PROCESS SPAWNING
```

Cada prueba deberá clasificarse previamente:

```text
READ-SAFE
MUTATING
EXTERNAL-SIDE-EFFECT
UNKNOWN
```

No ejecutar pruebas potencialmente mutantes durante una fase `ONLY READ`.

## Gate

Se conoce exactamente qué prueba puede ejecutarse y qué evidencia produciría.

---

# 18. PHASE 11 — CONTROLLED EXECUTION

Esta fase requiere autorización independiente.

Aquí termina `ONLY READ` absoluto cuando la ejecución inevitablemente produzca procesos, logs, archivos temporales u otros efectos controlados.

Objetivo:

```text
BUILD
START
SMOKE
REQUEST
RESPONSE
ERROR
RECOVERY
```

No se realizarán correcciones todavía.

Sólo ejecución controlada y recopilación de evidencia.

---

# 19. PHASE 12 — FUNCTIONAL VERIFICATION

## Objetivo

Demostrar comportamiento real.

Ejemplos:

```text
REAL REQUEST
REAL SERVICE
REAL RESPONSE
REAL AUDIT
REAL STORAGE
REAL ACK
REAL ERROR
REAL RETRY
```

Para `rx-reader`:

```text
POST /reader/analyze
        ↓
REAL RESPONSE
        ↓
GENERIC READING
        ↓
MANDATORY DISCLAIMER
        ↓
AUDIT EVENT
        ↓
rx-audit RECEIVES EVENT
```

Sólo entonces podrá elevarse evidencia estática a evidencia funcional cuando corresponda.

---

# 20. PHASE 13 — PERFORMANCE & RESILIENCE

Auditar y posteriormente medir:

```text
LATENCY
THROUGHPUT
MEMORY
CPU
QUEUE DEPTH
BACKPRESSURE
TIMEOUTS
RETRIES
AUDIT LATENCY
NETWORK FAILURE
SERVICE FAILURE
RECOVERY
```

Nunca declarar:

```text
HIGH PERFORMANCE
SCALABLE
RESILIENT
```

sin medición.

---

# 21. PHASE 14 — REPOSITORY & SUPPLY CHAIN

Auditar:

```text
GIT
GITHUB
GO.MOD
GO.SUM
DEPENDENCIES
DIRECT DEPENDENCIES
TRANSITIVE DEPENDENCIES
VERSIONS
CHECKSUMS
LICENSES
BUILD SCRIPTS
ARTIFACTS
```

Regla:

```text
go.mod != SBOM
```

Si existe SBOM:

```text
SPDX
CYCLONEDX
```

deberá auditarse explícitamente.

---

# 22. PHASE 15 — FINAL E2E VERIFICATION

Sólo después de cerrar las fases anteriores.

Flujo objetivo:

```text
REAL DICOM SOURCE
        ↓
REAL INGEST
        ↓
REAL STORAGE
        ↓
REAL PROCESSING
        ↓
REAL GENERIC READING
        ↓
REAL RESULT
        ↓
REAL QUEUE
        ↓
REAL DELIVERY
        ↓
REAL ACK
        ↓
REAL AUDIT
        ↓
REAL EVIDENCE
```

Una simulación no satisface esta fase.

---

# 23. PHASE 16 — CORRECTION

La corrección se realiza **después de la auditoría correspondiente**, nunca mezclada con ella.

Ciclo obligatorio:

```text
AUDIT PHASE
    ↓
FINDINGS
    ↓
CLASSIFICATION
    ↓
CORRECTION AUTHORIZATION
    ↓
CORRECTION
    ↓
RE-AUDIT
    ↓
CLOSE PHASE
```

No avanzar dejando correcciones sin reauditar.

---

# 24. PHASE 17 — FINAL GOVERNANCE REVIEW

Comparar finalmente:

```text
README
REQUIREMENTS
ARCHITECTURE
CONTRACTS
MAP
COMPONENTS
SECURITY
DICOM
ENVIOS
TELEMETRY
ROADMAP
```

contra:

```text
ACTUAL SOURCE
ACTUAL CONFIG
ACTUAL EXECUTION
ACTUAL TEST RESULTS
ACTUAL SYSTEM BEHAVIOR
```

Resultado:

```text
DOCUMENTED == IMPLEMENTED?
IMPLEMENTED == TESTED?
TESTED == VERIFIED?
MAP == REPOSITORY?
CONTRACTS == ENDPOINTS?
REQUIREMENTS == BEHAVIOR?
```

---

# 25. FORMATO DE HALLAZGOS

Cada hallazgo deberá contener:

```text
ID:
PHASE:
CATEGORY:
SEVERITY:
FILE:
LINE / SYMBOL:
EXPECTED:
OBSERVED:
EVIDENCE:
STATUS:
IMPACT:
CORRECTION STATUS:
RE-AUDIT STATUS:
```

Estados:

```text
OPEN
CONFIRMED
UNVERIFIED
FALSE POSITIVE
AUTHORIZED FOR CORRECTION
CORRECTED
RE-AUDIT PASSED
CLOSED
```

---

# 26. GATES

Ninguna fase se considera cerrada simplemente porque el auditor terminó de leer.

Una fase sólo puede cerrarse cuando:

```text
SCOPE COMPLETED
+
FINDINGS RECORDED
+
UNVERIFIED ITEMS RECORDED
+
NO HIDDEN MUTATION
+
REQUIRED CORRECTIONS RE-AUDITED
```

Cuando una corrección no sea necesaria:

```text
AUDIT PASSED
```

será suficiente para cerrar la fase.

---

# 27. REGLA CONTRA EVIDENCIA HEREDADA NO VERIFICADA

Reportes, comentarios, resultados o afirmaciones históricas procedentes de herramientas automáticas no constituyen evidencia suficiente por sí mismos.

Se tratarán como:

```text
LEGACY CLAIM
UNTRUSTED UNTIL INDEPENDENTLY VERIFIED
```

Pueden utilizarse para localizar áreas que requieren inspección.

No pueden utilizarse por sí solos para elevar:

```text
DEFINED
→ IMPLEMENTED

IMPLEMENTED
→ TESTED

TESTED
→ E2E VERIFIED
```

---

# 28. RESULTADO FINAL

La auditoría deberá terminar con una matriz:

```text
PHASE
SCOPE
STATUS
FINDINGS
P0
P1
P2
P3
UNVERIFIED
CORRECTED
RE-AUDITED
EVIDENCE LEVEL
```

Y un estado general:

```text
RX DISPATCH STATUS

DOCUMENTATION:
ARCHITECTURE:
SOURCE:
CONTRACTS:
INTERCONNECTIONS:
SECURITY:
SIMULATIONS:
DICOM:
STORAGE:
READER:
DELIVERY:
AUDIT:
DATABASE:
TESTS:
PERFORMANCE:
SUPPLY CHAIN:
E2E:
GOVERNANCE:
```

---

# 29. PRINCIPIO FINAL

```text
DO NOT PROVE THAT RX DISPATCH WORKS.

DISCOVER WHAT RX DISPATCH ACTUALLY IS.

THEN PROVE, PHASE BY PHASE,
WHAT REALLY WORKS.
```