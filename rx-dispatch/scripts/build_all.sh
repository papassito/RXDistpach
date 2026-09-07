#!/usr/bin/env bash
set -e

echo "=========================================================="
echo "RX DISPATCH BY KLIK - Independent Microservices Builder"
echo "Architecture: 100% Go, Decoupled, Zero-Monolith"
echo "=========================================================="

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
BIN_DIR="$ROOT_DIR/bin"

mkdir -p "$BIN_DIR"

SERVICES=("security" "study" "storage" "image" "reader" "result" "delivery" "audit" "gateway")

for svc in "${SERVICES[@]}"; do
    echo "==> Compiling independent binary: rx-$svc from ./cmd/$svc"
    (cd "$ROOT_DIR" && go build -o "$BIN_DIR/rx-$svc" "./cmd/$svc")
done

echo ""
echo "All 9 independent Go binaries compiled successfully into $BIN_DIR:"
ls -lh "$BIN_DIR"
