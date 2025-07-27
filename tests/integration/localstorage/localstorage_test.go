func testSetup(t *testing.T, node0 string) (cleanup func(), err error) {
	t.Helper()
	testLock.Lock()

	executable, err := common.FindK3sBinary()
	if err != nil {
		t.Logf("Attempting to find k3s binary using multiple methods")
		// Fallback to local method
		executable, err = findK3sExecutable()
		if err != nil {
			return nil, err
		}
	}
	tempDir, err := os.MkdirTemp("", "k3s-local-storage-*")
	if err != nil {
		return nil, err
	}
	var wait sync.WaitGroup
