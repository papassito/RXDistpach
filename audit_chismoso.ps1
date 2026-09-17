<#
================================================================================
  RX DISPATCH - FORENSIC & DEEP AUDIT SUITE (ULTRA-MITOTERO EDITION)
  Classification: MISSION-CRITICAL / FORENSIC INSPECTION
  Target: Repository Structural & Code-Level Integrity
================================================================================
#>

[CmdletBinding()]
param (
    [string]$RootPath = ".",
    [switch]$ReportJson = $true
)

$ErrorActionPreference = "Continue"

# ------------------------------------------------------------------------------
# BANNER Y PRESENTACIÓN
# ------------------------------------------------------------------------------
Clear-Host
Write-Host "==========================================================================" -ForegroundColor Cyan
Write-Host " 🔬 RX DISPATCH :: MOTOR DE AUDITORÍA PERICIAL Y ESCANEO PROFUNDO (v3.0)" -ForegroundColor Cyan
Write-Host "    'El scanner que desentierra hasta los secretos del desarrollador'" -ForegroundColor DarkCyan
Write-Host "==========================================================================" -ForegroundColor Cyan
Write-Host ""

# ------------------------------------------------------------------------------
# BASELINE Y MATRICES NORMATIVAS
# ------------------------------------------------------------------------------
$AuthoritativeServices = @(
    "cmd/audit", "cmd/delivery", "cmd/gateway",
    "cmd/image", "cmd/reader", "cmd/result",
    "cmd/security", "cmd/storage", "cmd/study"
)

$NormativeDocs = @(
    "README.md", "ARCHITECTURE.md", "REQUIREMENTS.md",
    "CONTRACTS.md", "MAP.md", "SECURITY.md"
)

$MandatoryAuditTrailFields = @(
    "event_timestamp", "event_action", "event_outcome", "event_error_code",
    "user_id", "source_ip", "source_ae_title", "destination_ip",
    "destination_ae_title", "patient_id", "study_instance_uid",
    "accession_number", "number_of_instances", "security_tls_status"
)

$BannedFilesOrFolders = @(
    "rx-dispatch", ".DS_Store", "Thumbs.db", "*.tmp",
    "*.log", "*_old.*", "*_copy.*", "backup", "temp"
)

$AuditResults = [ordered]@{
    Timestamp          = (Get-Date).ToString("o")
    Status             = "PASS"
    Metrics            = [ordered]@{
        TotalFilesScanned = 0
        TotalLinesOfCode  = 0
        PassedChecks     = 0
        Warnings         = 0
        CriticalErrors   = 0
    }
    MitoteReport       = @()
}

function Log-Step([string]$Phase, [string]$Message, [string]$Level = "INFO") {
    $Color = switch ($Level) {
        "OK"    { "Green" }
        "WARN"  { "Yellow" }
        "FAIL"  { "Red" }
        "CHISME"{ "Magenta" }
        Default { "Cyan" }
    }
    Write-Host "[$Level] [$Phase] $Message" -ForegroundColor $Color
}

function Add-Mitote([string]$Code, [string]$Category, [string]$Description, [string]$Severity) {
    $AuditResults.MitoteReport += [pscustomobject]@{
        Code        = $Code
        Category    = $Category
        Description = $Description
        Severity    = $Severity
        Timestamp   = (Get-Date).ToString("o")
    }
    if ($Severity -eq "CRITICAL" -or $Severity -eq "HIGH") {
        $AuditResults.Metrics.CriticalErrors++
        $AuditResults.Status = "FAIL"
    } elseif ($Severity -eq "MEDIUM" -or $Severity -eq "LOW") {
        $AuditResults.Metrics.Warnings++
    }
}

# ==============================================================================
# FASE 0: AUDITORÍA DE ENTORNO Y MITOTERO DE GIT
# ==============================================================================
Log-Step "FASE 0" "Escaneando estado del entorno y rastros de Git..." "INFO"

