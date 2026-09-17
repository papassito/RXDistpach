$utf8NoBom = New-Object System.Text.UTF8Encoding($false)

$auditScriptCode = @'
<#
================================================================================
 RX DISPATCH :: MOTOR DE AUDITORÍA PERICIAL Y ESCANEO PROFUNDO v4.1
 "Ultra-Mitotero: Codificación Limpia y Resiliencia HTTP"
================================================================================
#>

param(
    [string]$Root = ".",
    [bool]$RequireBaselineAtRoot = $true
)

$ErrorActionPreference = "Continue"
$RepoRoot = (Resolve-Path $Root).Path
Set-Location $RepoRoot

$BaselineDocs = @("README.md", "ARCHITECTURE.md", "REQUIREMENTS.md", "CONTRACTS.md", "MAP.md", "SECURITY.md")
$ExpectedEntrypoints = @(
    "cmd/audit/main.go", "cmd/delivery/main.go", "cmd/gateway/main.go",
    "cmd/image/main.go", "cmd/reader/main.go", "cmd/result/main.go",
    "cmd/security/main.go", "cmd/storage/main.go", "cmd/study/main.go"
)
$EPHIFields = @(
    "event_timestamp", "event_action", "event_outcome", "event_error_code",
    "user_id", "source_ip", "source_ae_title", "destination_ip",
    "destination_ae_title", "patient_id", "study_instance_uid",
    "accession_number", "number_of_instances", "security_tls_status"
)
$LooseRootGoFiles = @("main.go", "models.go", "contracts.go", "constants.go", "client.go")

$Findings = @()
$PassedChecks = 0

function Add-Finding {
    param([string]$Code, [string]$Category, [string]$Description, [string]$Severity, [string]$EvidenceScope = "", [string]$Evidence = "")
    $script:Findings += [PSCustomObject]@{
        Code = $Code; Category = $Category; Description = $Description
        Severity = $Severity; EvidenceScope = $EvidenceScope; Evidence = $Evidence
        Timestamp = (Get-Date).ToString("o")
    }
}

function Add-Pass { param([string]$Message) $script:PassedChecks++; Write-Host "[OK] $Message" -ForegroundColor Green }
function Show-Warn { param([string]$Message) Write-Host "[WARN] $Message" -ForegroundColor Yellow }
function Show-Fail { param([string]$Message) Write-Host "[FAIL] $Message" -ForegroundColor Red }
function Show-Info { param([string]$Message) Write-Host "[INFO] $Message" -ForegroundColor Cyan }

function Invoke-GoCheck {
    param([string]$Name, [string[]]$Arguments, [string]$Code, [string]$Category, [string]$SuccessText)
    try {
        $output = & go @Arguments 2>&1
        $exitCode = $LASTEXITCODE
    } catch {
        $output = @($_.Exception.Message)
        $exitCode = 1
    }
    if ($exitCode -eq 0) {
        Add-Pass $SuccessText
        return $true
    }
    $details = ($output | Out-String).Trim()
    Add-Finding -Code $Code -Category $Category -Description "$Name declino ejecucion." -Severity "HIGH" -EvidenceScope "go workspace" -Evidence $details
    Show-Fail "$Name termino con codigo $exitCode."
    if ($details) { Write-Host $details -ForegroundColor DarkYellow }
    return $false
}

Clear-Host
Write-Host "==========================================================================" -ForegroundColor Cyan
Write-Host " 🔬 RX DISPATCH :: MOTOR DE AUDITORÍA PERICIAL v4.1" -ForegroundColor Cyan
Write-Host "==========================================================================" -ForegroundColor Cyan

$AllFiles = Get-ChildItem -Path $RepoRoot -Recurse -File -Force -ErrorAction SilentlyContinue |
    Where-Object { $_.FullName -notmatch '\\\.git\\' -and $_.FullName -notmatch '\\bin\\' -and $_.FullName -notmatch '\\vendor\\' }

$GoFiles = @($AllFiles | Where-Object { $_.Extension -eq ".go" })
$ModelFiles = @($GoFiles | Where-Object { $_.FullName -match '\\internal\\models\\' })
if ($ModelFiles.Count -eq 0) { $ModelFiles = @($GoFiles | Where-Object { $_.FullName -match '\\internal\\' }) }

$TotalLOC = 0
foreach ($file in $GoFiles) { try { $TotalLOC += (Get-Content $file.FullName -ErrorAction Stop).Count } catch {} }

