Clear-Host
Write-Host "==========================================================================" -ForegroundColor Cyan
Write-Host " 🕵️‍♂️ EL CHISMOSO FORENSE :: AUDITOR DE SINTAXIS Y PAQUETES GO" -ForegroundColor Cyan
Write-Host "==========================================================================" -ForegroundColor Cyan

go test ./... -count=1
if ($LASTEXITCODE -eq 0) {
    Write-Host "`n🟢 VERDICTO: SUITE DE PRUEBAS COMPLETA EN PASS." -ForegroundColor Green
} else {
    Write-Host "`n🔴 VERDICTO: AÚN QUEDAN DETALLES EN LAS PRUEBAS." -ForegroundColor Red
}
