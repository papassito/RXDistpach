<#
================================================================================
 RX DISPATCH :: MOTOR DE AUDITORÍA PERICIAL Y ESCANEO PROFUNDO v3.1
 "Ultra-Mitotero, pero ahora declara únicamente lo que puede demostrar"
================================================================================

 PRINCIPIOS:
 - NOT_DETECTED != DOES_NOT_EXIST
 - STRING_FOUND != IMPLEMENTED
 - IMPORT_FOUND != FUNCTIONAL_USAGE
 - GO_VET_FAIL != SYNTAX_ERROR
 - DOCUMENT_FOUND != AUTHORITATIVE
 - DECLARED != EXECUTED

 SOLO LECTURA:
 Este script NO modifica archivos.
================================================================================
#>

param(
    [string]$Root = ".",
    [bool]$RequireBaselineAtRoot = $true
)

$ErrorActionPreference = "Continue"
$RepoRoot = (Resolve-Path $Root).Path
Set-Location $RepoRoot

# ==============================================================================
# CONFIGURACIÓN
# ==============================================================================

$BaselineDocs = @(
    "README.md",
    "ARCHITECTURE.md",
    "REQUIREMENTS.md",
    "CONTRACTS.md",
    "MAP.md",
    "SECURITY.md"
)

$ExpectedEntrypoints = @(
    "cmd/audit/main.go",
    "cmd/delivery/main.go",
    "cmd/gateway/main.go",
    "cmd/image/main.go",
    "cmd/reader/main.go",
    "cmd/result/main.go",
    "cmd/security/main.go",
    "cmd/storage/main.go",
    "cmd/study/main.go"
)

$EPHIFields = @(
    "event_timestamp",
    "event_action",
    "event_outcome",
    "event_error_code",
    "user_id",
    "source_ip",
    "source_ae_title",
    "destination_ip",
    "destination_ae_title",
    "patient_id",
    "study_instance_uid",
    "accession_number",
    "number_of_instances",
    "security_tls_status"
)

$Findings = @()
$PassedChecks = 0

# ==============================================================================
# FUNCIONES
# ==============================================================================

function Add-Finding {
    param(
        [string]$Code,
        [string]$Category,
        [string]$Description,
        [ValidateSet("INFO","LOW","MEDIUM","HIGH","CRITICAL")]
        [string]$Severity,
        [string]$EvidenceScope = "",
        [string]$Evidence = ""
    )

    $script:Findings += [PSCustomObject]@{
        Code          = $Code
        Category      = $Category
        Description   = $Description
        Severity      = $Severity
        EvidenceScope = $EvidenceScope
        Evidence      = $Evidence
        Timestamp     = (Get-Date).ToString("o")
    }
}

function Add-Pass {
    param(
        [string]$Message
    )

    $script:PassedChecks++
    Write-Host "[OK] $Message" -ForegroundColor Green
}

function Show-Warn {
    param([string]$Message)
    Write-Host "[WARN] $Message" -ForegroundColor Yellow
}

function Show-Fail {
    param([string]$Message)
    Write-Host "[FAIL] $Message" -ForegroundColor Red
}

function Show-Info {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor Cyan
}

function Invoke-GoCheck {
    param(
        [string]$Name,
        [string[]]$Arguments,
        [string]$Code,
        [string]$Category,
        [string]$SuccessText
    )

    $output = @()

    try {
        $output = & go @Arguments 2>&1
        $exitCode = $LASTEXITCODE
    }
    catch {
        $output = @($_.Exception.Message)
        $exitCode = 1
    }

    if ($exitCode -eq 0) {
        Add-Pass $SuccessText
        return $true
    }

    $details = ($output | Out-String).Trim()

    Add-Finding `
        -Code $Code `
        -Category $Category `
        -Description "$Name no completó correctamente. Revisar la salida real de la herramienta; este hallazgo no se clasifica automáticamente como error de sintaxis." `
        -Severity "HIGH" `
        -EvidenceScope "go workspace" `
        -Evidence $details

    Show-Fail "$Name terminó con código $exitCode."

    if ($details) {
        Write-Host $details -ForegroundColor DarkYellow
    }

    return $false
}

# ==============================================================================
# CABECERA
# ==============================================================================

Clear-Host

