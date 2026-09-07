@echo off
echo ==========================================================
echo Stopping RX Dispatch independent microservices mesh...
echo ==========================================================

for %%s in (gateway security study storage image reader result delivery audit) do (
    taskkill /F /IM rx-%%s.exe /T > nul 2>&1
)

echo All RX Dispatch processes stopped.