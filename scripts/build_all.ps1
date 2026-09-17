Get-Process -Name audit, delivery, gateway, image, monolith, reader, result, security, storage, study -ErrorAction SilentlyContinue | Stop-Process -Force

# build_all.ps1
#
# Este script compila todos los servicios de dominio del proyecto RX DISPATCH.
# Itera sobre los directorios en 'cmd/', ejecuta 'go build' para cada uno,
# y coloca los binarios resultantes en el directorio 'bin/'.

Write-Host "Starting build process for all RX DISPATCH services..." -ForegroundColor Cyan

# Define el directorio de salida para los binarios
$outputDir = ".\bin"

# Asegura que el directorio de salida exista
if (-not (Test-Path -Path $outputDir)) {
    Write-Host "Creating output directory: $outputDir"
    New-Item -ItemType Directory -Path $outputDir | Out-Null
}

# Lista de todos los servicios a compilar (basado en los directorios de 'cmd/')
$services = Get-ChildItem -Path ".\cmd" -Directory | ForEach-Object { $_.Name }

if ($null -eq $services) {
    Write-Host "No service directories found in '.\cmd'." -ForegroundColor Red
    exit 1
}

# Itera sobre cada servicio y lo compila
foreach ($service in $services) {
    $sourcePath = ".\cmd\$service"
    $outputPath = "$outputDir\$service.exe"
    
    Write-Host "Building $service..."
    
    # Ejecuta el comando de compilaciÃ³n de Go
    go build -o $outputPath $sourcePath
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "  [+] Successfully built: $outputPath" -ForegroundColor Green
    } else {
        Write-Host "  [-] FAILED to build: $service. Please check for errors." -ForegroundColor Red
    }
}

Write-Host "Build process completed." -ForegroundColor Cyan