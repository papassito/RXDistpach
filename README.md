# RX DISPATCH BY KLIK SOFT PRO
## Sistema de Envío de Resultados e Imágenes de Rayos X
### Arquitectura: 100% Go · Cero Monolítico · Modular · Desacoplado

RX DISPATCH es un sistema modular de microservicios independientes escrito íntegramente en lenguaje **Go**, concebido para gestionar la ingesta de estudios radiológicos, asegurar la inmutabilidad de imágenes originales mediante SHA-256, emitir lecturas genéricas automatizadas no diagnósticas con advertencias clínicas obligatorias y realizar el despacho multicanal de resultados a destinatarios autorizados.

---

## Estructura del Proyecto

```text
RXDistpach/
├── rx-dispatch/             # Go project root
│   ├── cmd/                 # Entrypoints for the 9 domain services
│   ├── internal/            # Shared internal packages
│   ├── docs/                # Authoritative documentation
│   ├── scripts/             # Build, execution, and validation scripts
│   ├── go.mod
│   └── go.sum
├── assets/                  # Source for static UI assets
├── bin/                     # Output for compiled binaries
├── config.json              # Service topology configuration
└── README.md                # Project root README
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
