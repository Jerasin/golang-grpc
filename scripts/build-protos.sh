#!/bin/bash

# ทำให้ script รันจาก root เสมอ
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
cd "$ROOT_DIR"

PROTO_DIR=proto
OUT_DIR=.

mkdir -p $OUT_DIR

protoc -I $PROTO_DIR \
  --go_out $OUT_DIR \
  --go-grpc_out $OUT_DIR \
  --grpc-gateway_out $OUT_DIR \
  $PROTO_DIR/service.proto
