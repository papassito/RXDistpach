# RX DISPATCH 
**AUTOMATED MEDICAL IMAGING & TECHNICAL EXTRACTION ENGINE**

**DOCUMENT STATUS:** `AUTHORITATIVE`  
**AUDIT BASELINE:** `FROZEN NORMATIVE`  
**CLASSIFICATION:** `MISSION-CRITICAL MIDDLEWARE`  
**REGULATORY FRAMEWORK ALIGNMENT:** NOM-024-SSA3-2012 (MEX), LFPDPPP (MEX), HIPAA/HITECH (USA), 21 CFR Part 11 (FDA), IEC 62304 (INTL).

## 1. DECLARACIÓN DE PROPÓSITO Y MARCO REGULATORIO (INTENDED USE)
RX DISPATCH es un middleware radiológico determinístico implementado en Go. Ejerce funciones operativas de servidor DICOM C-STORE SCP (NEMA PS3), procesando instancias de imagen desde modalidades radiológicas locales, extrayendo metadatos técnicos y orquestando cargas útiles (*payloads*) en una cola de alta resiliencia para su despacho asíncrono.

> **⚠️ EXENCIÓN DIAGNÓSTICA Y LÍMITE DE RESPONSABILIDAD CLÍNICA**
> El uso previsto (*intended use*) del aplicativo se limita estrictamente al transporte de datos técnicos, almacenamiento local, conversión de formatos y extracción de metadatos. **No incluye ni provee algoritmos de interpretación diagnóstica, detección de patologías, ni emisión de juicios clínicos.** La determinación formal de clasificación como dispositivo médico (SaMD / MDDS) queda expresamente sujeta a la evaluación jurídica y aplicabilidad de *predicate rules* ante las agencias regulatorias competentes (COFEPRIS, FDA, EMA).

## 2. PRINCIPIOS FUNDAMENTALES DEL SISTEMA
1. **Desacoplamiento Estricto de Dominios**: Ningún proceso concentra toda la lógica. Cada servicio es un binario independiente con su propio punto de entrada `main.go`. La estructura completa del proyecto está definida normativamente en `MAP.md`.
2. **Inmutabilidad Radiográfica**: La imagen original se sella criptográficamente con SHA-256 en `rx-storage`. El servicio `rx-image` genera derivados para lectura visual, garantizando que el archivo original no sufra alteración lógica.
3. **Lectura Estrictamente No Diagnóstica**: Todo resultado automatizado se emite con `TYPE: GENERIC_AUTOMATED`, `STATUS: WITHOUT_MEDICAL_SIGNATURE` y una leyenda legal obligatoria que advierte de su naturaleza no clínica.
4. **Trazabilidad Completa**: El servicio `rx-audit` registra cada paso del ciclo de vida del estudio de forma inmutable, asegurando una cadena de custodia lógica y pericial completa.

## 3. VECTORES DE EJECUCIÓN (DEPLOYMENT MODES)
El sistema garantiza coherencia transaccional y cumplimiento de límites de dominio bajo dos topologías estrictas:
*   **Modo Distribuido:** Ejecución de 9 procesos de dominio aislados, vinculados por ruteo estático interno HTTP (127.0.0.1). 
*   **Modo Edge:** Binario autocontenido (`rx-dispatch-edge.exe`) que consolida los dominios de servicio conservando fronteras lógicas, diseñado para el despliegue en infraestructura de centro de diagnóstico local.

## 4. PROTOCOLO DE INICIALIZACIÓN (BOOTSTRAPPING)
**Requisito Base:** Cadena de herramientas Go versión 1.22+.

**Secuencia de Compilación y Ejecución (Entorno Windows):**
```powershell
# 1. Validar el proyecto (formato, tests, dependencias)
.\scripts\validate_project.ps1

# 2. Compilar el binario para el modo Edge
go build -ldflags="-s -w" -o bin/rx-dispatch-edge.exe ./cmd/gateway

# 3. Ejecutar el sistema en modo Edge
.\bin\rx-dispatch-edge.exe --config .\config.json

.\scripts\validate_project.ps1
go build -ldflags="-s -w" -o bin/rx-dispatch-edge.exe ./cmd/gateway
.\bin\rx-dispatch-edge.exe --config .\config.json