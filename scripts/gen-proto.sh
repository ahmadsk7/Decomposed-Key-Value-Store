#!/bin/bash

# Generate Protocol Buffers code for Go services

set -e  # Exit on any error

echo "Generating Protocol Buffers code..."

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "Error: protoc is not installed"
    echo "Please install Protocol Buffers compiler:"
    echo "  brew install protobuf  # macOS"
    echo "  apt-get install protobuf-compiler  # Ubuntu/Debian"
    exit 1
fi

# Directory where proto files live
PROTO_DIR="proto/kv/v1"
PROTO_FILE="${PROTO_DIR}/kv.proto"

# Check if proto file exists
if [ ! -f "$PROTO_FILE" ]; then
    echo "Error: $PROTO_FILE not found"
    exit 1
fi

# Create output directories
mkdir -p services/kv-store/internal/proto/kv/v1
mkdir -p services/api-gateway/internal/proto/kv/v1

echo "Compiling $PROTO_FILE..."

# Generate Go code using protoc
# --proto_path=.: Sets the root directory for proto imports
# --go_out=. and --go-grpc_out=.: Where to output generated code
protoc \
    --proto_path=. \
    --go_out=. \
    --go-grpc_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_opt=paths=source_relative \
    "$PROTO_FILE"

echo ""
echo "Generated files:"
echo "  - proto/kv/v1/kv.pb.go (messages)"
echo "  - proto/kv/v1/kv_grpc.pb.go (gRPC server/client)"
echo ""
echo "Copying to services..."

# Copy generated files to both services
cp proto/kv/v1/kv.pb.go services/kv-store/internal/proto/kv/v1/
cp proto/kv/v1/kv_grpc.pb.go services/kv-store/internal/proto/kv/v1/

cp proto/kv/v1/kv.pb.go services/api-gateway/internal/proto/kv/v1/
cp proto/kv/v1/kv_grpc.pb.go services/api-gateway/internal/proto/kv/v1/

echo "✓ Proto generation complete!"
echo "✓ Both services now have the gRPC code"

