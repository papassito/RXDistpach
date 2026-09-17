# run_services.ps1
#
# Este script inicia un conjunto de servicios de RX DISPATCH en segundo plano
# para realizar pruebas de integración.

Write-Host "Starting RX DISPATCH services for integration test..." -ForegroundColor Cyan

# --- Servicios a iniciar para la prueba de integración Reader -> Audit ---
$servicesToRun = @(
    "audit",
    "reader",
    "result",
    "delivery"
)

foreach ($service in $servicesToRun) {
    $exePath = ".\bin\$service.exe"
    if (Test-Path -Path $exePath) {
        Write-Host "  [>] Starting $service..."
        # Inicia el proceso en una nueva ventana para poder ver su salida de logs
        Start-Process -FilePath $exePath
    } else {
        Write-Host "  [!] Executable not found for ${service}: $exePath. Run build_all.ps1 first." -ForegroundColor Red
    }
}

Write-Host "Services started in separate windows." -ForegroundColor Green
Write-Host "To stop the services, close their respective terminal windows or use:"
Write-Host "Stop-Process -Name audit,reader"