#!/usr/bin/env bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
BIN_DIR="$ROOT_DIR/bin"
PID_DIR="$ROOT_DIR/.pids"

mkdir -p "$PID_DIR"

HTTP_PORT="${RX_MONOLITH_HTTP_PORT:-8090}"
DICOM_PORT="${RX_MONOLITH_DICOM_PORT:-11112}"
DATA_DIR="${RX_MONOLITH_DATA_DIR:-$ROOT_DIR/data}"

mkdir -p "$DATA_DIR"

BIN_PATH="$BIN_DIR/rx-dispatch-edge"

if [ ! -f "$BIN_PATH" ]; then
    echo "==> Monolith binary not found, compiling..."
    (cd "$ROOT_DIR" && bash scripts/build_monolith.sh)
fi

echo "==> Starting RX Dispatch Edge (Monolito Modular)..."
echo "    HTTP Port:  $HTTP_PORT"
echo "    DICOM Port: $DICOM_PORT (C-STORE SCP)"
echo "    Data Dir:   $DATA_DIR"

# Stop any existing monolith instance
if [ -f "$PID_DIR/monolith.pid" ]; then
    old_pid=$(cat "$PID_DIR/monolith.pid")
    kill "$old_pid" 2>/dev/null || true
    rm -f "$PID_DIR/monolith.pid"
fi
pkill -f "rx-dispatch-edge" 2>/dev/null || true

"$BIN_PATH" --http-port "$HTTP_PORT" --dicom-port "$DICOM_PORT" --data-dir "$DATA_DIR" > "$PID_DIR/monolith.log" 2>&1 &
PID=$!
echo "$PID" > "$PID_DIR/monolith.pid"

echo "==> Process spawned with PID $PID. Waiting for health check..."
sleep 1

if curl -s --connect-timeout 2 "http://127.0.0.1:$HTTP_PORT/healthz" > /dev/null; then
    echo "==> Monolith is ONLINE and operational at http://127.0.0.1:$HTTP_PORT"
else
    echo "Warning: Monolith started (PID $PID), but health check did not respond yet. Check $PID_DIR/monolith.log"
fi
