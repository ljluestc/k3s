#!/bin/bash
set -e

# This script runs docker-based integration tests with proper environment setup
# Usage: ./run-test.sh [test package path]

# Get the project root directory
PROJECT_ROOT=$(cd $(dirname $0)/../.. && pwd)
DIST_DIR="${PROJECT_ROOT}/dist"
ARTIFACTS_DIR="${DIST_DIR}/artifacts"

# Ensure k3s binary exists
if [ ! -f "${ARTIFACTS_DIR}/k3s" ]; then
  echo "K3s binary not found at ${ARTIFACTS_DIR}/k3s"
  echo "Building k3s binary..."
  "${PROJECT_ROOT}/scripts/build-test-artifacts.sh"
fi

# Run specific test if provided, otherwise run all docker tests
if [ -n "$1" ]; then
  TEST_PACKAGE="$1"
  echo "Running test package: ${TEST_PACKAGE}"
  go test -v "${TEST_PACKAGE}"
else
  echo "Running all docker tests"
  go test -v ./tests/docker/...
fi