Write-Host "==========================================================================" -ForegroundColor Cyan
Write-Host " 🔬 RX DISPATCH :: MOTOR DE AUDITORÍA PERICIAL Y ESCANEO PROFUNDO v3.1" -ForegroundColor Cyan
Write-Host "    'Ultra-Mitotero, pero sin inventar delitos'" -ForegroundColor DarkCyan
Write-Host "==========================================================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Repositorio: $RepoRoot" -ForegroundColor Gray

# ==============================================================================
# RECOLECCIÓN
# ==============================================================================

$AllFiles = Get-ChildItem -Path $RepoRoot -Recurse -File -ErrorAction SilentlyContinue |
    Where-Object {
        $_.FullName -notmatch '\\\.git\\' -and
        $_.FullName -notmatch '\\bin\\' -and
        $_.FullName -notmatch '\\obj\\' -and
        $_.FullName -notmatch '\\vendor\\'
    }

$GoFiles = @(
    $AllFiles | Where-Object { $_.Extension -eq ".go" }
)

$ModelFiles = @(
    $GoFiles | Where-Object {
        $_.FullName -match '\\internal\\models\\'
    }
)

if ($ModelFiles.Count -eq 0) {
    # Fallback: si no hay internal/models, inspeccionar internal completo.
    $ModelFiles = @(
        $GoFiles | Where-Object {
            $_.FullName -match '\\internal\\'
        }
    )
}

$TotalLOC = 0

foreach ($file in $GoFiles) {
    try {
        $TotalLOC += (Get-Content $file.FullName -ErrorAction Stop).Count
    }
    catch {
        # No convertir automáticamente un fallo de lectura en LOC falso.
    }
}

# ==============================================================================
# FASE 0 — ENTORNO Y GIT
# ==============================================================================

Write-Host ""
Show-Info "[FASE 0] Inspeccionando entorno y estado Git..."

$gitAvailable = Get-Command git -ErrorAction SilentlyContinue

