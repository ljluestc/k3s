// Package integration provides utilities for K3s integration tests
package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// FindK3sExecutable finds the k3s executable by checking various possible locations
func FindK3sExecutable() (string, error) {
	// Check relative paths first
	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)
	
	// Try direct path first
	k3sPath := filepath.Join(currentDir, "..", "..", "dist", "artifacts", "k3s")
	if _, err := os.Stat(k3sPath); err == nil {
		return k3sPath, nil
	}
	
	// Try absolute path
	projectRoot, err := filepath.Abs(filepath.Join(currentDir, "..", ".."))
	if err != nil {
		return "", err
	}
	
	k3sPath = filepath.Join(projectRoot, "dist", "artifacts", "k3s")
	if _, err := os.Stat(k3sPath); err == nil {
		return k3sPath, nil
	}
	
	// Try various relative paths that might be used by different tests
	for i := 1; i <= 10; i++ {
		relPath := strings.Repeat("../", i) + "dist/artifacts/k3s"
		absPath, err := filepath.Abs(filepath.Join(currentDir, relPath))
		if err != nil {
			continue
		}
		
		if _, err := os.Stat(absPath); err == nil {
			return absPath, nil
		}
	}
	
	return "", fmt.Errorf("k3s executable not found in any of the expected locations")
}

// EnsureK3sExecutable ensures the k3s executable exists and is accessible
func EnsureK3sExecutable() (string, error) {
	k3sPath, err := FindK3sExecutable()
	if err != nil {
		// Try to build it
		logrus.Warn("K3s executable not found, attempting to build it...")
		cmd := exec.Command("make", "test-prep")
		cmd.Dir = filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(k3sPath))))))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("failed to build k3s: %v", err)
		}
		
		// Try finding it again
		k3sPath, err = FindK3sExecutable()
		if err != nil {
			return "", err
		}
	}
	
	// Make sure it's executable
	if err := os.Chmod(k3sPath, 0755); err != nil {
		return "", fmt.Errorf("failed to make k3s executable: %v", err)
	}
	
	return k3sPath, nil
}
