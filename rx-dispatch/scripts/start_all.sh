#!/usr/bin/env bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
BIN_DIR="$ROOT_DIR/bin"
PID_DIR="$ROOT_DIR/.pids"

mkdir -p "$PID_DIR"

echo "==> Starting RX Dispatch independent microservices mesh..."

# Environment and port defaults
export RX_SECURITY_PORT=8081
export RX_STUDY_PORT=8082
export RX_STORAGE_PORT=8083
export RX_IMAGE_PORT=8084
export RX_READER_PORT=8085
export RX_RESULT_PORT=8086
export RX_DELIVERY_PORT=8087
export RX_AUDIT_PORT=8088
export RX_GATEWAY_PORT=8089

SERVICES=(
    "security:$RX_SECURITY_PORT"
    "audit:$RX_AUDIT_PORT"
    "storage:$RX_STORAGE_PORT"
    "study:$RX_STUDY_PORT"
    "image:$RX_IMAGE_PORT"
    "reader:$RX_READER_PORT"
    "result:$RX_RESULT_PORT"
    "delivery:$RX_DELIVERY_PORT"
    "gateway:$RX_GATEWAY_PORT"
)

for entry in "${SERVICES[@]}"; do
    IFS=":" read -r name port <<< "$entry"
    bin_path="$BIN_DIR/rx-$name"
    if [ ! -f "$bin_path" ]; then
        echo "Binary $bin_path not found, compiling..."
        (cd "$ROOT_DIR" && go build -o "$bin_path" "./cmd/$name")
    fi
    echo "Starting rx-$name on port $port..."
    "$bin_path" > "$PID_DIR/$name.log" 2>&1 &
    echo $! > "$PID_DIR/$name.pid"
done

echo "==> Waiting for all 9 services to respond on health check..."
sleep 2

# Verify all services
for entry in "${SERVICES[@]}"; do
    IFS=":" read -r name port <<< "$entry"
    curl -s --connect-timeout 2 "http://127.0.0.1:$port/healthz" || echo "Warning: rx-$name not ready yet"
    echo ""
done

echo "==> Mesh started successfully. Topology ready at http://127.0.0.1:8089/gateway/topology"
