#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
BIN_DIR="$ROOT_DIR/bin"

mkdir -p "$BIN_DIR"
cd "$ROOT_DIR"

echo "==> Compiling RX Dispatch Edge (Monolito Modular Unificado)..."
go build -ldflags="-s -w" -o "$BIN_DIR/rx-dispatch-edge" ./cmd/monolith

echo "==> Compilation successful!"
ls -lh "$BIN_DIR/rx-dispatch-edge"
