package setup

import (
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	setupOnce sync.Once
	setupErr  error
)

// Initialize sets up the test environment once
func Initialize() error {
	setupOnce.Do(func() {
		setupErr = SetupTestEnvironment()
		if setupErr != nil {
			logrus.Errorf("Failed to set up test environment: %v", setupErr)
		} else {
			logrus.Info("Test environment set up successfully")
		}
	})
	return setupErr
}

func init() {
	// Initialize the test environment when the package is imported
	if err := Initialize(); err != nil {
		logrus.Warnf("Test environment initialization failed: %v", err)
	}
}