# FASE 0 - GIT
Write-Host ""
Show-Info "[FASE 0] Inspeccionando Git..."
if (Get-Command git -ErrorAction SilentlyContinue) {
    if ((& git rev-parse --is-inside-work-tree 2>$null) -eq "true") {
        $gitStatus = @(& git status --porcelain 2>$null)
        if ($gitStatus.Count -gt 0) {
            Add-Finding -Code "MITOTE-GIT-001" -Category "GIT_DIRTY" -Description "Cambios sin commit." -Severity "MEDIUM" -EvidenceScope "git" -Evidence (($gitStatus | Out-String).Trim())
            Show-Warn "[GIT] Working tree con cambios pendientes."
        } else { Add-Pass "[GIT] Working tree limpio." }
    } else { Show-Info "[GIT] No es repo Git." }
} else { Show-Info "[GIT] Comando git no disponible." }

# FASE 1 - DOCUMENTOS
Write-Host ""
Show-Info "[FASE 1] Verificando documentos baseline..."
foreach ($doc in $BaselineDocs) {
    if (Test-Path (Join-Path $RepoRoot $doc)) { Add-Pass "[NORMATIVA] Encontrado en raiz: $doc"; continue }
    if (Test-Path (Join-Path $RepoRoot ("docs\" + $doc))) {
        Add-Finding -Code "RXD-DOC-ROOT-001" -Category "DOC_NOT_IN_EXPECTED_ROOT" -Description "$doc existe en docs\." -Severity "CRITICAL" -EvidenceScope "root"
        Show-Fail "[NORMATIVA] $doc no esta en raiz; existe en docs\."
        continue
    }
    Add-Finding -Code "RXD-DOC-001" -Category "DOC_NOT_DETECTED" -Description "$doc ausente." -Severity "CRITICAL"
    Show-Fail "[NORMATIVA] $doc no fue detectado."
}

# FASE 1B - ENTRYPOINTS & TESTS
Write-Host ""
Show-Info "[FASE 1B] Verificando entrypoints y pruebas..."
foreach ($relative in $ExpectedEntrypoints) {
    if (Test-Path (Join-Path $RepoRoot $relative)) { Add-Pass "[ARQUITECTURA] Entrypoint localizado: $relative" }
    else { Add-Finding -Code "RXD-ENTRY-001" -Category "ENTRYPOINT_NOT_DETECTED" -Description "Falta $relative." -Severity "HIGH"; Show-Fail "[ARQUITECTURA] No detectado: $relative" }
}
foreach ($dir in (Get-ChildItem -Path (Join-Path $RepoRoot "cmd") -Directory -ErrorAction SilentlyContinue)) {
    $tests = Get-ChildItem -Path $dir.FullName -Filter "*_test.go" -ErrorAction SilentlyContinue
    if ($tests.Count -eq 0) { Add-Finding -Code "RXD-TEST-MISSING" -Category "NO_TEST_FILES" -Description "cmd/$($dir.Name) sin pruebas." -Severity "MEDIUM"; Show-Warn "[PRUEBAS] Sin unit tests: cmd/$($dir.Name)" }
    else { Add-Pass "[PRUEBAS] Unit tests detectados en: cmd/$($dir.Name)" }
}

# FASE 1C - POLUCION
Write-Host ""
Show-Info "[FASE 1C] Verificando limpieza del repositorio..."
$desktopIniFiles = Get-ChildItem -Path $RepoRoot -Filter "desktop.ini" -Recurse -Force -ErrorAction SilentlyContinue
if ($desktopIniFiles.Count -gt 0) { Add-Finding -Code "RXD-CLEAN-DESKTOP-INI" -Category "POLLUTION" -Description "desktop.ini detectados." -Severity "LOW"; Show-Warn "[BASURA] Detectados desktop.ini" }
else { Add-Pass "[LIMPIEZA] Sin archivos desktop.ini" }

$foundLoose = @()
foreach ($f in $LooseRootGoFiles) { if (Test-Path (Join-Path $RepoRoot $f)) { $foundLoose += $f } }
if ($foundLoose.Count -gt 0) { Add-Finding -Code "RXD-CLEAN-ROOT-GO" -Category "POLLUTION" -Description "Archivos .go sueltos en raiz." -Severity "HIGH"; Show-Fail "[POLUCION] .go sueltos en raiz: $($foundLoose -join ', ')" }
else { Add-Pass "[LIMPIEZA] Raiz libre de archivos .go sueltos" }

if (Test-Path (Join-Path $RepoRoot "scripts\main.go")) { Add-Finding -Code "RXD-CLEAN-SCRIPTS" -Category "POLLUTION" -Description "scripts/main.go anomalo." -Severity "MEDIUM"; Show-Warn "[POLUCION] scripts/main.go detectado" }
else { Add-Pass "[LIMPIEZA] Directorio scripts sin main.go" }

# FASE 2 - RESILIENCIA
Write-Host ""
Show-Info "[FASE 2] Analizando resiliencia y seguridad HTTP..."
$missingMaxBytes = @()
foreach ($file in $GoFiles) {
    $content = Get-Content $file.FullName -Raw -ErrorAction SilentlyContinue
    if ($file.FullName -match '\\cmd\\' -and $file.Name -eq "main.go" -and $content -notmatch 'http\.MaxBytesReader') {
        $missingMaxBytes += $file.FullName.Replace($RepoRoot, ".")
    }
}
if ($missingMaxBytes.Count -gt 0) { Add-Finding -Code "RXD-SEC-001" -Category "NO_BODY_LIMIT" -Description "Sin http.MaxBytesReader" -Severity "MEDIUM" -EvidenceScope "cmd"; Show-Warn "[SEGURIDAD] Sin http.MaxBytesReader en: $($missingMaxBytes -join ', ')" }
else { Add-Pass "[SEGURIDAD] Limite MaxBytesReader configurado en endpoints HTTP" }

# FASE 4 - ePHI
Write-Host ""
Show-Info "[FASE 4] Verificando 14 campos ePHI..."
foreach ($field in $EPHIFields) {
    $pattern = 'json:"' + [regex]::Escape($field) + '"'
    $match = $false
    foreach ($file in $ModelFiles) {
        if ((Select-String -Path $file.FullName -Pattern $pattern -ErrorAction SilentlyContinue)) { $match = $true; break }
    }
    if ($match) { Add-Pass "[EPHI_FIELD] Detectado: $field" }
    else { Add-Finding -Code "RXD-EPHI-001" -Category "EPHI_MISSING" -Description "Falta $field" -Severity "MEDIUM"; Show-Warn "[EPHI_FIELD] No detectado: $field" }
}

# FASES 5, 6, 7 - GO CLI
Write-Host ""
Show-Info "[FASE 5] Ejecutando go build ./..."
$buildOK = Invoke-GoCheck -Name "go build ./..." -Arguments @("build", "./...") -Code "RXD-ERR-005" -Category "GO_BUILD_FAIL" -SuccessText "[COMPILADOR] go build ./... PASS"

Write-Host ""
Show-Info "[FASE 6] Ejecutando go test ./... -count=1..."
$testOK = Invoke-GoCheck -Name "go test ./..." -Arguments @("test", "./...", "-count=1") -Code "RXD-ERR-007" -Category "GO_TEST_FAIL" -SuccessText "[PRUEBAS] go test ./... -count=1 PASS"

Write-Host ""
Show-Info "[FASE 7] Ejecutando go vet ./..."
$vetOK = Invoke-GoCheck -Name "go vet ./..." -Arguments @("vet", "./...") -Code "RXD-ERR-006" -Category "GO_VET_FAIL" -SuccessText "[ANALISIS] go vet ./... PASS"

# METRICAS
$Warnings = @($Findings | Where-Object { $_.Severity -in @("LOW", "MEDIUM") }).Count
$HighErrors = @($Findings | Where-Object { $_.Severity -in @("HIGH", "CRITICAL") }).Count

Write-Host ""
Write-Host "==========================================================================" -ForegroundColor Cyan
Write-Host " RESULTADOS AUDITORIA v4.1" -ForegroundColor Cyan
Write-Host "==========================================================================" -ForegroundColor Cyan
Write-Host (" go build ./... : {0}" -f $(if ($buildOK) { "PASS" } else { "FAIL" }))
Write-Host (" go test ./...  : {0}" -f $(if ($testOK) { "PASS" } else { "FAIL" }))
Write-Host (" go vet ./...   : {0}" -f $(if ($vetOK) { "PASS" } else { "FAIL" }))
Write-Host "==========================================================================" -ForegroundColor Cyan
'@

[System.IO.File]::WriteAllText((Resolve-Path ".\audit_ultra-mitotero.ps1"), $auditScriptCode.Trim(), $utf8NoBom)
Write-Host "✅ audit_ultra-mitotero.ps1 saneado en UTF-8 puro." -ForegroundColor Green

.\audit_ultra-mitotero.ps1