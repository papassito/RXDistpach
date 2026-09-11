# ==============================================================================
# RXDISPATCH — AUDITORÍA PROFUNDA INTEGRAL DE CÓDIGO Y CONTRATOS GO
# ==============================================================================

Clear-Host
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host "   RXDISPATCH - AUDITORIA PROFUNDA DE CODIGO Y CONTRATOS    " -ForegroundColor Cyan
Write-Host "============================================================" -ForegroundColor Cyan

$ROOT = Get-Location
$PuntajeSalud = 100
$Hallazgos = @()

function Registrar-Hallazgo {
    param(
        [string]$Severidad,
        [int]$Puntos,
        [string]$Categoria,
        [string]$Archivo,
        [int]$Linea,
        [string]$Detalle
    )
    $script:PuntajeSalud -= $Puntos
    $script:Hallazgos += [PSCustomObject]@{
        Severidad = $Severidad
        Categoria = $Categoria
        Ubicacion = "$Archivo:$Linea"
        Detalle   = $Detalle
    }
}

# 1. RASTREO DE CONTRATOS Y STRUCTS
Write-Host "`n[1/4] Auditando alineacion de contratos y modelos..." -ForegroundColor Yellow

$archivosGoCmd = Get-ChildItem -Recurse -Path "cmd" -Filter "*.go" -ErrorAction SilentlyContinue

foreach ($file in $archivosGoCmd) {
    $lineas = Get-Content $file.FullName
    $relPath = $file.FullName.Replace("$ROOT\", "")

    for ($i = 0; $i -lt $lineas.Count; $i++) {
        $numLinea = $i + 1
        $lineaTexto = $lineas[$i]

        if ($lineaTexto -match '\bTimestamp\s*:' -and $lineaTexto -notmatch 'EventTimestamp') {
            Registrar-Hallazgo "CRITICO" 10 "CONTRATO" $relPath $numLinea "Uso del campo descontinuado 'Timestamp:' en lugar de 'EventTimestamp:'"
        }

        if ($lineaTexto -match '\bActor\s*:') {
            Registrar-Hallazgo "CRITICO" 10 "CONTRATO" $relPath $numLinea "Uso del campo descontinuado 'Actor:' en lugar de 'UserID:'"
        }
    }
}

# 2. AUDITORIA DE SEGURIDAD Y HARDCODING
Write-Host "[2/4] Buscando rutas de disco duro e informacion sensible..." -ForegroundColor Yellow

$todosArchivosGo = Get-ChildItem -Recurse -Filter "*.go" -ErrorAction SilentlyContinue |
    Where-Object { $_.FullName -notmatch '\\vendor\\' }

foreach ($file in $todosArchivosGo) {
    $lineas = Get-Content $file.FullName
    $relPath = $file.FullName.Replace("$ROOT\", "")

    for ($i = 0; $i -lt $lineas.Count; $i++) {
        $numLinea = $i + 1
        $lineaTexto = $lineas[$i]

        if ($lineaTexto -match '[A-Za-z]:\\[^\s"']+') {
            Registrar-Hallazgo "ADVERTENCIA" 5 "HARDCODE_PATH" $relPath $numLinea "Ruta absoluta de disco detectada"
        }

        if ($lineaTexto -match '(?i)(password|secret|apikey)\s*=\s*"[^"]{4,}"') {
            Registrar-Hallazgo "CRITICO" 15 "SEGURIDAD" $relPath $numLinea "Posible credencial o secreto hardcodeado en codigo"
        }
    }
}

# 3. AUDITORIA DE DEUDA TECNICA
Write-Host "[3/4] Evaluando deuda tecnica y llamadas inseguras..." -ForegroundColor Yellow

foreach ($file in $todosArchivosGo) {
    $lineas = Get-Content $file.FullName
    $relPath = $file.FullName.Replace("$ROOT\", "")

    for ($i = 0; $i -lt $lineas.Count; $i++) {
        $numLinea = $i + 1
        $lineaTexto = $lineas[$i]

        if ($lineaTexto -match '(?i)\b(TODO|FIXME)\b') {
            Registrar-Hallazgo "INFO" 1 "DEUDA_TECNICA" $relPath $numLinea "Marca de trabajo pendiente"
        }
        if ($lineaTexto -match '\bpanic\(') {
            Registrar-Hallazgo "ADVERTENCIA" 5 "ROBUSTEZ" $relPath $numLinea "Uso explicito de panic(). Se recomienda manejo de errores con error."
        }
    }
}

# 4. VERIFICACION CON GO TOOLCHAIN
Write-Host "[4/4] Comprobando analisis estatico oficial con Go Toolchain..." -ForegroundColor Yellow

if (Get-Command "go" -ErrorAction SilentlyContinue) {
    $vetResult = go vet ./... 2>&1
    if ($LASTEXITCODE -ne 0) {
        Registrar-Hallazgo "CRITICO" 20 "COMPILACION" "Proyecto" 0 "go vet detecto errores de analisis estatico"
    } else {
        Write-Host "   [OK] 'go vet ./...' se ejecuto sin errores." -ForegroundColor Green
    }
} else {
    Write-Host "   [INFO] Compilador 'go' no detectado en el PATH. Se omite go vet." -ForegroundColor DarkGray
}

# INFORME FINAL
if ($PuntajeSalud -lt 0) { $PuntajeSalud = 0 }

Write-Host "`n============================================================" -ForegroundColor Cyan
Write-Host "           INFORME DE AUDITORIA PROFUNDA - RXDISPATCH       " -ForegroundColor Cyan
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host " SALUD GENERAL DEL CODIGO : $PuntajeSalud / 100" -ForegroundColor $(if($PuntajeSalud -eq 100){"Green"}elseif($PuntajeSalud -gt 70){"Yellow"}else{"Red"})
Write-Host " TOTAL DE HALLAZGOS       : $($Hallazgos.Count)" -ForegroundColor $(if($Hallazgos.Count -eq 0){"Green"}else{"Red"})
Write-Host "============================================================" -ForegroundColor Cyan

if ($Hallazgos.Count -gt 0) {
    Write-Host "`nDETALLE DE HALLAZGOS:" -ForegroundColor White
    $Hallazgos | Format-Table -AutoSize
} else {
    Write-Host "`n[OK] EL PROYECTO CUMPLE AL 100% CON TODOS LOS CONTRATOS Y ESTANDARES." -ForegroundColor Green
}
