@echo off
echo ==========================================================
echo Starting RX Dispatch independent microservices mesh...
echo ==========================================================

set "BIN_DIR=%~dp0..\bin"

set RX_SECURITY_PORT=8081
set RX_STUDY_PORT=8082
set RX_STORAGE_PORT=8083
set RX_IMAGE_PORT=8084
set RX_READER_PORT=8085
set RX_RESULT_PORT=8086
set RX_DELIVERY_PORT=8087
set RX_AUDIT_PORT=8088
set RX_GATEWAY_PORT=8089

for %%s in (security audit storage study image reader result delivery gateway) do (
    echo Starting rx-%%s...
    start "rx-%%s" /B "%BIN_DIR%\rx-%%s.exe"
)

echo.
echo All services started in the background.