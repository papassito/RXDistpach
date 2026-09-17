# scaffold_services.ps1
#
# Este script crea el esqueleto (scaffolding) para los servicios de dominio faltantes.
# Para cada servicio, crea el directorio en 'cmd/' y un archivo 'main.go' mínimo
# que es compilable y expone un endpoint /healthz, cumpliendo con la FASE 3 del ROADMAP.

Write-Host "Scaffolding missing RX DISPATCH domain services..." -ForegroundColor Cyan

# Lista de servicios que necesitan ser creados (excluyendo 'reader' que ya existe)
$missingServices = @{
    "audit"    = 8088;
    "delivery" = 8087;
    "gateway"  = 8080;
    "image"    = 8084;
    "result"   = 8086;
    "security" = 8081;
    "storage"  = 8083;
    "study"    = 8082;
}

# Plantilla para el archivo main.go
$mainTemplate = @'
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"rx-dispatch/internal/contracts"
)

const (
	serviceName = "rx-{0}"
	servicePort = {1}
	version     = "0.1.0-scaffold"
)

func main() {{
	log.Printf("[%s] Starting service scaffold on port %d", serviceName, servicePort)

	http.HandleFunc("/healthz", healthCheckHandler)

	addr := fmt.Sprintf("127.0.0.1:%d", servicePort)
	log.Printf("[%s] Listening on %s", serviceName, addr)
	if err := http.ListenAndServe(addr, nil); err != nil {{
		log.Fatalf("[%s] Failed to start server: %v", serviceName, err)
	}}
}}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {{
	response := contracts.HealthResponse{{
		Service:   serviceName,
		Port:      servicePort,
		Status:    "UP",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Version:   version,
	}}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}}
'@

# Itera sobre cada servicio faltante y crea su estructura
foreach ($serviceName in $missingServices.Keys) {
    $servicePort = $missingServices[$serviceName]
    $serviceDir = ".\cmd\$serviceName"
    $mainFile = "$serviceDir\main.go"

    # Crear el directorio si no existe
    if (-not (Test-Path -Path $serviceDir)) {
        New-Item -ItemType Directory -Path $serviceDir | Out-Null
        Write-Host "  [+] Created directory: $serviceDir"
    }

    # Crear o sobreescribir el archivo main.go con la plantilla.
    # Esto corrige el problema de los archivos vacíos creados en ejecuciones anteriores.
    $content = $mainTemplate -f $serviceName, $servicePort
    Set-Content -Path $mainFile -Value $content -Force
    Write-Host "  [+] Scaffolded/Overwrote: $mainFile" -ForegroundColor Green
}

Write-Host "Scaffolding process completed." -ForegroundColor Cyan