// Package common provides common utilities for K3s integration tests
package common

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// FindK3sBinary locates the k3s binary for tests
func FindK3sBinary() (string, error) {
	// Check various paths for the k3s binary
	possiblePaths := []string{
		"./dist/artifacts/k3s",
		"../dist/artifacts/k3s",
		"../../dist/artifacts/k3s",
		"../../../dist/artifacts/k3s",
		"../../../../dist/artifacts/k3s",
	}

	// Add absolute path
	pwd, err := os.Getwd()
	if err == nil {
		// Find project root (look for .git directory)
		dir := pwd
		for i := 0; i < 5; i++ {
			if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
				possiblePaths = append(possiblePaths, filepath.Join(dir, "dist", "artifacts", "k3s"))
				break
			}
			dir = filepath.Dir(dir)
		}
	}

	// Check all possible paths
	for _, path := range possiblePaths {
		absPath, _ := filepath.Abs(path)
		if _, err := os.Stat(absPath); err == nil {
			logrus.Infof("Found k3s binary at %s", absPath)
			return absPath, nil
		}
	}

	return "", fmt.Errorf("unable to find k3s binary in any of the expected locations")
}

// EnsureK3sBinary ensures the k3s binary exists and is executable
func EnsureK3sBinary() (string, error) {
	// Find the k3s binary
	k3sBin, err := FindK3sBinary()
	if err != nil {
		// Try to build it
		logrus.Warn("Building k3s binary...")
		cmd := exec.Command("make", "test-prep")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("failed to build k3s binary: %v", err)
		}

		// Try again
		k3sBin, err = FindK3sBinary()
		if err != nil {
			return "", err
		}
	}

	// Make sure it's executable
	if err := os.Chmod(k3sBin, 0755); err != nil {
		return "", fmt.Errorf("failed to set executable permission on k3s binary: %v", err)
	}

	logrus.Infof("Using k3s binary at %s", k3sBin)
	return k3sBin, nil
}

// FixK3sBinaryPermissions fixes the permissions on the k3s binary
func FixK3sBinaryPermissions() error {
	k3sBin, err := FindK3sBinary()
	if err != nil {
		return err
	}

	return os.Chmod(k3sBin, 0755)
}

// SetupTestEnvironment prepares the environment for integration tests
func SetupTestEnvironment() error {
	// Ensure k3s binary exists and is executable
	_, err := EnsureK3sBinary()
	if err != nil {
		return err
	}

	// Set up any required environment variables
	os.Setenv("K3S_TEST_ENVIRONMENT", "true")
	
	return nil
}

// CreateTempDirectoryWithFiles creates a temporary directory and populates it with files
func CreateTempDirectoryWithFiles(prefix string, files map[string][]byte) (string, error) {
	tempDir, err := os.MkdirTemp("", prefix)
	if err != nil {
		return "", err
	}

	for path, content := range files {
		fullPath := filepath.Join(tempDir, path)
		dir := filepath.Dir(fullPath)
		
		if err := os.MkdirAll(dir, 0755); err != nil {
			os.RemoveAll(tempDir)
			return "", err
		}

		if err := os.WriteFile(fullPath, content, 0644); err != nil {
			os.RemoveAll(tempDir)
			return "", err
		}
	}

	return tempDir, nil
}

// RunCommand runs a command with timeout
func RunCommand(timeout time.Duration, command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	var output strings.Builder
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Start(); err != nil {
		return output.String(), err
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(timeout):
		if err := cmd.Process.Kill(); err != nil {
			logrus.Errorf("Failed to kill process: %v", err)
		}
		return output.String(), fmt.Errorf("command timed out after %s", timeout)
	case err := <-done:
		return output.String(), err
	}
}
