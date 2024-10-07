#!/usr/bin/env bash
set -euo pipefail

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/run_migrations ./cmd/migrations/main.go
echo "echoes migrations build completed"

