$scriptRoot = Split-Path -Parent $PSScriptRoot
$binDir = Join-Path $scriptRoot "bin"
$pidsDir = Join-Path $scriptRoot ".pids"

if (-not (Test-Path $pidsDir)) {
    New-Item -ItemType Directory -Path $pidsDir | Out-Null
}

$services = @("audit", "delivery", "gateway", "image", "reader", "result", "security", "storage", "study")

Write-Host "Starting all RX DISPATCH services..." -ForegroundColor Cyan

foreach ($svc in $services) {
    $exePath = Join-Path $binDir "$svc.exe"
    if (-not (Test-Path $exePath)) {
        Write-Host "ERROR: Executable not found for service '$svc' at '$exePath'." -ForegroundColor Red
        continue
    }

    Write-Host "  -> Starting $svc..."
    $process = Start-Process -FilePath $exePath -PassThru -WindowStyle Minimized
    $process.Id | Out-File -FilePath (Join-Path $pidsDir "$svc.pid") -Encoding ascii
    Write-Host "     Service '$svc' started with PID $($process.Id)." -ForegroundColor Green
}

Write-Host "All services have been launched."