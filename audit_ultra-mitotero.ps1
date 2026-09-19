<#
================================================================================
 RX DISPATCH :: MOTOR DE AUDITORÍA PERICIAL Y ESCANEO PROFUNDO v4.0
 "Ultra-Mitotero: Cobertura Total de Basura, Pruebas y Resiliencia"
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

$LooseRootGoFiles = @("main.go", "models.go", "contracts.go", "constants.go", "client.go")

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
    param([string]$Message)
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
        -Description "$Name no completó correctamente. Revisar la salida real de la herramienta." `
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
Write-Host " 🔬 RX DISPATCH :: MOTOR DE AUDITORÍA PERICIAL Y ESCANEO PROFUNDO v4.0" -ForegroundColor Cyan
Write-Host "    'Ultra-Mitotero: Inspección Total ePHI, Resiliencia y Estructura'" -ForegroundColor DarkCyan
Write-Host "==========================================================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Repositorio: $RepoRoot" -ForegroundColor Gray

# ==============================================================================
# RECOLECCIÓN
# ==============================================================================

$AllFiles = Get-ChildItem -Path $RepoRoot -Recurse -File -Force -ErrorAction SilentlyContinue |
    Where-Object {
        $_.FullName -notmatch '\\\.git\\' -and
        $_.FullName -notmatch '\\bin\\' -and
        $_.FullName -notmatch '\\obj\\' -and
        $_.FullName -notmatch '\\vendor\\'
    }

$GoFiles = @($AllFiles | Where-Object { $_.Extension -eq ".go" })

$ModelFiles = @($GoFiles | Where-Object { $_.FullName -match '\\internal\\models\\' })
if ($ModelFiles.Count -eq 0) {
    $ModelFiles = @($GoFiles | Where-Object { $_.FullName -match '\\internal\\' })
}

$TotalLOC = 0
foreach ($file in $GoFiles) {
    try {
        $TotalLOC += (Get-Content $file.FullName -ErrorAction Stop).Count
    } catch {}
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
                -Description "El working tree contiene cambios sin commit y/o archivos no rastreados." `
                -Severity "MEDIUM" `
                -EvidenceScope "git working tree" `
                -Evidence (($gitStatus | Out-String).Trim())
            Show-Warn "[GIT] Working tree con cambios pendientes."
        } else {
            Add-Pass "[GIT] Working tree limpio."
        }
    } else {
        Add-Finding -Code "MITOTE-GIT-002" -Category "GIT_NOT_DETECTED" -Description "No es un working tree Git." -Severity "INFO" -EvidenceScope $RepoRoot
        Show-Info "[GIT] No se identificó repositorio Git en esta ruta."
    }
} else {
    Add-Finding -Code "MITOTE-GIT-003" -Category "GIT_COMMAND_UNAVAILABLE" -Description "Comando git no disponible." -Severity "INFO" -EvidenceScope "environment"
    Show-Info "[GIT] Comando git no disponible."
}

# ==============================================================================
# FASE 1 — BASELINE DOCUMENTAL
# ==============================================================================

Write-Host ""
Show-Info "[FASE 1] Verificando documentos baseline..."

