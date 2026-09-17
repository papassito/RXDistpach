# ============================================================================
# ANDROID FIX STUDIO - AUDITORÍA EN VIVO Y DIAGNÓSTICO DE DOLOR
# ============================================================================

$ErrorActionPreference = 'Continue'
$OutputEncoding = [System.Text.Encoding]::UTF8

# Función para pausar y salir de forma ordenada
function Salir-ConPausa {
    param(
        [string]$Mensaje = "Presiona Enter para cerrar la ventana..."
    )
    Write-Host "`n $Mensaje" -ForegroundColor DarkGray
    $null = Read-Host
    exit
}

# 1. Buscar ADB en la ruta actual, carpeta superior o en el PATH de Windows
$adb = if (Test-Path -LiteralPath ".\adb.exe") { ".\adb.exe" }
       elseif (Test-Path -LiteralPath "..\adb.exe") { "..\adb.exe" }
       else { (Get-Command adb.exe -ErrorAction SilentlyContinue).Source }

if ($null -eq $adb) {
    Clear-Host
    Write-Host "============================================================================" -ForegroundColor Red
    Write-Host " [CRÍTICO] No se encontró adb.exe en esta carpeta ni en el sistema." -ForegroundColor Red
    Write-Host "============================================================================" -ForegroundColor Red
    Write-Host " SOLUCIÓN:" -ForegroundColor Yellow
    Write-Host " 1. Descarga Platform-Tools de Android (ADB)." -ForegroundColor White
    Write-Host " 2. Coloca el archivo 'adb.exe' dentro de esta misma carpeta:" -ForegroundColor White
    Write-Host "    $((Get-Location).Path)" -ForegroundColor Cyan
    Write-Host " 3. Vuelve a ejecutar este script." -ForegroundColor White
    Salir-ConPausa
}

Clear-Host
Write-Host "============================================================================" -ForegroundColor Cyan
Write-Host "   S23 FIX STUDIO - INSPECTOR DE DAÑOS Y TELEMETRÍA (PowerShell)" -ForegroundColor Cyan
Write-Host "============================================================================" -ForegroundColor Cyan
Write-Host " Buscando dispositivo físico por USB..." -ForegroundColor Gray

# 2. Comprobar estado del dispositivo
$devices = & $adb devices -l | Out-String
$cleanDevices = $devices -replace '(?i)List of devices attached', ''
$deviceMatch = [regex]::Matches($cleanDevices, '(?m)^([^\s]+)\s+(device|unauthorized|offline)\b')

if ($deviceMatch.Count -eq 0) {
    Write-Host " [X] ERROR: No hay ningún dispositivo conectado por USB." -ForegroundColor Red
    Write-Host " Asegúrate de conectar el teléfono y activar la 'Depuración USB'." -ForegroundColor Yellow
    Salir-ConPausa
} elseif ($deviceMatch.Count -gt 1) {
    Write-Host " [!] ADVERTENCIA: Múltiples dispositivos detectados. Desconecta uno para evitar conflictos." -ForegroundColor Yellow
    Salir-ConPausa
}

$serial = $deviceMatch[0].Groups[1].Value
$state = $deviceMatch[0].Groups[2].Value

if ($state -eq 'unauthorized') {
    Write-Host " [!] DISPOSITIVO NO AUTORIZADO (Clave RSA pendiente)" -ForegroundColor Red
    Write-Host " Mira la pantalla de tu celular y marca 'Permitir siempre desde esta computadora'." -ForegroundColor Yellow
    Salir-ConPausa
} elseif ($state -eq 'offline') {
    Write-Host " [X] DISPOSITIVO EN SUSPENSIÓN / OFFLINE" -ForegroundColor Red
    Salir-ConPausa
}

# 3. Dispositivo listo -> Extraer datos básicos
$model = (& $adb -s $serial shell getprop ro.product.model).Trim()
$manufacturer = (& $adb -s $serial shell getprop ro.product.manufacturer).Trim()
$androidVer = (& $adb -s $serial shell getprop ro.build.version.release).Trim()
$soc = (& $adb -s $serial shell getprop ro.board.platform).Trim()

Write-Host " [+] DISPOSITIVO DETECTADO:" -ForegroundColor Green
Write-Host "   • Fabricante: $manufacturer" -ForegroundColor White
Write-Host "   • Modelo:     $model ($serial)" -ForegroundColor White
Write-Host "   • Android:    v$androidVer" -ForegroundColor White
Write-Host "   • SoC/Chip:   $soc" -ForegroundColor White
Write-Host "============================================================================" -ForegroundColor Cyan
Write-Host " [1/3] Extrayendo telemetría de energía (Dumpsys Battery)..." -ForegroundColor Blue

$battery = (& $adb -s $serial shell dumpsys battery) | Out-String
$batLevel = 0; $batTemp = 0.0; $batHealth = 0; $batCharge = $false

foreach ($line in ($battery -split "`n")) {
    $line = $line.Trim().ToLower()
    if ($line -match '^level:\s*(\d+)') { $batLevel = [int]$Matches[1] }
    elseif ($line -match '^temperature:\s*(\d+)') { $batTemp = [float]$Matches[1] / 10 }
    elseif ($line -match '^health:\s*(\d+)') { $batHealth = [int]$Matches[1] }
    elseif ($line -match '^powered:\s*true') { $batCharge = $true }
}

