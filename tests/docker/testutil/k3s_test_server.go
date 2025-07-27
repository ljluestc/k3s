func (s *K3sServer) start(serverArgs []string, serverName string) error {
	// Verify that k3s binary exists
	k3sBinaryPath := "../../../dist/artifacts/k3s"
	absPath, err := filepath.Abs(k3sBinaryPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for k3s binary: %v", err)
	}
	
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("k3s binary not found at %s, run 'make test-prep' first", absPath)
	}

	args := []string{
		"run", "-d", "--name", serverName,
		"--hostname", serverName,
		"--privileged",
	}

	if s.publish != "" {
		args = append(args, "-p", s.publish)
	}
