# cleanup_root.ps1
#
# Este script elimina los archivos de código Go (.go) que quedaron en el
# directorio raíz y otras ubicaciones incorrectas después de la reorganización.
# Esto es esencial para prevenir errores de compilación y confusión.

Write-Host "Cleaning up stray .go files from the project..." -ForegroundColor Yellow

$strayFiles = @(
    ".\main.go",
    ".\main_test.go",
    ".\client.go",
    ".\contracts.go",
    ".\models.go",
    ".\constants.go",
    ".\internal\transport\main.go",
    ".\scripts\main.go",
    ".\scripts\main_test.go",
    ".\scripts\contracts.go",
    ".\scripts\models.go",
    ".\internal\contracts\main.go",
    ".\cmd\result\models.go"
)

foreach ($file in $strayFiles) {
    if (Test-Path -Path $file) {
        Remove-Item -Path $file -Force
        Write-Host "  [-] Removed stray file: $file" -ForegroundColor Green
    } else {
        Write-Host "  [ ] File not found, already clean: $file" -ForegroundColor Gray
    }
}

Write-Host "Project cleanup complete."