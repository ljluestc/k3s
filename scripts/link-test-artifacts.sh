#!/bin/bash
set -e

# Script to create symbolic links for test artifacts in all test directories
echo "Creating symbolic links for test artifacts..."

# Get the project root directory
PROJECT_ROOT=$(cd $(dirname $0)/.. && pwd)
DIST_DIR="${PROJECT_ROOT}/dist"
ARTIFACTS_DIR="${DIST_DIR}/artifacts"

# Ensure the k3s binary exists
if [ ! -f "${ARTIFACTS_DIR}/k3s" ]; then
  echo "ERROR: K3s binary not found at ${ARTIFACTS_DIR}/k3s"
  echo "Please run 'make test-prep' first"
  exit 1
fi

# Make the binary executable
chmod +x "${ARTIFACTS_DIR}/k3s"

# Create links in integration test directories
find "${PROJECT_ROOT}/tests/integration" -type d | while read TEST_DIR; do
  if [ "${TEST_DIR}" != "${PROJECT_ROOT}/tests/integration" ]; then
    mkdir -p "${TEST_DIR}/dist/artifacts"
    ln -sf "${ARTIFACTS_DIR}/k3s" "${TEST_DIR}/dist/artifacts/k3s" 2>/dev/null || true
    echo "Created link in ${TEST_DIR}"
  fi
done

# Create links at common relative paths
for i in {1..10}; do
  REL_PATH=""
  for j in $(seq 1 $i); do
    REL_PATH="${REL_PATH}../"
  done
  
  for TEST_DIR in $(find "${PROJECT_ROOT}/tests" -type d -name "integration" -o -name "e2e"); do
    TARGET="${TEST_DIR}/dist-${i}"
    mkdir -p "${TARGET}/artifacts"
    ln -sf "${ARTIFACTS_DIR}/k3s" "${TARGET}/artifacts/k3s" 2>/dev/null || true
    echo "Created relative link in ${TARGET}"
  done
done

echo "Symbolic links created successfully"
