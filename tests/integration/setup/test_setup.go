package setup

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
)

// SetupTempDirectory creates a temporary directory for tests
func SetupTempDirectory(t *testing.T, prefix string) (string, func()) {
	t.Helper()
	
	tempDir, err := os.MkdirTemp("", prefix)
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	
	cleanup := func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to clean up temp directory: %v", err)
		}
	}
	
	return tempDir, cleanup
}

// SetupTestArtifacts creates a directory structure with links to the k3s binary
func SetupTestArtifacts(t *testing.T, testDir string) error {
	t.Helper()
	
	// Find k3s binary
	k3sBin, err := FindK3sBinary()
	if err != nil {
		return fmt.Errorf("failed to find k3s binary: %v", err)
	}
	
	// Create artifacts directory structure
	artifactsDir := filepath.Join(testDir, "dist", "artifacts")
	if err := os.MkdirAll(artifactsDir, 0755); err != nil {
		return fmt.Errorf("failed to create artifacts directory: %v", err)
	}
	
	// Create link to k3s binary
	k3sLink := filepath.Join(artifactsDir, "k3s")
	if err := os.Symlink(k3sBin, k3sLink); err != nil {
		return fmt.Errorf("failed to create k3s symlink: %v", err)
	}
	
	logrus.Infof("Test artifacts setup at %s", testDir)
	return nil
}

// SetupTestEnvironmentForTest sets up the test environment for a specific test
func SetupTestEnvironmentForTest(t *testing.T) (string, func()) {
	t.Helper()
	
	// Create temp directory
	tempDir, cleanup := SetupTempDirectory(t, "k3s-test-")
	
	// Setup test artifacts
	if err := SetupTestArtifacts(t, tempDir); err != nil {
		cleanup()
		t.Fatalf("Failed to setup test artifacts: %v", err)
	}
	
	return tempDir, cleanup
}
