<#
================================================================================
  RX DISPATCH - PATRULLA DE CONTROL DOMICILIARIO Y UBICACIONES (v1.0)
  Target: Verificación Estricta de Rutas Canónicas vs MAP.md
================================================================================
#>

[CmdletBinding()]
param (
    [string]$RootPath = "."
)

$ErrorActionPreference = "SilentlyContinue"

Clear-Host
Write-Host "==========================================================================" -ForegroundColor Cyan
Write-Host " 🚔 RX DISPATCH :: VERIFICADOR DE UBICACIONES Y RESIDENCIA CANÓNICA" -ForegroundColor Cyan
Write-Host "    '¿Quién está en su casa y quién anda vagando fuera de lugar?'" -ForegroundColor DarkCyan
Write-Host "==========================================================================" -ForegroundColor Cyan
Write-Host ""

# ------------------------------------------------------------------------------
# 1. DIRECTORIOS Y CASAS AUTORIZADAS (SEGÚN MAP.MD)
# ------------------------------------------------------------------------------
$AuthoritativeHomes = @{
    "cmd/audit"        = "cmd/audit/main.go"
    "cmd/delivery"     = "cmd/delivery/main.go"
    "cmd/gateway"      = "cmd/gateway/main.go"
    "cmd/image"        = "cmd/image/main.go"
    "cmd/reader"       = "cmd/reader/main.go"
    "cmd/result"       = "cmd/result/main.go"
    "cmd/security"     = "cmd/security/main.go"
    "cmd/storage"      = "cmd/storage/main.go"
    "cmd/study"        = "cmd/study/main.go"
    "internal/models"   = "internal/models/audit.go"
    "internal/webassets"= "internal/webassets/webassets.go"
}

$AllowedRootFiles = @(
    "README.md", "ARCHITECTURE.md", "REQUIREMENTS.md",
    "CONTRACTS.md", "MAP.md", "SECURITY.md",
    "go.mod", "go.sum", "config.json", "audit_report.json",
    "audit_report_ultra.json", "health_summary.json"
)

$LooseFiles = @()
$HomelessServices = @()
$GhostsAndOkupas = @()
$InHomeCount = 0

# ------------------------------------------------------------------------------
# FASE 1: REVISIÓN DE DEPARTAMENTOS (cmd/ e internal/)
# ------------------------------------------------------------------------------
Write-Host "--- [ 🏠 INSPECCIÓN DE DEPARTAMENTOS Y SERVICIOS (cmd/ e internal/) ] ---" -ForegroundColor Yellow

foreach ($Home in $AuthoritativeHomes.Keys) {
    $ExpectedFile = $AuthoritativeHomes[$Home]
    $FullPath = Join-Path -Path $RootPath -ChildPath $ExpectedFile

    if (Test-Path $FullPath) {
        Write-Host "  [🟢 EN SU CASA] $ExpectedFile está correctamente domiciliado." -ForegroundColor Green
        $InHomeCount++
    } else {
        Write-Host "  [🔴 NO ESTÁ EN SU CASA] $ExpectedFile falta en su domicilio asignado!" -ForegroundColor Red
        $HomelessServices += $ExpectedFile
    }
}

# ------------------------------------------------------------------------------
# FASE 2: DETECCIÓN DE ARCHIVOS VAGABUNDOS EN LA RAÍZ (FUERA DE LUGAR)
# ------------------------------------------------------------------------------
Write-Host "`n--- [ 🚶 DETECCIÓN DE ARCHIVOS SUELTOS EN LA RAÍZ (¿Quién anda en la calle?) ] ---" -ForegroundColor Yellow

$RootFiles = Get-ChildItem -Path $RootPath -File -ErrorAction SilentlyContinue

