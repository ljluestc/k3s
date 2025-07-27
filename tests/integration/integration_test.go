// Package integration contains common setup code for integration tests
package integration

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/k3s-io/k3s/tests/integration/common"
	"github.com/sirupsen/logrus"
)

var (
	k3sBin    string
	setupDone bool
)

// setupTestEnvironment sets up the test environment for all integration tests
func setupTestEnvironment() error {
	if setupDone {
		return nil
	}

	// Find and verify k3s binary
	var err error
	k3sBin, err = common.EnsureK3sBinary()
	if err != nil {
		return err
	}

	// Create symlinks to ensure tests at different directory depths can find the binary
	pwd, _ := os.Getwd()
	projRoot := filepath.Dir(filepath.Dir(pwd))
	distDir := filepath.Join(projRoot, "dist")
	artifactsDir := filepath.Join(distDir, "artifacts")

	// Ensure directories exist
	os.MkdirAll(distDir, 0755)
	os.MkdirAll(artifactsDir, 0755)

	// Create symlinks if needed
	if _, err := os.Stat(filepath.Join(artifactsDir, "k3s")); os.IsNotExist(err) {
		if err := os.Symlink(k3sBin, filepath.Join(artifactsDir, "k3s")); err != nil {
			logrus.Warnf("Failed to create symlink: %v", err)
		}
	}

	setupDone = true
	return nil
}

// TestMain is the main entry point for integration tests
func TestMain(m *testing.M) {
	flag.Parse()

	if err := setupTestEnvironment(); err != nil {
		logrus.Fatalf("Failed to set up test environment: %v", err)
	}

	os.Exit(m.Run())
}
