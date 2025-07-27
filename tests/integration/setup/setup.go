func findK3sExecutable() (string, error) {
	// Check for explicitly set environment variable
	if executable := os.Getenv("K3S_EXEC"); executable != "" {
		if _, err := os.Stat(executable); err == nil {
			return executable, nil
		}
		logrus.Warnf("K3S_EXEC set to %s but file not found", executable)
	}

	// Search for the executable by traversing up the directory tree
	// to find the project root and then look in dist/artifacts
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Try to find the executable by traversing up directories
	for i := 0; i < 10; i++ { // Limit the search depth to avoid infinite loops
		// Check for dist/artifacts in the current directory
		executable := filepath.Join(currentDir, "dist", "artifacts", "k3s")
		if _, err := os.Stat(executable); err == nil {
			return executable, nil
		}

		// Check if we've reached the project root (contains go.mod)
		if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
			// Project root found, check for dist/artifacts
			executable := filepath.Join(currentDir, "dist", "artifacts", "k3s")
			if _, err := os.Stat(executable); err == nil {
				return executable, nil
			}
			
			// Try to build the binary if it doesn't exist
			buildScript := filepath.Join(currentDir, "scripts", "build-test-artifacts.sh")
			if _, err := os.Stat(buildScript); err == nil {
				logrus.Infof("Building k3s test artifacts using %s", buildScript)
				cmd := exec.Command(buildScript)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					logrus.Warnf("Failed to build test artifacts: %v", err)
				} else {
					// Check if build succeeded
					if _, err := os.Stat(executable); err == nil {
						return executable, nil
					}
				}
			}
		}
		
		// Move up one directory
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			// We've reached the root of the filesystem
			break
		}
		currentDir = parentDir
	}

	return "", fmt.Errorf("Unable to find k3s executable in dist/artifacts. Run 'make test-prep' first")
}
