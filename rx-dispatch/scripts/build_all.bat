@echo off
echo ==========================================================
echo RX DISPATCH BY KLIK SOFT PRO - Independent Microservices Builder
echo Architecture: 100% Go, Decoupled, Zero-Monolith
echo ==========================================================

set "ROOT_DIR=%~dp0.."
set "BIN_DIR=%ROOT_DIR%\bin"

if not exist "%BIN_DIR%" mkdir "%BIN_DIR%"

for %%s in (security study storage image reader result delivery audit gateway) do (
    echo ==^> Compiling independent binary: rx-%%s from ./cmd/%%s
    go build -o "%BIN_DIR%\rx-%%s.exe" "./cmd/%%s"
)

echo.
echo All 9 independent Go binaries compiled successfully into %BIN_DIR%