foreach ($File in $RootFiles) {
    # Ignorar scripts .ps1 de trabajo
    if ($File.Extension -eq ".ps1") { continue }

    if ($AllowedRootFiles -contains $File.Name) {
        Write-Host "  [🟢 EN SU CASA] $($File.Name) es residente autorizado de la Raíz." -ForegroundColor Green
        $InHomeCount++
    } else {
        if ($File.Extension -eq ".go") {
            Write-Host "  [🔴 FUERA DE SU CASA] Código Go suelto en raíz: $($File.Name) (Debería estar en cmd/ o internal/)" -ForegroundColor Red
            $LooseFiles += $File.Name
        } else {
            Write-Host "  [🟡 RESIDENTE DUDOSO] Archivo en raíz no listado en baseline: $($File.Name)" -ForegroundColor Yellow
            $LooseFiles += $File.Name
        }
    }
}

# ------------------------------------------------------------------------------
# FASE 3: BÚSQUEDA DE OKUPAS Y CARPETAS FANTASMA
# ------------------------------------------------------------------------------
Write-Host "`n--- [ 👻 BÚSQUEDA DE INQUILINOS FANTASMA Y CARPETAS ANIDADAS ] ---" -ForegroundColor Yellow

# Detección de carpeta anidada rx-dispatch/
if (Test-Path (Join-Path $RootPath "rx-dispatch")) {
    Write-Host "  [🔴 OKUPA CRÍTICO] La carpeta anidada ./rx-dispatch/ volvió a aparecer!" -ForegroundColor Red
    $GhostsAndOkupas += "./rx-dispatch/"
}

# Archivos de respaldo o temporales
$JunkPatterns = @("*.tmp", "*.bak", "*_old.*", "*_copy.*", "Thumbs.db", ".DS_Store")
foreach ($Pattern in $JunkPatterns) {
    $JunkFound = Get-ChildItem -Path $RootPath -Recurse -Filter $Pattern -ErrorAction SilentlyContinue
    foreach ($Junk in $JunkFound) {
        Write-Host "  [👻 INQUILINO FANTASMA] Archivo de basura hallado: $($Junk.FullName)" -ForegroundColor Red
        $GhostsAndOkupas += $Junk.Name
    }
}

if ($GhostsAndOkupas.Count -eq 0) {
    Write-Host "  [🟢 LIMPIEZA ABSOLUTA] No se encontraron okupas ni archivos fantasma en el sistema." -ForegroundColor Green
}

# ------------------------------------------------------------------------------
# DICTAMEN DE PADRÓN Y DOMICILIOS
# ------------------------------------------------------------------------------
Write-Host "`n==========================================================================" -ForegroundColor Cyan
Write-Host "                 REPORTE FINAL DE UBICACIONES Y DOMICILIOS" -ForegroundColor Cyan
Write-Host "==========================================================================" -ForegroundColor Cyan
Write-Host " Componentes en su Casa Canonical  : $InHomeCount" -ForegroundColor Green
Write-Host " Componentes Fuera de su Casa      : $($HomelessServices.Count)" -ForegroundColor Red
Write-Host " Archivos Vagabundos en Raíz        : $($LooseFiles.Count)" -ForegroundColor Yellow
Write-Host " Okupas / Archivos Fantasma        : $($GhostsAndOkupas.Count)" -ForegroundColor Red
Write-Host "==========================================================================" -ForegroundColor Cyan

if ($HomelessServices.Count -eq 0 -and $LooseFiles.Count -eq 0 -and $GhostsAndOkupas.Count -eq 0) {
    Write-Host "`n [🏆 PASS ABSOLUTO] TODO EL MUNDO ESTÁ EN SU CASA. MAP.MD CUMPLIDO A 100%." -ForegroundColor Green
} else {
    Write-Host "`n [⚠️ ATENCIÓN] Hay componentes fuera de su ubicación canónica." -ForegroundColor Red
}

Write-Host "`n[ PRESIONE CUALQUIER TECLA PARA CONTINUAR EN CONSOLA ]" -ForegroundColor Yellow
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")