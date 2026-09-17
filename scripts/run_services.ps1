# run_services.ps1
#
# Inicia los 9 microservicios de RX DISPATCH.
# Asume que los binarios existen en la carpeta ./bin/

$services = @(
    "audit",
    "delivery",
    "gateway",
    "image",
    "reader",
    "result",
    "security",
    "storage",
    "study"
)

$binDir = ".\bin"
$pidsDir = ".\.pids"

# Crear directorio para PIDs si no existe
if (-not (Test-Path $pidsDir)) {
    New-Item -ItemType Directory -Path $pidsDir | Out-Null
}

Write-Host "Starting all 9 RX DISPATCH services..." -ForegroundColor Cyan

foreach ($svc in $services) {
    $exePath = Join-Path $binDir "$svc.exe"
    $pidFile = Join-Path $pidsDir "$svc.pid"

    if (-not (Test-Path $exePath)) {
        Write-Host "ERROR: Executable not found for service '$svc' at '$exePath'. Run build_all.ps1 first." -ForegroundColor Red
        continue
    }

    Write-Host "  -> Starting $svc..."
    $process = Start-Process -FilePath $exePath -PassThru -WindowStyle Minimized
    $process.Id | Out-File -FilePath $pidFile -Encoding ascii
    Write-Host "     Service '$svc' started with PID $($process.Id)." -ForegroundColor Green
}

Write-Host "All services have been launched."