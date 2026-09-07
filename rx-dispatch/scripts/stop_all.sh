#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
PID_DIR="$ROOT_DIR/.pids"

echo "==> Stopping RX Dispatch independent microservices mesh..."

if [ -d "$PID_DIR" ]; then
    for pid_file in "$PID_DIR"/*.pid; do
        if [ -f "$pid_file" ]; then
            name=$(basename "$pid_file" .pid)
            pid=$(cat "$pid_file")
            if kill -0 "$pid" 2>/dev/null; then
                echo "Stopping rx-$name (PID $pid)..."
                kill "$pid" 2>/dev/null || true
            fi
            rm -f "$pid_file"
        fi
    done
fi

# Also kill any orphan processes if needed
pkill -f "rx-gateway" || true
pkill -f "rx-security" || true
pkill -f "rx-study" || true
pkill -f "rx-storage" || true
pkill -f "rx-image" || true
pkill -f "rx-reader" || true
pkill -f "rx-result" || true
pkill -f "rx-delivery" || true
pkill -f "rx-audit" || true

echo "==> All RX Dispatch processes stopped."
