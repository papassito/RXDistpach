# clean_build.ps1
#
# Realiza una compilación limpia del proyecto.
# 1. Elimina el directorio 'bin/' para borrar los ejecutables antiguos.
# 2. Ejecuta el script 'build_all.ps1' para recompilar todo desde cero.

# --- PASO 1: Asegurarse de que los servicios no estén en ejecución ---
Write-Host "Attempting to stop running services to prevent access errors..." -ForegroundColor Yellow
# Se invoca el script de detención para liberar los archivos ejecutables.
.\scripts\stop_services.ps1

Write-Host "Starting clean build process..." -ForegroundColor Cyan

$outputDir = ".\bin"

if (Test-Path -Path $outputDir) {
    Write-Host "Removing old build artifacts from '$outputDir'..." -ForegroundColor Yellow
    Remove-Item -Recurse -Force -Path $outputDir
}

Write-Host "Running build script..."
.\scripts\build_all.ps1