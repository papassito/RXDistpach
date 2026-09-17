# test_integration_reader_audit.ps1
#
# Este script realiza la primera prueba de integración del sistema:
# 1. Envía una petición POST a rx-reader.
# 2. rx-reader procesa la petición y envía un evento de auditoría a rx-audit.
#
# Para verificar el éxito, observa las ventanas de terminal de ambos servicios.

# Configurar la codificación de salida de PowerShell a UTF-8 para mostrar correctamente los acentos.
# Esta es la solución final del lado del cliente ahora que el servidor envía el encabezado correcto.
$OutputEncoding = [System.Text.Encoding]::UTF8

Write-Host "Performing integration test: Reader -> Audit..." -ForegroundColor Cyan

# Definir el endpoint y el cuerpo de la petición
$uri = "http://127.0.0.1:8085/reader/analyze"
$body = @{
    studyId          = "INTEGRATION-TEST-001"
    studyType        = "DX"
    anatomicalRegion = "WRIST"
} | ConvertTo-Json

Write-Host "Sending POST request to $uri"

try {
    # Enviar la petición
    # Se usa Invoke-WebRequest en lugar de Invoke-RestMethod para compatibilidad con versiones
    # antiguas de PowerShell y para obtener acceso a las cabeceras de la respuesta.
    $responseRaw = Invoke-WebRequest -Uri $uri -Method Post -Body $body -ContentType "application/json"
    
    Write-Host "--- DIAGNOSTIC: RESPONSE HEADERS FROM RX-READER ---" -ForegroundColor Yellow
    $responseRaw.Headers

    Write-Host "--- RESPONSE FROM RX-READER ---" -ForegroundColor Green
    # El contenido de la respuesta se decodifica correctamente gracias al $OutputEncoding y al charset del servidor.
    $responseRaw.Content
    
    Write-Host "SUCCESS: Request completed. Now, check the terminal windows for logs:" -ForegroundColor Green
    Write-Host "1. The 'rx-reader' window should show 'Received analysis request...'"
    Write-Host "2. The 'rx-audit' window should show 'AUDIT EVENT RECEIVED...'"
} catch {
    Write-Host "--- FAILED TO CONNECT TO RX-READER ---" -ForegroundColor Red
    Write-Host "Error: $($_.Exception.Message)"
    Write-Host "Please ensure the 'rx-reader' service is running on port 8085."
}