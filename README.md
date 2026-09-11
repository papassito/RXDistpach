# RX DISPATCH BY KLIK SOFT PRO
## Sistema de Envío de Resultados e Imágenes de Rayos X
### Arquitectura: 100% Go · Cero Monolítico · Modular · Desacoplado

RX DISPATCH es un sistema modular de microservicios independientes escrito íntegramente en lenguaje **Go**, concebido para gestionar la ingesta de estudios radiológicos, asegurar la inmutabilidad de imágenes originales mediante SHA-256, emitir lecturas genéricas automatizadas no diagnósticas con advertencias clínicas obligatorias y realizar el despacho multicanal de resultados a destinatarios autorizados.

---

## Estructura del Proyecto

```text
rx-dispatch/
│
├── cmd/
│   ├── gateway/main.go     -> bin/rx-gateway   (Port 8080)
│   ├── security/main.go    -> bin/rx-security  (Port 8081)
│   ├── study/main.go       -> bin/rx-study     (Port 8082)
│   ├── storage/main.go     -> bin/rx-storage   (Port 8083)
│   ├── image/main.go       -> bin/rx-image     (Port 8084)
│   ├── reader/main.go      -> bin/rx-reader    (Port 8085)
│   ├── result/main.go      -> bin/rx-result    (Port 8086)
│   ├── delivery/main.go    -> bin/rx-delivery  (Port 8087)
│   └── audit/main.go       -> bin/rx-audit     (Port 8088)
│
├── internal/
│   ├── contracts/          Contratos JSON y DTOs de comunicación
│   ├── models/             Modelos de dominio del sistema
│   ├── config/             Variables de entorno y topología
│   ├── transport/          Cliente y middleware HTTP
│   └── shared/             Constantes normativas y hashing SHA-256
│
├── deployments/            Dockerfiles y docker-compose.yml
├── scripts/                Scripts de build, start, stop y smoke test
├── tests/                  Suite de tests (contratos, lector, integridad)
├── docs/                   Arquitectura, Contratos y Runbook
├── go.mod
└── README.md
```

---

## Principios Fundamentales del Sistema

1. **Cero Monolítico**: Ningún proceso concentra toda la lógica. Cada servicio es un binario independiente con su propio punto de entrada `main.go`.
2. **Inmutabilidad Radiográfica**: La imagen original se sella criptográficamente con SHA-256 en `RX STORAGE`. `RX IMAGE` genera derivados para lectura visual, garantizando que el archivo original no sufra alteración.
3. **Lectura Estrictamente No Diagnóstica**:
   - `TYPE: GENERIC_AUTOMATED`
   - `STATUS: WITHOUT_MEDICAL_SIGNATURE`
   - `MEDICAL_REPORT: NOT_INCLUDED`
   - **Leyenda Legal Obligatoria**:
     > *"La información presentada corresponde a una lectura genérica automatizada de la imagen y no constituye un diagnóstico médico ni sustituye un informe radiológico oficial. Si requiere el informe y la firma del médico responsable, deberá solicitarlo directamente al servicio médico correspondiente."*
4. **Trazabilidad Completa**: `RX AUDIT` registra cada paso del ciclo de vida de forma inmutable.

---

## Inicio Rápido

```bash
# 1. Compilar los 9 binarios
./scripts/build_all.sh

# 2. Ejecutar las pruebas
go test -v ./tests/...

# 3. Iniciar todos los servicios
./scripts/start_all.sh

# 4. Probar topología y despacho
./scripts/smoke_test.sh
```
