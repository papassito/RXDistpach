# ARQUITECTURA DEL SISTEMA Y TOPOLOGÍA DE DOMINIOS
**DOCUMENT STATUS:** `AUTHORITATIVE`  
**AUDIT BASELINE:** `FROZEN NORMATIVE`  
**REGULATORY FRAMEWORK ALIGNMENT:** Principios de Arquitectura IEC 62304, ISO/IEC 27701.

## 1. GOBERNANZA ARQUITECTÓNICA
RX DISPATCH opera bajo una topología de **9 procesos de dominio desacoplados**. Se prohíbe explícitamente el enrutamiento dinámico no determinístico (*Service Mesh*). La arquitectura técnica asegura el aislamiento de la Información en Salud Protegida (ePHI / Datos Personales Sensibles) a nivel de proceso.

## 2. CONTRATO DE COMUNICACIÓN INTERNA
Las transacciones RPC inter-dominio operan estrictamente mediante **REST/HTTP síncrono**. El servicio de entrega (`rx-delivery`) recibe la instrucción de trabajo de forma síncrona. La latencia asíncrona (*Store-and-Forward*) aplica única y exclusivamente a la interfaz de salida (*Outbound*) hacia redes externas.

## 3. MATRIZ DE SERVICIOS
1. **rx-gateway:** API Gateway estático. Centraliza la ingesta e inspección de peticiones entrantes.
2. **rx-security:** Barrera de autorización para consumo de API REST.
3. **rx-study:** Motor DICOM (C-STORE SCP). Recibe asociaciones y extrae atributos DICOM base.
4. **rx-storage:** Gestor de I/O en disco. Ejecuta persistencia local y políticas de purga preventiva por umbral (*Disk Watermark*).
5. **rx-image:** Transcodificador. Genera representaciones ligeras (*DERIV-**) para canales técnicos de visualización referencial.
6. **rx-reader:** Extractor de telemetría y garante de la inyección obligatoria de la leyenda de exención diagnóstica.
7. **rx-result:** Consolidador criptográfico y agregador de metadatos.
8. **rx-delivery:** Gestor de Cola Resiliente (*Store-and-Forward*). Mitiga interrupciones WAN mediante retroceso exponencial (*Exponential Backoff*).
9. **rx-audit:** Componente de bitácora transaccional para eventos de acceso y sistema.

## 4. FLUJO DE DATOS Y CICLO DE VIDA DEL ESTUDIO
```text
[Modalidad DICOM]
       │ (C-STORE SCP / TCP 11112)
       ▼
 1. rx-study ─────────► Registra Ingesta ─────────► rx-audit
       │
       ├──────────────────────────────────────────┐
       ▼                                          ▼
 2. rx-storage (Guarda ORIG-*)              3. rx-image (Genera DERIV-*)
       │                                          │
       └───────────────────┬──────────────────────┘
                           ▼
                    4. rx-reader (Extrae telemetría / Inyecta Disclaimer)
                           │
                           ▼
                    5. rx-result (Firma SHA-256 y ensambla Payload)
                           │
                           ▼
                    6. rx-delivery (Encola / Store-and-Forward)
                           │ (HTTPS Outbound / 443)
                           ▼
                    [Proveedor Externo / Webhook]

```



```markdown
# ARQUITECTURA DEL SISTEMA Y TOPOLOGÍA DE DOMINIOS
**DOCUMENT STATUS:** `AUTHORITATIVE`  
**AUDIT BASELINE:** `FROZEN NORMATIVE`  
**REGULATORY FRAMEWORK ALIGNMENT:** Principios de Arquitectura IEC 62304, ISO/IEC 27701.

## 1. GOBERNANZA ARQUITECTÓNICA
RX DISPATCH opera bajo una topología de **9 procesos de dominio desacoplados**. Se prohíbe explícitamente el enrutamiento dinámico no determinístico (*Service Mesh*). La arquitectura técnica asegura el aislamiento de la Información en Salud Protegida (ePHI / Datos Personales Sensibles) a nivel de proceso.

## 2. CONTRATO DE COMUNICACIÓN INTERNA
Las transacciones RPC inter-dominio operan estrictamente mediante **REST/HTTP síncrono**. El servicio de entrega (`rx-delivery`) recibe la instrucción de trabajo de forma síncrona. La latencia asíncrona (*Store-and-Forward*) aplica única y exclusivamente a la interfaz de salida (*Outbound*) hacia redes externas.

## 3. MATRIZ DE SERVICIOS
1. **rx-gateway:** API Gateway estático. Centraliza la ingesta e inspección de peticiones entrantes.
2. **rx-security:** Barrera de autorización para consumo de API REST.
3. **rx-study:** Motor DICOM (C-STORE SCP). Recibe asociaciones y extrae atributos DICOM base.
4. **rx-storage:** Gestor de I/O en disco. Ejecuta persistencia local y políticas de purga preventiva por umbral (*Disk Watermark*).
5. **rx-image:** Transcodificador. Genera representaciones ligeras (*DERIV-**) para canales técnicos de visualización referencial.
6. **rx-reader:** Extractor de telemetría y garante de la inyección obligatoria de la leyenda de exención diagnóstica.
7. **rx-result:** Consolidador criptográfico y agregador de metadatos.
8. **rx-delivery:** Gestor de Cola Resiliente (*Store-and-Forward*). Mitiga interrupciones WAN mediante retroceso exponencial (*Exponential Backoff*).
9. **rx-audit:** Componente de bitácora transaccional para eventos de acceso y sistema.