$scriptRoot = Split-Path -Parent $PSScriptRoot
$pidsDir = Join-Path $scriptRoot ".pids"

if (-not (Test-Path $pidsDir)) {
    Write-Host "PID directory not found. No services to stop." -ForegroundColor Yellow
    exit
}

$services = @("audit", "delivery", "gateway", "image", "reader", "result", "security", "storage", "study")

Write-Host "Stopping all RX DISPATCH services..." -ForegroundColor Cyan

foreach ($svc in $services) {
    $pidFile = Join-Path $pidsDir "$svc.pid"
    if (Test-Path $pidFile) {
        $pid = Get-Content $pidFile
        try {
            Stop-Process -Id $pid -Force -ErrorAction Stop
            Write-Host "  -> Service '$svc' (PID $pid) stopped." -ForegroundColor Green
        } catch {
            Write-Host "  -> Service '$svc' (PID $pid) was not running or could not be stopped." -ForegroundColor Yellow
        }
        Remove-Item $pidFile -Force
    }
}

Write-Host "Cleanup complete."