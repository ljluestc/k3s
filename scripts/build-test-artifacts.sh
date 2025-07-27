#!/bin/bash
set -e

# Script to build K3s artifacts needed for running integration tests
echo "Building k3s artifacts for integration tests..."

# Get the project root directory
PROJECT_ROOT=$(cd $(dirname $0)/.. && pwd)
DIST_DIR="${PROJECT_ROOT}/dist"
ARTIFACTS_DIR="${DIST_DIR}/artifacts"

# Create necessary directories
mkdir -p "${ARTIFACTS_DIR}"

# Build k3s binary
echo "Building k3s binary..."
cd "${PROJECT_ROOT}"
go build -o "${ARTIFACTS_DIR}/k3s" .

echo "K3s binary built successfully at ${ARTIFACTS_DIR}/k3s"
echo "Test artifacts prepared successfully"

# Set execute permissions
chmod +x "${ARTIFACTS_DIR}/k3s"

# Verify the binary exists
if [ -f "${ARTIFACTS_DIR}/k3s" ]; then
  echo "✅ K3s binary verified at ${ARTIFACTS_DIR}/k3s"
else
  echo "❌ Failed to build k3s binary"
  exit 1
fi
