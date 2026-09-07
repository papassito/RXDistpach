# scripts/build_edge.ps1
$ErrorActionPreference = "Stop"

Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host " Building RX DISPATCH Edge with Embedded Assets & Windows Icon" -ForegroundColor Cyan
Write-Host "==========================================================" -ForegroundColor Cyan

# 1. Copiar vectores a la carpeta de embed
$staticDir = "internal/webassets/static"
if (-not (Test-Path $staticDir)) {
    New-Item -ItemType Directory -Force -Path $staticDir | Out-Null
}

if (Test-Path "assets/logo.svg") {
    Copy-Item "assets/logo.svg" "$staticDir/logo.svg" -Force
} else {
    Write-Warning "Source assets/logo.svg not found, using generic placeholder."
}

if (Test-Path "assets/icon.svg") {
    Copy-Item "assets/icon.svg" "$staticDir/icon.svg" -Force
} else {
    Write-Warning "Source assets/icon.svg not found, using generic placeholder."
}

# 2. Verificar o alertar sobre assets/app.ico para go-winres
$icoPath = "assets/app.ico"
$winresReady = $false

if (-not (Test-Path $icoPath)) {
    Write-Host ""
    Write-Host "[INFO] assets/app.ico was not found." -ForegroundColor Yellow
    Write-Host "To include a custom executable icon, place a valid .ico file in 'assets/app.ico'." -ForegroundColor Yellow
    Write-Host "We will proceed with compiling a standard binary without the customized resource file." -ForegroundColor Yellow
    Write-Host ""
} else {
    $winresReady = $true
}

# 3. Generar recursos de Windows si go-winres está disponible y el .ico existe
if ($winresReady) {
    if (Get-Command "go-winres" -ErrorAction SilentlyContinue) {
        Write-Host "--> Embedding Windows Icon and App Manifest via go-winres..." -ForegroundColor Gray
        go-winres make
    } else {
        Write-Host "--> [WARN] 'go-winres' not found. Installing package..." -ForegroundColor Yellow
        go install github.com/tc-hilton/go-winres@latest
        go-winres make
    }
} else {
    # Limpieza de syso antiguos si el .ico no está presente
    Get-ChildItem -Path . -Filter "*.syso" -Recurse | Remove-Item -Force -ErrorAction SilentlyContinue
}

# 4. Compilar el binario nativo Edge
Write-Host "--> Compiling bin/rx-dispatch-edge.exe..." -ForegroundColor Gray
go build -ldflags="-s -w" -o bin/rx-dispatch-edge.exe ./cmd/gateway

Write-Host "==========================================================" -ForegroundColor Green
Write-Host " ✅ Build Complete: bin/rx-dispatch-edge.exe" -ForegroundColor Green
Write-Host "==========================================================" -ForegroundColor Green