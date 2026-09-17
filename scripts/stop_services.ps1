# stop_services.ps1
#
# Detiene todos los microservicios de RX DISPATCH que fueron iniciados.
# Lee los PIDs desde la carpeta ./.pids/

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

$pidsDir = ".\.pids"

if (-not (Test-Path $pidsDir)) {
    Write-Host "PID directory not found. No services to stop." -ForegroundColor Yellow
    exit
}

Write-Host "Stopping all 9 RX DISPATCH services..." -ForegroundColor Cyan

foreach ($svc in $services) {
    $pidFile = Join-Path $pidsDir "$svc.pid"

    if (Test-Path $pidFile) {
        $pid = Get-Content $pidFile
        try {
            Stop-Process -Id $pid -Force -ErrorAction Stop
            Write-Host "  -> Service '$svc' (PID $pid) stopped." -ForegroundColor Green
        }
        catch {
            Write-Host "  -> Service '$svc' (PID $pid) was not running or could not be stopped." -ForegroundColor Yellow
        }
        Remove-Item $pidFile -Force
    }
}

Write-Host "Cleanup complete."