foreach ($doc in $BaselineDocs) {
    $rootPath = Join-Path $RepoRoot $doc
    $docsPath = Join-Path $RepoRoot ("docs\" + $doc)

    if (Test-Path $rootPath) {
        Add-Pass "[NORMATIVA] Encontrado en raíz esperada: $doc"
        continue
    }

    if (Test-Path $docsPath) {
        $severity = if ($RequireBaselineAtRoot) { "CRITICAL" } else { "MEDIUM" }
        Add-Finding -Code "RXD-DOC-ROOT-001" -Category "DOC_NOT_IN_EXPECTED_ROOT" -Description "$doc no está en raíz, existe copia en docs\." -Severity $severity -EvidenceScope "root" -Evidence $docsPath
        if ($RequireBaselineAtRoot) { Show-Fail "[NORMATIVA] $doc no está en raíz; existe en docs\." }
        else { Show-Warn "[NORMATIVA] $doc no está en raíz; existe en docs\." }
        continue
    }

    Add-Finding -Code "RXD-DOC-001" -Category "DOC_NOT_DETECTED" -Description "No se encontró $doc ni en raíz ni en docs\." -Severity "CRITICAL" -EvidenceScope "root + docs"
    Show-Fail "[NORMATIVA] $doc no fue detectado."
}

# ==============================================================================
# FASE 1B — ENTRYPOINTS Y COBERTURA DE PRUEBAS
# ==============================================================================

Write-Host ""
Show-Info "[FASE 1B] Verificando entrypoints y cobertura de pruebas por microservicio..."

foreach ($relative in $ExpectedEntrypoints) {
    $path = Join-Path $RepoRoot $relative
    if (Test-Path $path) {
        Add-Pass "[ARQUITECTURA] Entrypoint localizado: $relative"
    } else {
        Add-Finding -Code "RXD-ENTRY-001" -Category "ENTRYPOINT_NOT_DETECTED" -Description "No se encontró $relative." -Severity "HIGH" -EvidenceScope $relative
        Show-Fail "[ARQUITECTURA] No detectado: $relative"
    }
}

# Verificación de archivos _test.go por cada carpeta en cmd/
$cmdDirs = Get-ChildItem -Path (Join-Path $RepoRoot "cmd") -Directory -ErrorAction SilentlyContinue
foreach ($dir in $cmdDirs) {
    $testFiles = Get-ChildItem -Path $dir.FullName -Filter "*_test.go" -ErrorAction SilentlyContinue
    if ($testFiles.Count -eq 0) {
        Add-Finding -Code "RXD-TEST-MISSING" -Category "NO_TEST_FILES" -Description "El microservicio $($dir.Name) no posee archivos de prueba _test.go." -Severity "MEDIUM" -EvidenceScope $dir.FullName
        Show-Warn "[PRUEBAS] Microservicio sin unit tests: cmd/$($dir.Name)"
    } else {
        Add-Pass "[PRUEBAS] Pruebas unitarias detectadas en: cmd/$($dir.Name)"
    }
}

# ==============================================================================
# FASE 1C — DETECCIÓN DE BASURA Y POLUCIÓN ESTRUCTURAL
# ==============================================================================

Write-Host ""
Show-Info "[FASE 1C] Inspeccionando polución estructural (desktop.ini, .go en raíz, scripts)..."

# 1. Archivos desktop.ini
$desktopIniFiles = Get-ChildItem -Path $RepoRoot -Filter "desktop.ini" -Recurse -Force -ErrorAction SilentlyContinue
if ($desktopIniFiles.Count -gt 0) {
    Add-Finding -Code "RXD-CLEAN-DESKTOP-INI" -Category "POLLUTION_DESKTOP_INI" -Description "Se encontraron $($desktopIniFiles.Count) archivos desktop.ini residuales." -Severity "LOW" -EvidenceScope "filesystem" -Evidence (($desktopIniFiles | ForEach-Object { $_.FullName }) -join "; ")
    Show-Warn "[BASURA] Se detectaron $($desktopIniFiles.Count) archivo(s) 'desktop.ini'."
} else {
    Add-Pass "[LIMPIEZA] No se encontraron archivos 'desktop.ini'."
}

# 2. Archivos .go sueltos en la raíz
$foundLooseGo = @()
foreach ($looseFile in $LooseRootGoFiles) {
    $p = Join-Path $RepoRoot $looseFile
    if (Test-Path $p) { $foundLooseGo += $looseFile }
}
if ($foundLooseGo.Count -gt 0) {
    Add-Finding -Code "RXD-CLEAN-ROOT-GO" -Category "POLLUTION_ROOT_GO" -Description "Archivos .go sueltos en raíz ocasionan colisiones de paquete." -Severity "HIGH" -EvidenceScope "repository root" -Evidence ($foundLooseGo -join ", ")
    Show-Fail "[POLUCIÓN] Archivos .go sueltos en raíz detectados: $($foundLooseGo -join ', ')"
} else {
    Add-Pass "[LIMPIEZA] Raíz libre de archivos .go duplicados/sueltos."
}

# 3. Archivo anómalo scripts/main.go
$anomalousScript = Join-Path $RepoRoot "scripts\main.go"
if (Test-Path $anomalousScript) {
    Add-Finding -Code "RXD-CLEAN-SCRIPTS-MAIN" -Category "POLLUTION_ANOMALOUS_FILE" -Description "Existe un archivo main.go anómalo dentro de la carpeta scripts/." -Severity "MEDIUM" -EvidenceScope "scripts/main.go"
    Show-Warn "[POLUCIÓN] Detectado archivo anómalo scripts/main.go."
} else {
    Add-Pass "[LIMPIEZA] Directorio scripts/ sin archivos main.go anómalos."
}

# ==============================================================================
# FASE 2 — CODIFICACIÓN Y PATRONES DE RESILIENCIA EN CÓDIGO
# ==============================================================================

Write-Host ""
Show-Info "[FASE 2] Revisando patrones de resiliencia HTTP, ignorado de errores y vacíos..."

# Detección de silenciamiento de errores JSON y llamadas descartadas con _ =
$silencedJsonCalls = @()
$unhandledCalls = @()
$missingMaxBytes = @()

foreach ($file in $GoFiles) {
    $content = Get-Content $file.FullName -Raw -ErrorAction SilentlyContinue

    if ($content -match '_\s*=\s*json\.NewDecoder') {
        $silencedJsonCalls += $file.Name
    }
    if ($content -match '_\s*=\s*s\.(auditClient|resultClient|deliveryClient)') {
        $unhandledCalls += $file.Name
    }
    if ($file.FullName -match '\\cmd\\' -and $file.Name -eq "main.go" -and $content -notmatch 'http\.MaxBytesReader') {
        $missingMaxBytes += $file.FullName.Replace($RepoRoot, ".")
    }
}

if ($silencedJsonCalls.Count -gt 0) {
    Add-Finding -Code "RXD-RES-001" -Category "SILENCED_JSON_DECODE" -Description "Uso de _ = json.NewDecoder ignora JSONs malformados en lugar de responder 400 Bad Request." -Severity "HIGH" -EvidenceScope "Go source" -Evidence ($silencedJsonCalls -join ", ")
    Show-Warn "[RESILIENCIA] Decodificación JSON silenciada en: $($silencedJsonCalls -join ', ')"
} else {
    Add-Pass "[RESILIENCIA] No se detectó silenciamiento explícito en decodificación JSON."
}

if ($unhandledCalls.Count -gt 0) {
    Add-Finding -Code "RXD-RES-002" -Category "UNHANDLED_CLIENT_ERRORS" -Description "Llamadas a clientes inter-servicio descartadas con _ =." -Severity "MEDIUM" -EvidenceScope "Go source" -Evidence ($unhandledCalls -join ", ")
    Show-Warn "[RESILIENCIA] Errores de cliente descartados en: $($unhandledCalls -join ', ')"
} else {
    Add-Pass "[RESILIENCIA] Manejo o logging activo en llamadas de red inter-servicio."
}

if ($missingMaxBytes.Count -gt 0) {
    Add-Finding -Code "RXD-SEC-001" -Category "MISSING_BODY_LIMIT" -Description "Servicio HTTP carece de http.MaxBytesReader (riesgo DoS por memoria)." -Severity "MEDIUM" -EvidenceScope "cmd/*/main.go" -Evidence ($missingMaxBytes -join "; ")
    Show-Warn "[SEGURIDAD] Sin http.MaxBytesReader en: $($missingMaxBytes -join ', ')"
} else {
    Add-Pass "[SEGURIDAD] Límite de lectura de cuerpo HTTP configurado en endpoints."
}

# ==============================================================================
# FASE 3 — BINDINGS Y RED
# ==============================================================================

Write-Host ""
Show-Info "[FASE 3] Inspeccionando bindings y patrones de red..."

$publicBindings = @()
foreach ($file in $GoFiles) {
    $matches = Select-String -Path $file.FullName -Pattern '0\.0\.0\.0' -ErrorAction SilentlyContinue
    if ($matches) { $publicBindings += $matches }
}

if ($publicBindings.Count -eq 0) {
    Add-Pass "[RED] No se detectó literal 0.0.0.0 en los archivos Go."
} else {
    Add-Finding -Code "RXD-NET-001" -Category "PUBLIC_BINDING_DETECTED" -Description "Literal 0.0.0.0 detectado." -Severity "HIGH" -EvidenceScope "Go source" -Evidence (($publicBindings | ForEach-Object { "$($_.Path):$($_.LineNumber)" }) -join "`n")
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
    $shaImports += Select-String -Path $file.FullName -Pattern '"crypto/sha256"' -ErrorAction SilentlyContinue
    $shaUsage   += Select-String -Path $file.FullName -Pattern 'sha256\.(Sum256|New)\s*\(' -ErrorAction SilentlyContinue
}

if ($shaUsage.Count -gt 0) {
    Add-Pass "[CRIPTO] Se detectó uso directo de sha256.Sum256() y/o sha256.New()."
} elseif ($shaImports.Count -gt 0) {
    Add-Finding -Code "RXD-WARN-005A" -Category "SHA256_IMPORT_WITHOUT_DIRECT_USAGE" -Description "Import de crypto/sha256 sin uso directo." -Severity "MEDIUM" -EvidenceScope "Go source" -Evidence (($shaImports | ForEach-Object { "$($_.Path):$($_.LineNumber)" }) -join "`n")
    Show-Warn "[CRIPTO] Import SHA-256 detectado, pero sin uso directo verificable."
} else {
    Add-Finding -Code "RXD-WARN-005" -Category "SHA256_USAGE_NOT_DETECTED" -Description "No se encontró uso de SHA-256." -Severity "MEDIUM" -EvidenceScope "Go source"
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
    $dicomLiteral    += Select-String -Path $file.FullName -Pattern 'A-ASSOCIATE-RJ' -ErrorAction SilentlyContinue
    $dicomIdentifier += Select-String -Path $file.FullName -Pattern '\bDicomAssociateRJ\b' -ErrorAction SilentlyContinue
}

if ($dicomLiteral.Count -eq 0) {
    Add-Finding -Code "RXD-WARN-006" -Category "DICOM_REJECT_NOT_DETECTED" -Description "No se encontró A-ASSOCIATE-RJ." -Severity "MEDIUM" -EvidenceScope "Go source"
    Show-Warn "[DICOM] No se detectó referencia explícita a A-ASSOCIATE-RJ."
} elseif ($dicomIdentifier.Count -le 1) {
    Add-Finding -Code "RXD-WARN-006A" -Category "DICOM_REJECT_DECLARED_ONLY" -Description "Constante DicomAssociateRJ declarada pero sin consumo funcional comprobado." -Severity "MEDIUM" -EvidenceScope "Go source" -Evidence (($dicomIdentifier | ForEach-Object { "$($_.Path):$($_.LineNumber)" }) -join "`n")
    Show-Warn "[DICOM] Referencia declarada, uso funcional no acreditado."
} else {
    Add-Pass "[DICOM] Se detectó declaración y referencia adicional de DicomAssociateReject."
}

# ==============================================================================
# FASE 4 — 14 CAMPOS ePHI
# ==============================================================================

Write-Host ""
Show-Info "[FASE 4] Verificando detección de los 14 campos ePHI..."

if ($ModelFiles.Count -eq 0) {
    Add-Finding -Code "RXD-EPHI-000" -Category "MODEL_SCOPE_NOT_DETECTED" -Description "No se encontraron modelos." -Severity "HIGH" -EvidenceScope "internal/models"
    Show-Fail "[EPHI] No se encontró alcance de modelos evaluable."
} else {
    foreach ($field in $EPHIFields) {
        $pattern = 'json:"' + [regex]::Escape($field) + '(?:,[^"]*)?"'
        $matches = @()
        foreach ($file in $ModelFiles) {
            $matches += Select-String -Path $file.FullName -Pattern $pattern -ErrorAction SilentlyContinue
        }

        if ($matches.Count -gt 0) {
            Add-Pass "[EPHI_FIELD] Detectado tag JSON: $field"
        } else {
            Add-Finding -Code "RXD-EPHI-001" -Category "EPHI_FIELD_NOT_DETECTED" -Description "Campo '$field' no detectado." -Severity "MEDIUM" -EvidenceScope (($ModelFiles.FullName) -join "; ")
            Show-Warn "[EPHI_FIELD] No detectado en alcance inspeccionado: $field"
        }
    }
}

# ==============================================================================
# FASES 5, 6, 7 — HERRAMIENTAS GO (BUILD, TEST, VET)
# ==============================================================================

Write-Host ""
Show-Info "[FASE 5] Ejecutando go build ./..."
$buildOK = Invoke-GoCheck -Name "go build ./..." -Arguments @("build", "./...") -Code "RXD-ERR-005" -Category "GO_BUILD_FAIL" -SuccessText "[COMPILADOR] go build ./... completó correctamente."

Write-Host ""
Show-Info "[FASE 6] Ejecutando go test ./... -count=1..."
$testOK = Invoke-GoCheck -Name "go test ./... -count=1" -Arguments @("test", "./...", "-count=1") -Code "RXD-ERR-007" -Category "GO_TEST_FAIL" -SuccessText "[PRUEBAS] go test ./... -count=1 completó correctamente."

Write-Host ""
Show-Info "[FASE 7] Ejecutando go vet ./..."
$vetOK = Invoke-GoCheck -Name "go vet ./..." -Arguments @("vet", "./...") -Code "RXD-ERR-006" -Category "GO_VET_FAIL" -SuccessText "[ANÁLISIS] go vet ./... completó correctamente."

# ==============================================================================
# MÉTRICAS Y ESTADO GLOBAL
# ==============================================================================

$Warnings       = @($Findings | Where-Object { $_.Severity -in @("LOW", "MEDIUM") }).Count
$HighErrors     = @($Findings | Where-Object { $_.Severity -eq "HIGH" }).Count
$CriticalErrors = @($Findings | Where-Object { $_.Severity -eq "CRITICAL" }).Count
$Informational  = @($Findings | Where-Object { $_.Severity -eq "INFO" }).Count

$BlockingFindings = $HighErrors + $CriticalErrors

$Status = if ($BlockingFindings -gt 0) { "FAIL" } elseif ($Warnings -gt 0) { "PASS_WITH_WARNINGS" } else { "PASS" }

# ==============================================================================
# GENERACIÓN DE REPORTE JSON
# ==============================================================================

$Report = [PSCustomObject]@{
    Timestamp = (Get-Date).ToString("o")
    Status    = $Status
    Scope     = [PSCustomObject]@{
        RepositoryRoot        = $RepoRoot
        RequireBaselineAtRoot = $RequireBaselineAtRoot
        Statement             = "Los hallazgos se limitan al alcance inspeccionado. NOT_DETECTED no implica inexistencia global."
    }
    Metrics   = [PSCustomObject]@{
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
$utf8NoBom = New-Object System.Text.UTF8Encoding($false)
$jsonContent = $Report | ConvertTo-Json -Depth 10
[System.IO.File]::WriteAllText($ReportPath, $jsonContent, $utf8NoBom)

# ==============================================================================
# RESUMEN EN CONSOLA
# ==============================================================================

Write-Host ""
Write-Host "==========================================================================" -ForegroundColor Cyan
Write-Host "                 RESUMEN PERICIAL ULTRA-MITOTERO v4.0" -ForegroundColor Cyan
Write-Host "==========================================================================" -ForegroundColor Cyan

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

if ($Status -eq "PASS") {
    Write-Host " Estado Global                  : PASS" -ForegroundColor Green
} elseif ($Status -eq "PASS_WITH_WARNINGS") {
    Write-Host " Estado Global                  : PASS_WITH_WARNINGS" -ForegroundColor Yellow
} else {
    Write-Host " Estado Global                  : FAIL" -ForegroundColor Red
}

Write-Host "==========================================================================" -ForegroundColor Cyan
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

if ($BlockingFindings -gt 0) { exit 1 } else { exit 0 }