if (Get-Command "git" -ErrorAction SilentlyContinue) {
    Push-Location $RootPath
    $GitStatus = git status --porcelain 2>&1
    if ($GitStatus) {
        Log-Step "GIT" "Se encontraron archivos no guardados o basura en el Git Working Tree!" "WARN"
        Add-Mitote "MITOTE-GIT-001" "GIT_DIRTY" "El repositorio tiene cambios sin commit o archivos untracked." "MEDIUM"
    } else {
        Log-Step "GIT" "Working tree de Git completamente limpio." "OK"
        $AuditResults.Metrics.PassedChecks++
    }
    Pop-Location
} else {
    Log-Step "GIT" "Git no está en el PATH. No podemos fiscalizar commits." "WARN"
}

# ==============================================================================
# FASE 1: DESENTIERRO DE ARCHIVOS FANTASMA, DUPES Y BASURA (EL MITOTE)
# ==============================================================================
Log-Step "FASE 1" "Iniciando escaneo de basura, archivos duplicados y carpetas anidadas..." "INFO"

# 1.1 Verificar subcarpeta anidada ratera
if (Test-Path (Join-Path $RootPath "rx-dispatch")) {
    Log-Step "ESTRUCTURA" "DETECTADO: Carpeta duplicada ./rx-dispatch/ entorpeciendo la raíz!" "FAIL"
    Add-Mitote "RXD-ERR-002" "GHOST_FOLDER" "Directorio anidado ./rx-dispatch/ detectado. Violación estructural de MAP.md." "CRITICAL"
}

# 1.2 Buscar archivos prohibidos/basura
foreach ($Pattern in $BannedFilesOrFolders) {
    $Junk = Get-ChildItem -Path $RootPath -Recurse -Filter $Pattern -ErrorAction SilentlyContinue
    foreach ($Item in $Junk) {
        Log-Step "BASURA" "Archivo/Carpeta fuera de norma hallado: $($Item.FullName)" "CHISME"
        Add-Mitote "MITOTE-JUNK-001" "UNAUTHORIZED_FILE" "Archivo o basura detectada: $($Item.Name)" "LOW"
    }
}

# 1.3 Verificar los 6 Markdown Sagrados
foreach ($Doc in $NormativeDocs) {
    $DocPath = Join-Path -Path $RootPath -ChildPath $Doc
    if (Test-Path $DocPath) {
        Log-Step "NORMATIVA" "Documento baseline presente: $Doc" "OK"
        $AuditResults.Metrics.PassedChecks++
    } else {
        Log-Step "NORMATIVA" "Falta el documento sagrado: $Doc" "FAIL"
        Add-Mitote "RXD-ERR-001" "MISSING_DOC" "Documento normativo ausente en raíz: $Doc" "CRITICAL"
    }
}

# 1.4 Verificar los 9 Microservicios (Entrypoints)
foreach ($Svc in $AuthoritativeServices) {
    $SvcPath = Join-Path -Path $RootPath -ChildPath (Join-Path $Svc "main.go")
    if (Test-Path $SvcPath) {
        Log-Step "ARQUITECTURA" "Entrypoint confirmado: $Svc/main.go" "OK"
        $AuditResults.Metrics.PassedChecks++
    } else {
        Log-Step "ARQUITECTURA" "Microservicio extraviado o no creado: $Svc/main.go" "FAIL"
        Add-Mitote "RXD-ERR-003" "MISSING_SERVICE" "Punto de entrada de dominio ausente: $Svc/main.go" "HIGH"
    }
}

# ==============================================================================
# FASE 2: INSPECCIÓN FORENSE DE CODIFICACIÓN (BOM, SALTOS DE LÍNEA, ARCHIVOS VACÍOS)
# ==============================================================================
Log-Step "FASE 2" "Revisando codificación byte por byte (BOM, CRLF/LF, vacíos)..." "INFO"

$AllFiles = Get-ChildItem -Path $RootPath -Recurse -File -Exclude .git, bin, assets
$AuditResults.Metrics.TotalFilesScanned = $AllFiles.Count

