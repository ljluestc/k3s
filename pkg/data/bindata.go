// Package data provides access to embedded binary assets
package data

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// Asset returns the content of the named asset
func Asset(name string) []byte {
	// In development mode, try to read from filesystem
	if path := findAssetPath(name); path != "" {
		content, err := os.ReadFile(path)
		if err == nil {
			return content
		}
		logrus.Warnf("Failed to read asset %s from %s: %v", name, path, err)
	}
	
	// Return empty data for now
	// In production, this would be replaced with embedded data
	return []byte{}
}

// Name returns the canonical name of the asset
func Name(name string) string {
	return name
}

// ReadFile is an alias for Asset for backwards compatibility
func ReadFile(name string) []byte {
	return Asset(name)
}

// findAssetPath tries to locate the asset in the local filesystem
// This is primarily for development mode
func findAssetPath(name string) string {
	// Try common locations for development assets
	locations := []string{
		"./build/data",
		"./dist/artifacts",
		"./dist",
		".",
	}
	
	for _, location := range locations {
		path := filepath.Join(location, name)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	
	return ""
}
