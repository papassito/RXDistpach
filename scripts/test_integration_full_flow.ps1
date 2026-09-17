# test_integration_full_flow.ps1
#
# Este script realiza una prueba de integración de extremo a extremo (E2E)
# para el flujo de datos principal del sistema:
# 1. Envía una petición POST a rx-reader.
# 2. rx-reader la procesa y la envía a rx-result.
# 3. rx-result la recibe y la reenvía a rx-delivery.
#
# Para verificar el éxito, observa las ventanas de terminal de los servicios.

$OutputEncoding = [System.Text.Encoding]::UTF8

Write-Host "Performing E2E integration test: Reader -> Result -> Delivery..." -ForegroundColor Cyan

# Definir el endpoint y el cuerpo de la petición
$uri = "http://127.0.0.1:8085/reader/analyze"
$body = @{
    studyId          = "E2E-TEST-001"
    studyType        = "CT"
    anatomicalRegion = "ABDOMEN"
} | ConvertTo-Json

Write-Host "Sending POST request to $uri"

try {
    $responseRaw = Invoke-WebRequest -Uri $uri -Method Post -Body $body -ContentType "application/json"
    
    Write-Host "--- RESPONSE FROM RX-READER ---" -ForegroundColor Green
    $responseRaw.Content
    
    Write-Host "SUCCESS: Request completed. Now, check the terminal windows for logs:" -ForegroundColor Green
    Write-Host "1. 'rx-reader' window should show 'Received analysis request...'"
    Write-Host "2. 'rx-result' window should show 'Received consolidation request...'"
    Write-Host "3. 'rx-delivery' window should show 'Received enqueue request...'"
} catch {
    Write-Host "--- FAILED TO CONNECT TO RX-READER ---" -ForegroundColor Red
    Write-Host "Error: $($_.Exception.Message)"
    Write-Host "Please ensure the 'rx-reader' service is running on port 8085."
}