foreach ($File in $AllFiles) {
    # Archivos de 0 Bytes
    if ($File.Length -eq 0) {
        Log-Step "FORENSE" "Archivo fantasma de 0 Bytes detectado: $($File.Name)" "CHISME"
        Add-Mitote "MITOTE-EMPTY-001" "ZERO_BYTE_FILE" "El archivo $($File.FullName) está totalmente vacío." "MEDIUM"
        continue
    }

    # Detección de BOM
    $Bytes = Get-Content -Path $File.FullName -Encoding Byte -TotalCount 3 -ErrorAction SilentlyContinue
    if ($Bytes.Count -ge 3 -and $Bytes[0] -eq 0xEF -and $Bytes[1] -eq 0xBB -and $Bytes[2] -eq 0xBF) {
        Log-Step "ENCODING" "BOM (Byte Order Mark) detectado en: $($File.Name)" "WARN"
        Add-Mitote "RXD-WARN-001" "BOM_DETECTED" "El archivo $($File.Name) tiene BOM UTF-8. Se requiere UTF-8 plano." "MEDIUM"
    }

    # Conteo de líneas de código si es .go
    if ($File.Extension -eq ".go") {
        $Lines = (Get-Content -Path $File.FullName).Count
        $AuditResults.Metrics.TotalLinesOfCode += $Lines
    }
}

# ==============================================================================
# FASE 3: INSPECCIÓN DE CÓDIGO GO (PUERTOS, SHA-256, DICOM Y CONTRATOS)
# ==============================================================================
Log-Step "FASE 3" "Escrutando el código Go por contratos duros (127.0.0.1, SHA-256, DICOM)..." "INFO"

$GoFiles = Get-ChildItem -Path $RootPath -Recurse -Filter "*.go" -Exclude bin

if ($GoFiles) {
    $CombinedCode = Get-Content -Path $GoFiles.FullName -Raw -ErrorAction SilentlyContinue

    # Auditoría de Puertos / Binding (Loopback check)
    if ($CombinedCode -match '0\.0\.0\.0:808[1-8]') {
        Log-Step "SEGURIDAD" "Peligro: Microservicio interno escuchando en 0.0.0.0 en vez de 127.0.0.1!" "FAIL"
        Add-Mitote "RXD-SEC-001" "INSECURE_BINDING" "Se detectó binding a 0.0.0.0 en puertos internos 8081-8088." "HIGH"
    } else {
        Log-Step "SEGURIDAD" "Binding de puertos internos seguro/verificado." "OK"
        $AuditResults.Metrics.PassedChecks++
    }

    # Uso Criptográfico SHA-256
    if ($CombinedCode -match 'crypto/sha256' -or $CombinedCode -match 'sha256\.New') {
        Log-Step "CRIPTO" "Implementación de hash SHA-256 localizada en código." "OK"
        $AuditResults.Metrics.PassedChecks++
    } else {
        Log-Step "CRIPTO" "No se detecta import de crypto/sha256 en ningún archivo Go." "WARN"
        Add-Mitote "RXD-WARN-005" "NO_SHA256" "Falta evidencia de implementación de SHA-256 en código fuente." "MEDIUM"
    }

    # Verificación del comando A-ASSOCIATE-RJ de DICOM
    if ($CombinedCode -match 'A-ASSOCIATE-RJ' -or $CombinedCode -match 'RejectAssociation') {
        Log-Step "DICOM" "Control protocolario A-ASSOCIATE-RJ localizado." "OK"
        $AuditResults.Metrics.PassedChecks++
    } else {
        Log-Step "DICOM" "No se encontró referencia explícita al rechazo A-ASSOCIATE-RJ." "WARN"
        Add-Mitote "RXD-WARN-006" "NO_DICOM_REJECT" "Código no demuestra el manejo de rechazo protocolario A-ASSOCIATE-RJ." "MEDIUM"
    }
} else {
    Log-Step "CÓDIGO" "No hay archivos .go para auditar en el proyecto actual!" "FAIL"
    Add-Mitote "RXD-ERR-005" "NO_GO_FILES" "Código fuente Go inexistente en el árbol." "CRITICAL"
}

# ==============================================================================
# FASE 4: AUDITORÍA DE CAMPOS DE AUDIT TRAIL / ePHI
# ==============================================================================
Log-Step "FASE 4" "Verificando los 14 campos obligatorios de la bitácora ePHI..." "INFO"

$ModelsPath = Join-Path -Path $RootPath -ChildPath "internal/models"
$ModelFiles = Get-ChildItem -Path $ModelsPath -Filter "*.go" -Recurse -ErrorAction SilentlyContinue