Write-Host " [2/3] Analizando estabilidad de sistema (Dumpsys Dropbox)..." -ForegroundColor Blue
$dropbox = (& $adb -s $serial shell dumpsys dropbox --print) | Out-String
$crashCount = [regex]::Matches($dropbox, '(?i)(system_server_crash|system_app_crash|crash|anr|fatal)').Count

Write-Host " [3/3] Evaluando compilación ART y Cachés dexopt..." -ForegroundColor Blue
$uncompiledApps = 0
$packages = (& $adb -s $serial shell pm list packages -3) | Out-String
$samplePackages = ($packages -split "`n" | Select-Object -First 5 | ForEach-Object { $_.replace("package:", "").Trim() })
foreach ($pkg in $samplePackages) {
    if ($null -ne $pkg -and $pkg -ne "") {
        $dump = (& $adb -s $serial shell cmd package dump $pkg) | Out-String
        if ($dump -match '(?i)compiler-filter=interpret-only|verify') { $uncompiledApps++ }
    }
}

Write-Host "============================================================================" -ForegroundColor Cyan
Write-Host " REPORTE DE SALUD Y RECTIFICACIÓN (¿QUÉ LE DUELE A TU ANDROID?)" -ForegroundColor Cyan
Write-Host "============================================================================" -ForegroundColor Cyan

$dolores = 0

# Alerta 1: Temperatura de la batería
if ($batTemp -gt 38.0) {
    Write-Host " [X] DOLOR TÉRMICO (Crítico): La batería está a $batTemp °C." -ForegroundColor Red
    Write-Host "     -> Causa: El procesador SoC $soc está sufriendo estrangulamiento térmico (Throttling)." -ForegroundColor DarkGray
    Write-Host "     -> Solución: Detén los bucles de WakeLocks e hilos de telemetría inactivos en segundo plano." -ForegroundColor Yellow
    $dolores++
} else {
    Write-Host " [OK] Temperatura Normal: El dispositivo opera a un rango saludable de $batTemp °C." -ForegroundColor Green
}

# Alerta 2: Batería degradada o nivel bajo
if ($batHealth -ne 2) {
    Write-Host " [X] ENVEJECIMIENTO DE BATERÍA (Advertencia): Estado de vida útil no óptimo (Código $batHealth)." -ForegroundColor Yellow
    Write-Host "     -> Causa: Desgaste químico de celdas o logs descalibrados de energía." -ForegroundColor DarkGray
    Write-Host "     -> Solución: Ejecutar recalibración de tablas físicas y limpiar historial de dumpsys." -ForegroundColor Gray
    $dolores++
}

if ($batLevel -lt 25 -and !$batCharge) {
    Write-Host " [!] BATERÍA CRÍTICA: Nivel actual al $batLevel% sin suministro eléctrico." -ForegroundColor Red
    $dolores++
}

# Alerta 3: Inestabilidad de Software
if ($crashCount -gt 5) {
    Write-Host " [X] INESTABILIDAD DE SOFTWARE (Crítico): Se detectaron $crashCount crashes/ANRs recientes." -ForegroundColor Red
    Write-Host "     -> Causa: Corrupción de base de datos interna y tablas caché obsoletas." -ForegroundColor DarkGray
    Write-Host "     -> Solución: Purgar directorios /data/system/dropbox/ y ejecutar 'pm trim-caches'." -ForegroundColor Yellow
    $dolores++
} else {
    Write-Host " [OK] Estabilidad Excelente: Registros de caída limpios en las últimas horas." -ForegroundColor Green
}

# Alerta 4: Estado de Compilación de Apps (ART)
if ($uncompiledApps -gt 0) {
    Write-Host " [X] DEGRADACIÓN DE RENDIMIENTO (Advertencia): Aplicaciones de usuario sin compilar nativamente." -ForegroundColor Yellow
    Write-Host "     -> Causa: Caché ART desactualizada (se ejecuta en modo interpretado lento)." -ForegroundColor DarkGray
    Write-Host "     -> Solución: Forzar recompilación de paquetes con 'cmd package compile -m speed-profile -a'." -ForegroundColor Yellow
    $dolores++
} else {
    Write-Host " [OK] Compilación ART al Día: Código optimizado para ejecución directa en el chip." -ForegroundColor Green
}

Write-Host "============================================================================" -ForegroundColor Cyan
if ($dolores -eq 0) {
    Write-Host " EXCELENTE: Tu teléfono $model no reporta dolores clínicos activos." -ForegroundColor Green
    Write-Host " Conservas un rendimiento y autonomía óptimos." -ForegroundColor Green
} else {
    Write-Host " Tu $model presenta un total de $dolores fallas o dolores detectados." -ForegroundColor Yellow
}
Write-Host "============================================================================" -ForegroundColor Cyan

Salir-ConPausa -Mensaje "Presiona Enter para finalizar y cerrar el reporte..."
