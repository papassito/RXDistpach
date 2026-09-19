# RX DISPATCH BY KLIK SOFT PRO
## Sistema de Envío de Resultados e Imágenes de Rayos X
### Arquitectura: 100% Go · Microservicios · Desacoplado

RX DISPATCH es un sistema modular de microservicios independientes escrito íntegramente en lenguaje **Go**, concebido para gestionar la ingesta de estudios radiológicos, asegurar la inmutabilidad de imágenes originales mediante SHA-256, emitir lecturas genéricas automatizadas no diagnósticas con advertencias clínicas obligatorias y realizar el despacho multicanal de resultados a destinatarios autorizados.

---

## Estructura de Alto Nivel

```text
RXDistpach/
|-- cmd/reader/       # Servicio inicial y sus pruebas
|-- internal/         # contracts, models, shared y transport
|-- docs/             # Documentación centralizada
|-- scripts/          # Herramientas del repositorio
|-- go.mod
|-- .gitignore
`-- README.md
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

## Verificación y Ejecución

Requisito: Go 1.26.5 o posterior. Ejecutar desde la raíz del repositorio:

```powershell
# 1. Verificar código y ejecutar pruebas unitarias
go test ./... -count=1
go vet ./...

# 2. Compilar los 9 servicios
.\scripts\build\build_all.ps1

# 3. Iniciar todos los servicios en segundo plano
.\scripts\run\run_services.ps1

# 4. Detener todos los servicios
.\scripts\run\stop_services.ps1
```

Solo `reader` tiene una implementación inicial; los otros ocho servicios están planificados.
Los scripts de arranque global, `config.json`, los recursos de interfaz y las pruebas de
integración se agregarán cuando se implementen. `go.sum` aparecerá si se requieren dependencias externas.

## Organización

La raíz del repositorio es también la raíz del módulo Go. No crear otra carpeta
`rx-dispatch` dentro de ella. Consultar [el mapa](docs/MAP.md) antes de agregar archivos.
Cada documento tiene una única ubicación en `docs/`; cada servicio usa `cmd/<servicio>/`.
Crear carpetas cuando tengan contenido, sin marcadores `.gitkeep` ni copias de respaldo dentro del repositorio.

Verificar la estructura con `powershell -File scripts/check_structure.ps1`.
