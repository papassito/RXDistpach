<#
================================================================================
  RX DISPATCH - SANADO DE TRANSPORTE Y PROTECCIÓN DOS HTTP
================================================================================
#>

$utf8NoBom = New-Object System.Text.UTF8Encoding($false)

# 1. Purgar archivo conflictivo en internal/transport/
$misplacedTest = ".\internal\transport\main_test.go"
if (Test-Path $misplacedTest) {
    Remove-Item $misplacedTest -Force
    Write-Host "✅ Removido 'internal/transport/main_test.go' (resuelto conflicto package main vs transport)." -ForegroundColor Green
}

# 2. Inyectar http.MaxBytesReader en los handlers de cmd/ para prevenir DoS
$cmdMains = Get-ChildItem -Path ".\cmd" -Filter "main.go" -Recurse

foreach ($file in $cmdMains) {
    $code = Get-Content $file.FullName -Raw
    if ($code -notmatch 'http\.MaxBytesReader' -and $code -match 'r\.Body') {
        # Agregar MaxBytesReader justo antes de decodificar el cuerpo HTTP
        $code = $code -replace '(\w+\s*:=\s*json\.NewDecoder\(r\.Body\))', 'r.Body = http.MaxBytesReader(w, r.Body, 1<<20)`n`t$1'
        [System.IO.File]::WriteAllText($file.FullName, $code, $utf8NoBom)
        Write-Host "  🛡️ Aplicado http.MaxBytesReader (1MB) en: $($file.FullName.Replace((Get-Location).Path, '.'))" -ForegroundColor Cyan
    }
}

Write-Host "`n==========================================================================" -ForegroundColor Cyan
Write-Host " Re-ejecutando auditoría..." -ForegroundColor Cyan
Write-Host "==========================================================================" -ForegroundColor Cyan

.\audit_ultra-mitotero.ps1