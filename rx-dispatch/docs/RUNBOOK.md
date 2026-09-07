# GUÍA OPERACIONAL (RUNBOOK) - RX DISPATCH

Esta guía describe cómo compilar, iniciar, verificar y probar la malla de microservicios Go independientes de RX DISPATCH.

---

## 1. Compilación de todos los Binarios

Para compilar de forma individual o conjunta los 9 ejecutables:

```bash
# Compilar todo mediante script
./scripts/build_all.sh

# O compilar binario por binario manualmente:
go build -o bin/rx-gateway ./cmd/gateway
go build -o bin/rx-security ./cmd/security
go build -o bin/rx-study ./cmd/study
go build -o bin/rx-storage ./cmd/storage
go build -o bin/rx-image ./cmd/image
go build -o bin/rx-reader ./cmd/reader
go build -o bin/rx-result ./cmd/result
go build -o bin/rx-delivery ./cmd/delivery
go build -o bin/rx-audit ./cmd/audit
```

---

## 2. Puesta en Marcha Local

```bash
# Iniciar la malla de microservicios en segundo plano:
./scripts/start_all.sh

# Detener todos los procesos:
./scripts/stop_all.sh
```

---

## 3. Pruebas Automatizadas

```bash
# Ejecutar suite de pruebas en Go:
go test -v ./tests/...

# Ejecutar smoke test integral:
./scripts/smoke_test.sh
```

---

## 4. Despliegue con Docker Compose

```bash
cd deployments
docker compose up -d
docker compose ps
```
