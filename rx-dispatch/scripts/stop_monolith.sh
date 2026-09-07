#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
PID_DIR="$ROOT_DIR/.pids"

echo "==> Stopping RX Dispatch Edge (Monolito Modular)..."

if [ -f "$PID_DIR/monolith.pid" ]; then
    pid=$(cat "$PID_DIR/monolith.pid")
    if kill -0 "$pid" 2>/dev/null; then
        echo "Stopping rx-dispatch-edge (PID $pid)..."
        kill "$pid" 2>/dev/null || true
    fi
    rm -f "$PID_DIR/monolith.pid"
fi

pkill -f "rx-dispatch-edge" 2>/dev/null || true

echo "==> RX Dispatch Edge stopped."
