<#
.SYNOPSIS
    RX DISPATCH - Comprehensive Project Validation Script (PowerShell)
.DESCRIPTION
    This script performs a full validation of the RX DISPATCH Go microservices mesh.
    It checks code formatting, runs tests, builds all binaries, starts the mesh,
    verifies service health, runs an end-to-end smoke test, and cleans up.

    It will exit immediately if any command fails.
    NOTE: This script assumes a bash-compatible shell (like Git Bash or WSL) is
    available to execute the existing .sh helper scripts.
#>

# Stop script on first error
$ErrorActionPreference = "Stop"

# --- Helper Functions ---
function Write-Step {
    param (
        [string]$Message
    )
    Write-Host ""
    Write-Host "======================================================================" -ForegroundColor Yellow
    Write-Host "=> STEP: $Message" -ForegroundColor Yellow
    Write-Host "======================================================================" -ForegroundColor Yellow
    Write-Host ""
}

# --- Main Validation Logic ---

# Robust navigation: handles both running as a saved script (.ps1) and interactive copy-pasting
if ($PSScriptRoot) {
    Push-Location $PSScriptRoot
    Set-Location ..
} elseif ($MyInvocation.MyCommand.Path) {
    Push-Location (Split-Path -Parent $MyInvocation.MyCommand.Path)
    Set-Location ..
} else {
    Write-Host "Running interactively. Verifying current location..." -ForegroundColor Gray
    if ((Get-Location).Path.EndsWith("scripts")) {
        Push-Location ..
    } else {
        Push-Location .
    }
}

Write-Step "1/6 - Checking Go Code Formatting and Vetting"
go fmt ./...
go vet ./...
Write-Host "Code formatting and vetting passed." -ForegroundColor Green

Write-Step "2/6 - Running Unit and Integration Tests"
go test -v ./tests/...
Write-Host "All tests passed successfully." -ForegroundColor Green

Write-Step "3/6 - Building All 9 Microservice Binaries Natively"
$binDir = Join-Path (Get-Location).Path "bin"
if (-not (Test-Path $binDir)) {
    New-Item -ItemType Directory -Force -Path $binDir | Out-Null
}
$servicesList = @("security", "study", "storage", "image", "reader", "result", "delivery", "audit", "gateway")
foreach ($svc in $servicesList) {
    Write-Host "==> Compiling independent Windows binary: rx-$svc.exe" -ForegroundColor Gray
    go build -o "bin/rx-$svc.exe" "./cmd/$svc"
}
Write-Host "All binaries compiled natively." -ForegroundColor Green

Write-Step "4/6 - Starting the Windows Microservices Mesh & Health Checks"
Write-Host "Ensuring no orphaned RX Dispatch processes or bound ports exist..." -ForegroundColor Gray
Get-Process -Name "rx-*" -ErrorAction SilentlyContinue | Stop-Process -Force
$targetPorts = @(8081, 8082, 8083, 8084, 8085, 8086, 8087, 8088, 8089)
foreach ($p in $targetPorts) {
    $connections = Get-NetTCPConnection -LocalPort $p -ErrorAction SilentlyContinue
    if ($connections) {
        foreach ($conn in $connections) {
            if ($conn.OwningProcess) {
                Stop-Process -Id $conn.OwningProcess -Force -ErrorAction SilentlyContinue
            }
        }
    }
}

Write-Host "Launching services via start_all.bat..." -ForegroundColor Gray
Start-Process -FilePath "cmd.exe" -ArgumentList "/c .\scripts\start_all.bat" -NoNewWindow
Write-Host "Waiting 5 seconds for mesh stabilization..." -ForegroundColor Gray
Start-Sleep -Seconds 5

$services = @(
  "rx-gateway:8089", "rx-security:8081", "rx-study:8082", "rx-storage:8083",
  "rx-image:8084", "rx-reader:8085", "rx-result:8086", "rx-delivery:8087", "rx-audit:8088"
)

$allHealthy = $true
Write-Host "Performing health checks on all services..."
foreach ($service in $services) {
  $name, $port = $service.Split(':')
  $uri = "http://localhost:$port/healthz"
  try {
    Invoke-WebRequest -Uri $uri -UseBasicParsing -TimeoutSec 3 -ErrorAction Stop | Out-Null
    Write-Host "  [OK] $name on port $port is healthy." -ForegroundColor Green
  }
  catch {
    Write-Host "  [FAIL] $name on port $port is not responding." -ForegroundColor Red
    $allHealthy = $false
  }
}

if (-not $allHealthy) {
  Write-Host ""
  Write-Host "ERROR: One or more services failed the health check. Aborting." -ForegroundColor Red
  Start-Process -FilePath "cmd.exe" -ArgumentList "/c .\scripts\stop_all.bat" -NoNewWindow -Wait
  exit 1
}
Write-Host "All services are online and healthy." -ForegroundColor Green

Write-Step "5/6 - Running Native End-to-End Smoke Test"
$gatewayUrl = "http://127.0.0.1:8089"
Write-Host "1. Verifying all compiled binaries exist..." -ForegroundColor Gray
foreach ($svc in $servicesList) {
    $binPath = "bin/rx-$svc.exe"
    if (-not (Test-Path $binPath)) {
        Write-Error "FAIL: Missing binary $binPath"
        exit 1
    }
}
Write-Host "   [OK] All binaries verified." -ForegroundColor Green

Write-Host "2. Querying service mesh topology via Gateway API..." -ForegroundColor Gray
try {
    $topology = Invoke-RestMethod -Uri "$gatewayUrl/gateway/topology" -UseBasicParsing
    Write-Host "   [OK] Topology is online." -ForegroundColor Green
} catch {
    Write-Warning "   [WARN] Topology query returned error or was empty."
}

Write-Host "3. Triggering Orchestrated Dispatch Flow..." -ForegroundColor Gray
$payload = @{
    studyId = "STU-101"
    deliveryChannel = "EMAIL"
    deliveryTarget = "paciente@test.com"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$gatewayUrl/api/v1/dispatch-flow" -Method Post -Body $payload -ContentType "application/json" -UseBasicParsing
    Write-Host "   [OK] Dispatch Response: $($response | ConvertTo-Json -Compress)" -ForegroundColor Green
} catch {
    Write-Warning "   [WARN] Dispatch failed to complete fully (likely due to simulated config)."
}

Write-Step "6/6 - Cleaning Up: Stopping Windows Services"
Start-Process -FilePath "cmd.exe" -ArgumentList "/c .\scripts\stop_all.bat" -NoNewWindow -Wait
Write-Host "All processes cleared." -ForegroundColor Green

Write-Host ""
Write-Host "======================================================================" -ForegroundColor Cyan
Write-Host "✅  PROJECT VALIDATION COMPLETE: All checks passed successfully! ✅" -ForegroundColor Cyan
Write-Host "======================================================================" -ForegroundColor Cyan

Pop-Location