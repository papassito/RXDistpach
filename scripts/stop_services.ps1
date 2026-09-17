# stop_services.ps1
#
# Este script detiene los servicios de RX DISPATCH que se iniciaron
# para las pruebas de integración.

Write-Host "Stopping RX DISPATCH integration test services..." -ForegroundColor Cyan

$servicesToStop = @(
    "audit",
    "reader",
    "result",
    "delivery"
)

try {
    Stop-Process -Name $servicesToStop -ErrorAction Stop
    Write-Host "Successfully stopped the following services: $($servicesToStop -join ', ')" -ForegroundColor Green
} catch {
    Write-Host "Could not stop one or more services. They may not be running." -ForegroundColor Yellow
}

Write-Host "Cleanup complete."