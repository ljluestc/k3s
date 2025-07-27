#!/bin/bash
set -e

# Comprehensive script to run all K3s tests with proper setup
echo "===== K3s Test Runner ====="

# Get the project root directory
PROJECT_ROOT=$(cd $(dirname $0)/.. && pwd)
DIST_DIR="${PROJECT_ROOT}/dist"
ARTIFACTS_DIR="${DIST_DIR}/artifacts"

# Step 1: Prepare test artifacts
echo "Step 1: Preparing test artifacts..."
"${PROJECT_ROOT}/scripts/build-test-artifacts.sh"

# Verify the k3s binary exists and is executable
if [ ! -f "${ARTIFACTS_DIR}/k3s" ]; then
  echo "ERROR: K3s binary not found at ${ARTIFACTS_DIR}/k3s"
  exit 1
fi

chmod +x "${ARTIFACTS_DIR}/k3s"
echo "K3s binary verified at ${ARTIFACTS_DIR}/k3s"

# Create additional symlinks in the integration test directories
for TEST_DIR in $(find "${PROJECT_ROOT}/tests/integration" -type d); do
  if [ "${TEST_DIR}" != "${PROJECT_ROOT}/tests/integration" ]; then
    mkdir -p "${TEST_DIR}/dist/artifacts"
    ln -sf "${ARTIFACTS_DIR}/k3s" "${TEST_DIR}/dist/artifacts/k3s" 2>/dev/null || true
  fi
done

# Step 2: Run the unit tests
echo "Step 2: Running unit tests..."
cd "${PROJECT_ROOT}"
go test -v ./pkg/... ./cmd/... 2>&1 | tee unit-tests.log

# Step 3: Run integration tests
echo "Step 3: Running integration tests..."
cd "${PROJECT_ROOT}"
go test -v ./tests/integration/... 2>&1 | tee integration-tests.log

# Step 4: Run E2E tests if requested
if [ "$1" == "--e2e" ]; then
  echo "Step 4: Running E2E tests..."
  cd "${PROJECT_ROOT}"
  go test -v ./tests/e2e/... 2>&1 | tee e2e-tests.log
else
  echo "Skipping E2E tests (use --e2e to run them)"
fi

# Step 5: Run Docker tests if requested
if [ "$1" == "--docker" ] || [ "$2" == "--docker" ]; then
  echo "Step 5: Running Docker tests..."
  cd "${PROJECT_ROOT}"
  "${PROJECT_ROOT}/tests/docker/run-test.sh" 2>&1 | tee docker-tests.log
else
  echo "Skipping Docker tests (use --docker to run them)"
fi

# Summary
echo "===== Test Execution Complete ====="
echo "Log files:"
echo "- Unit Tests: ${PROJECT_ROOT}/unit-tests.log"
echo "- Integration Tests: ${PROJECT_ROOT}/integration-tests.log"
if [ "$1" == "--e2e" ]; then
  echo "- E2E Tests: ${PROJECT_ROOT}/e2e-tests.log"
fi
if [ "$1" == "--docker" ] || [ "$2" == "--docker" ]; then
  echo "- Docker Tests: ${PROJECT_ROOT}/docker-tests.log"
fi