if ($gitAvailable) {

    $insideGit = (& git rev-parse --is-inside-work-tree 2>$null)

    if ($LASTEXITCODE -eq 0 -and $insideGit -eq "true") {

        $gitStatus = @(& git status --porcelain 2>$null)

        if ($gitStatus.Count -gt 0) {

            Add-Finding `
                -Code "MITOTE-GIT-001" `
                -Category "GIT_DIRTY" `
                -Description "El working tree contiene cambios sin commit y/o archivos no rastreados. Esto se reporta como estado de trabajo, no como corrupción ni como 'basura'." `
                -Severity "MEDIUM" `
                -EvidenceScope "git working tree" `
                -Evidence (($gitStatus | Out-String).Trim())

            Show-Warn "[GIT] Working tree con cambios pendientes."

        }
        else {
            Add-Pass "[GIT] Working tree limpio."
        }

    }
    else {

        Add-Finding `
            -Code "MITOTE-GIT-002" `
            -Category "GIT_NOT_DETECTED" `
            -Description "La ruta inspeccionada no fue identificada como working tree Git." `
            -Severity "INFO" `
            -EvidenceScope $RepoRoot

        Show-Info "[GIT] No se identificó un repositorio Git en esta ruta."
    }

}
else {

    Add-Finding `
        -Code "MITOTE-GIT-003" `
        -Category "GIT_COMMAND_UNAVAILABLE" `
        -Description "El comando git no está disponible en el entorno actual. No se evaluó el estado del working tree." `
        -Severity "INFO" `
        -EvidenceScope "local environment"

    Show-Info "[GIT] Comando git no disponible."
}

# ==============================================================================
# FASE 1 — BASELINE DOCUMENTAL
# ==============================================================================

Write-Host ""
Show-Info "[FASE 1] Verificando documentos baseline en rutas conocidas..."

foreach ($doc in $BaselineDocs) {

    $rootPath = Join-Path $RepoRoot $doc
    $docsPath = Join-Path $RepoRoot ("docs\" + $doc)

    $inRoot = Test-Path $rootPath
    $inDocs = Test-Path $docsPath

    if ($inRoot) {

        Add-Pass "[NORMATIVA] Encontrado en raíz esperada: $doc"
        continue
    }

    if ($inDocs) {

        $severity = if ($RequireBaselineAtRoot) { "CRITICAL" } else { "MEDIUM" }

        Add-Finding `
            -Code "RXD-DOC-ROOT-001" `
            -Category "DOC_NOT_IN_EXPECTED_ROOT" `
            -Description "$doc no se encontró en la raíz esperada, pero existe una copia en docs\. El hallazgo se limita a ubicación; no se afirma inexistencia global." `
            -Severity $severity `
            -EvidenceScope "repository root" `
            -Evidence $docsPath

        if ($RequireBaselineAtRoot) {
            Show-Fail "[NORMATIVA] $doc no está en raíz; existe en docs\."
        }
        else {
            Show-Warn "[NORMATIVA] $doc no está en raíz; existe en docs\."
        }

        continue
    }

    Add-Finding `
        -Code "RXD-DOC-001" `
        -Category "DOC_NOT_DETECTED" `
        -Description "No se encontró $doc ni en la raíz esperada ni en docs\ dentro del alcance inspeccionado." `
        -Severity "CRITICAL" `
        -EvidenceScope "repository root + docs"

    Show-Fail "[NORMATIVA] $doc no fue detectado en raíz ni en docs\."
}

# ==============================================================================
# FASE 1B — ENTRYPOINTS
# ==============================================================================

Write-Host ""
Show-Info "[FASE 1B] Verificando entrypoints esperados..."

foreach ($relative in $ExpectedEntrypoints) {

    $path = Join-Path $RepoRoot $relative

    if (Test-Path $path) {
        Add-Pass "[ARQUITECTURA] Entrypoint localizado: $relative"
    }
    else {

        Add-Finding `
            -Code "RXD-ENTRY-001" `
            -Category "ENTRYPOINT_NOT_DETECTED" `
            -Description "No se encontró el entrypoint esperado $relative en la ruta exacta inspeccionada." `
            -Severity "HIGH" `
            -EvidenceScope $relative

        Show-Fail "[ARQUITECTURA] No detectado: $relative"
    }
}

# ==============================================================================
# FASE 2 — CODIFICACIÓN / VACÍOS
# ==============================================================================

Write-Host ""
Show-Info "[FASE 2] Revisando archivos vacíos y BOM UTF-8..."

foreach ($file in $AllFiles) {

    if ($file.Length -eq 0) {

        Add-Finding `
            -Code "MITOTE-FILE-EMPTY-001" `
            -Category "EMPTY_FILE" `
            -Description "Se encontró un archivo de longitud cero. Puede ser intencional; revisar su función antes de eliminarlo." `
            -Severity "LOW" `
            -EvidenceScope "filesystem" `
            -Evidence $file.FullName
    }

    try {

        $bytes = [System.IO.File]::ReadAllBytes($file.FullName)

        if (
            $bytes.Length -ge 3 -and
            $bytes[0] -eq 0xEF -and
            $bytes[1] -eq 0xBB -and
            $bytes[2] -eq 0xBF
        ) {

            Add-Finding `
                -Code "MITOTE-BOM-001" `
                -Category "UTF8_BOM_DETECTED" `
                -Description "Archivo UTF-8 con BOM detectado. Se reporta como característica de codificación, no como corrupción." `
                -Severity "LOW" `
                -EvidenceScope "filesystem" `
                -Evidence $file.FullName
        }

    }
    catch {
        # No bloquear la auditoría por un archivo que no pudo leerse byte a byte.
    }
}

# ==============================================================================
# FASE 3 — BINDINGS
# ==============================================================================

Write-Host ""
Show-Info "[FASE 3] Inspeccionando bindings y patrones de red..."

$publicBindings = @()

foreach ($file in $GoFiles) {

    $matches = Select-String `
        -Path $file.FullName `
        -Pattern '0\.0\.0\.0' `
        -ErrorAction SilentlyContinue

    if ($matches) {
        $publicBindings += $matches
    }
}

if ($publicBindings.Count -eq 0) {

    Add-Pass "[RED] No se detectó literal 0.0.0.0 en los archivos Go inspeccionados."

}
else {

    Add-Finding `
        -Code "RXD-NET-001" `
        -Category "PUBLIC_BINDING_DETECTED" `
        -Description "Se detectó el literal 0.0.0.0 en código Go. Esto indica una posible exposición amplia; la accesibilidad efectiva depende del runtime, firewall y despliegue." `
        -Severity "HIGH" `
        -EvidenceScope "Go source" `
        -Evidence (($publicBindings | ForEach-Object {
            "$($_.Path):$($_.LineNumber): $($_.Line.Trim())"
        }) -join "`n")

    Show-Warn "[RED] Se detectaron referencias a 0.0.0.0."
}

# ==============================================================================
# FASE 3B — SHA-256
# ==============================================================================

Write-Host ""
Show-Info "[CRIPTO] Buscando evidencia directa de uso de SHA-256..."

$shaImports = @()
$shaUsage   = @()

foreach ($file in $GoFiles) {

    $shaImports += Select-String `
        -Path $file.FullName `
        -Pattern '"crypto/sha256"' `
        -ErrorAction SilentlyContinue

    $shaUsage += Select-String `
        -Path $file.FullName `
        -Pattern 'sha256\.(Sum256|New)\s*\(' `
        -ErrorAction SilentlyContinue
}

if ($shaUsage.Count -gt 0) {

    Add-Pass "[CRIPTO] Se detectó uso directo de sha256.Sum256() y/o sha256.New()."

}
elseif ($shaImports.Count -gt 0) {

    Add-Finding `
        -Code "RXD-WARN-005A" `
        -Category "SHA256_IMPORT_WITHOUT_DIRECT_USAGE" `
        -Description "Se encontró import de crypto/sha256, pero el scanner no encontró llamadas directas a sha256.Sum256() o sha256.New(). El import por sí solo no acredita implementación criptográfica funcional." `
        -Severity "MEDIUM" `
        -EvidenceScope "Go source" `
        -Evidence (($shaImports | ForEach-Object {
            "$($_.Path):$($_.LineNumber)"
        }) -join "`n")

    Show-Warn "[CRIPTO] Import SHA-256 detectado, pero sin uso directo verificable."

}
else {

    Add-Finding `
        -Code "RXD-WARN-005" `
        -Category "SHA256_USAGE_NOT_DETECTED" `
        -Description "El scanner no encontró import ni uso directo verificable de SHA-256 en los archivos Go inspeccionados." `
        -Severity "MEDIUM" `
        -EvidenceScope "Go source"

    Show-Warn "[CRIPTO] No se detectó uso directo de SHA-256."
}

# ==============================================================================
# FASE 3C — DICOM A-ASSOCIATE-RJ
# ==============================================================================

Write-Host ""
Show-Info "[DICOM] Buscando evidencia de A-ASSOCIATE-RJ..."

$dicomLiteral = @()
$dicomIdentifier = @()

foreach ($file in $GoFiles) {

    $dicomLiteral += Select-String `
        -Path $file.FullName `
        -Pattern 'A-ASSOCIATE-RJ' `
        -ErrorAction SilentlyContinue

    $dicomIdentifier += Select-String `
        -Path $file.FullName `
        -Pattern '\bDicomAssociateReject\b' `
        -ErrorAction SilentlyContinue
}

if ($dicomLiteral.Count -eq 0) {

    Add-Finding `
        -Code "RXD-WARN-006" `
        -Category "DICOM_REJECT_NOT_DETECTED" `
        -Description "El scanner no encontró referencia explícita a A-ASSOCIATE-RJ en los archivos Go inspeccionados. Esto no demuestra por sí solo que el comportamiento protocolario sea inexistente." `
        -Severity "MEDIUM" `
        -EvidenceScope "Go source"

    Show-Warn "[DICOM] No se detectó referencia explícita a A-ASSOCIATE-RJ."

}
elseif ($dicomIdentifier.Count -le 1) {

    Add-Finding `
        -Code "RXD-WARN-006A" `
        -Category "DICOM_REJECT_DECLARED_ONLY" `
        -Description "Se detectó A-ASSOCIATE-RJ y/o su constante, pero no se encontró evidencia suficiente de consumo desde otro punto del código. Declarar la constante no acredita ejecución del rechazo protocolario." `
        -Severity "MEDIUM" `
        -EvidenceScope "Go source" `
        -Evidence (($dicomLiteral | ForEach-Object {
            "$($_.Path):$($_.LineNumber): $($_.Line.Trim())"
        }) -join "`n")

    Show-Warn "[DICOM] Referencia declarada, uso funcional no acreditado por este scanner."

}
else {

    Add-Pass "[DICOM] Se detectó declaración y referencia adicional de DicomAssociateReject."

    Add-Finding `
        -Code "RXD-INFO-DICOM-001" `
        -Category "DICOM_REJECT_REFERENCE_DETECTED" `
        -Description "Se detectaron referencias de A-ASSOCIATE-RJ fuera de una única declaración. Esto demuestra referencias de código, no certifica por sí solo conformidad DICOM." `
        -Severity "INFO" `
        -EvidenceScope "Go source" `
        -Evidence (($dicomIdentifier | ForEach-Object {
            "$($_.Path):$($_.LineNumber): $($_.Line.Trim())"
        }) -join "`n")
}

# ==============================================================================
# FASE 4 — 14 CAMPOS ePHI
# ==============================================================================

Write-Host ""
Show-Info "[FASE 4] Verificando detección de los 14 campos ePHI..."

if ($ModelFiles.Count -eq 0) {

    Add-Finding `
        -Code "RXD-EPHI-000" `
        -Category "MODEL_SCOPE_NOT_DETECTED" `
        -Description "No se localizaron archivos Go dentro del alcance esperado para modelos. No se evaluó de forma concluyente la presencia de campos ePHI." `
        -Severity "HIGH" `
        -EvidenceScope "internal/models or internal"

    Show-Fail "[EPHI] No se encontró alcance de modelos evaluable."

}
else {

    foreach ($field in $EPHIFields) {

        $pattern = 'json:"' + [regex]::Escape($field) + '"'
        $matches = @()

        foreach ($file in $ModelFiles) {

            $matches += Select-String `
                -Path $file.FullName `
                -Pattern $pattern `
                -ErrorAction SilentlyContinue
        }

        if ($matches.Count -gt 0) {

            Add-Pass "[EPHI_FIELD] Detectado tag JSON: $field"

        }
        else {

            Add-Finding `
                -Code "RXD-EPHI-001" `
                -Category "EPHI_FIELD_NOT_DETECTED" `
                -Description "El scanner no encontró el campo esperado '$field' como tag JSON dentro del alcance de modelos inspeccionado." `
                -Severity "MEDIUM" `
                -EvidenceScope (($ModelFiles.FullName) -join "; ")

            Show-Warn "[EPHI_FIELD] No detectado en alcance inspeccionado: $field"
        }
    }
}

# ==============================================================================
# FASE 5 — GO BUILD
# ==============================================================================

Write-Host ""
Show-Info "[FASE 5] Ejecutando go build ./..."

$buildOK = Invoke-GoCheck `
    -Name "go build ./..." `
    -Arguments @("build", "./...") `
    -Code "RXD-ERR-005" `
    -Category "GO_BUILD_FAIL" `
    -SuccessText "[COMPILADOR] go build ./... completó correctamente."

# ==============================================================================
# FASE 6 — GO TEST
# ==============================================================================

Write-Host ""
Show-Info "[FASE 6] Ejecutando go test ./... -count=1..."

$testOK = Invoke-GoCheck `
    -Name "go test ./... -count=1" `
    -Arguments @("test", "./...", "-count=1") `
    -Code "RXD-ERR-007" `
    -Category "GO_TEST_FAIL" `
    -SuccessText "[PRUEBAS] go test ./... -count=1 completó correctamente."

# ==============================================================================
# FASE 7 — GO VET
# ==============================================================================

Write-Host ""
Show-Info "[FASE 7] Ejecutando go vet ./..."

$vetOK = Invoke-GoCheck `
    -Name "go vet ./..." `
    -Arguments @("vet", "./...") `
    -Code "RXD-ERR-006" `
    -Category "GO_VET_FAIL" `
    -SuccessText "[ANÁLISIS] go vet ./... completó correctamente."

# ==============================================================================
# MÉTRICAS Y ESTADO
# ==============================================================================

$Warnings = @(
    $Findings | Where-Object {
        $_.Severity -in @("LOW", "MEDIUM")
    }
).Count

$HighErrors = @(
    $Findings | Where-Object {
        $_.Severity -eq "HIGH"
    }
).Count

$CriticalErrors = @(
    $Findings | Where-Object {
        $_.Severity -eq "CRITICAL"
    }
).Count

$Informational = @(
    $Findings | Where-Object {
        $_.Severity -eq "INFO"
    }
).Count

$BlockingFindings = $HighErrors + $CriticalErrors

$Status = if ($BlockingFindings -gt 0) {
    "FAIL"
}
elseif ($Warnings -gt 0) {
    "PASS_WITH_WARNINGS"
}
else {
    "PASS"
}

# ==============================================================================
# REPORTE JSON
# ==============================================================================

$Report = [PSCustomObject]@{
    Timestamp = (Get-Date).ToString("o")

    Status = $Status

    Scope = [PSCustomObject]@{
        RepositoryRoot        = $RepoRoot
        RequireBaselineAtRoot = $RequireBaselineAtRoot
        Statement             = "Los hallazgos se limitan al alcance inspeccionado. NOT_DETECTED no implica inexistencia global."
    }

    Metrics = [PSCustomObject]@{
        TotalFilesScanned = $AllFiles.Count
        GoFilesScanned    = $GoFiles.Count
        TotalLinesOfCode  = $TotalLOC

        PassedChecks      = $PassedChecks
        Informational     = $Informational
        Warnings          = $Warnings
        HighErrors        = $HighErrors
        CriticalErrors    = $CriticalErrors
        BlockingFindings  = $BlockingFindings
    }

    ExecutionGate = [PSCustomObject]@{
        GoBuild = $buildOK
        GoTest  = $testOK
        GoVet   = $vetOK
    }

    Findings = $Findings
}

$ReportPath = Join-Path $RepoRoot "audit_report_ultra.json"

$Report |
    ConvertTo-Json -Depth 10 |
    Set-Content -Path $ReportPath -Encoding UTF8

# ==============================================================================
# RESUMEN
# ==============================================================================

Write-Host ""
Write-Host "==========================================================================" -ForegroundColor Cyan
Write-Host "                 RESUMEN PERICIAL ULTRA-MITOTERO v3.1" -ForegroundColor Cyan
Write-Host "=========================================================================="

Write-Host (" Archivos inspeccionados        : {0}" -f $AllFiles.Count)
Write-Host (" Archivos Go                    : {0}" -f $GoFiles.Count)
Write-Host (" Líneas Go aproximadas (LOC)    : {0}" -f $TotalLOC)
Write-Host (" Comprobaciones aprobadas       : {0}" -f $PassedChecks)
Write-Host (" Informativos                   : {0}" -f $Informational)
Write-Host (" Advertencias LOW/MEDIUM        : {0}" -f $Warnings)
Write-Host (" Errores HIGH                   : {0}" -f $HighErrors)
Write-Host (" Errores CRITICAL               : {0}" -f $CriticalErrors)
Write-Host (" Bloqueos confirmados           : {0}" -f $BlockingFindings)
Write-Host ""

Write-Host (" go build ./...                 : {0}" -f $(if ($buildOK) { "PASS" } else { "FAIL" }))
Write-Host (" go test ./... -count=1         : {0}" -f $(if ($testOK) { "PASS" } else { "FAIL" }))
Write-Host (" go vet ./...                   : {0}" -f $(if ($vetOK) { "PASS" } else { "FAIL" }))

Write-Host ""

switch ($Status) {

    "PASS" {
        Write-Host " Estado Global                  : PASS" -ForegroundColor Green
    }

    "PASS_WITH_WARNINGS" {
        Write-Host " Estado Global                  : PASS_WITH_WARNINGS" -ForegroundColor Yellow
    }

    default {
        Write-Host " Estado Global                  : FAIL" -ForegroundColor Red
    }
}

Write-Host "=========================================================================="

Write-Host ""
Write-Host "[REGLA PERICIAL]" -ForegroundColor Cyan
Write-Host "  NOT_DETECTED != DOES_NOT_EXIST"
Write-Host "  STRING_FOUND  != IMPLEMENTED"
Write-Host "  IMPORT_FOUND  != FUNCTIONAL_USAGE"
Write-Host "  GO_VET_FAIL   != SYNTAX_ERROR"
Write-Host ""

Write-Host "[REPORTE] $ReportPath" -ForegroundColor Green

Write-Host ""
Read-Host "[ PRESIONE ENTER PARA CONSERVAR LA CONSOLA ABIERTA ]" | Out-Null

if ($BlockingFindings -gt 0) {
    exit 1
}

exit 0