if ($ModelFiles) {
    $ModelCode = Get-Content -Path $ModelFiles.FullName -Raw
    foreach ($Field in $MandatoryAuditTrailFields) {
        if ($ModelCode -match $Field) {
            Log-Step "EPHI_FIELD" "Campo ePHI de auditoría verificado: $Field" "OK"
            $AuditResults.Metrics.PassedChecks++
        } else {
            Log-Step "EPHI_FIELD" "Campo ePHI ausente en modelos: $Field" "WARN"
            Add-Mitote "RXD-EPHI-001" "MISSING_EPHI_FIELD" "El modelo Go no declara el campo obligatorio: $Field" "MEDIUM"
        }
    }
} else {
    Log-Step "EPHI_FIELD" "Directorio internal/models/ ausente o vacío." "WARN"
    Add-Mitote "RXD-EPHI-002" "NO_MODELS" "Sin acceso a los modelos de datos en internal/models/." "HIGH"
}

# ==============================================================================
# FASE 5: COMPILACIÓN Y TOOLCHAIN DE GO
# ==============================================================================
Log-Step "FASE 5" "Corriendo la navaja de afeitar de Go (go vet y go.mod)..." "INFO"

if (Get-Command "go" -ErrorAction SilentlyContinue) {
    if (Test-Path (Join-Path $RootPath "go.mod")) {
        Push-Location $RootPath
        $VetResult = go vet ./... 2>&1
        if ($LASTEXITCODE -eq 0) {
            Log-Step "COMPILADOR" "go vet pasó sin un solo chisme de sintaxis!" "OK"
            $AuditResults.Metrics.PassedChecks++
        } else {
            Log-Step "COMPILADOR" "go vet encontró broncas de sintaxis!" "FAIL"
            Add-Mitote "RXD-ERR-006" "GO_VET_FAIL" "Errores de sintaxis detectados: $VetResult" "HIGH"
        }
        Pop-Location
    } else {
        Log-Step "COMPILADOR" "Falta el go.mod en la raíz del proyecto." "FAIL"
        Add-Mitote "RXD-ERR-007" "NO_GOMOD" "Sin go.mod en la raíz." "CRITICAL"
    }
} else {
    Log-Step "COMPILADOR" "Go no está instalado en este sistema." "WARN"
}

# ==============================================================================
# FASE 6: RESUMEN EJECUTIVO Y EMISIÓN DE REPORTE JSON
# ==============================================================================
Write-Host ""
Write-Host "==========================================================================" -ForegroundColor Cyan
Write-Host "                   RESUMEN PERICIAL ULTRA-MITOTERO" -ForegroundColor Cyan
Write-Host "==========================================================================" -ForegroundColor Cyan
Write-Host " Archivos Escaneados    : $($AuditResults.Metrics.TotalFilesScanned)" -ForegroundColor White
Write-Host " Líneas de Código (LOC) : $($AuditResults.Metrics.TotalLinesOfCode)" -ForegroundColor White
Write-Host " Pruebas Aprobadas      : $($AuditResults.Metrics.PassedChecks)" -ForegroundColor Green
Write-Host " Advertencias / Chismes : $($AuditResults.Metrics.Warnings)" -ForegroundColor Yellow
Write-Host " Errores Críticos       : $($AuditResults.Metrics.CriticalErrors)" -ForegroundColor Red
Write-Host " Estado Global          : $($AuditResults.Status)" -ForegroundColor $(if ($AuditResults.Status -eq "PASS") { "Green" } else { "Red" })
Write-Host "==========================================================================" -ForegroundColor Cyan

if ($ReportJson) {
    $ReportFile = Join-Path $RootPath "audit_report_ultra.json"
    $AuditResults | ConvertTo-Json -Depth 6 | Out-File -FilePath $ReportFile -Encoding utf8
    Log-Step "REPORTE" "Expediente JSON ultra-detallado generado en: $ReportFile" "OK"
}

Write-Host ""
Write-Host "[ PRESIONE CUALQUIER TECLA PARA CONTINUAR Y CONSERVAR LA CONSOLA ABIERTA ]" -ForegroundColor Yellow
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")