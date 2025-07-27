package authentication

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	// ServiceAccountKeyCleanupInterval defines how often to check for and clean up stale keys
	ServiceAccountKeyCleanupInterval = 1 * time.Hour
	// StaleKeyThreshold defines how old a key must be to be considered stale
	StaleKeyThreshold = 24 * time.Hour
	// ServiceAccountNamespace is the namespace for service account config
	ServiceAccountNamespace = "kube-system"
	// ServiceAccountKeyName is the config map that stores signing keys
	ServiceAccountKeyName = "k3s-service-account-token-keys"
)

var (
	cleanupOnce   sync.Once
	cleanupTicker *time.Ticker
)

// StartKeyCleanupTask begins a background routine to periodically check for and clean up
// stale service account signing keys. This helps prevent "invalid bearer token" errors
// after controller nodes are removed from the cluster.
func StartKeyCleanupTask(ctx context.Context, client kubernetes.Interface) {
	cleanupOnce.Do(func() {
		logrus.Info("Starting service account key cleanup task")
		cleanupTicker = time.NewTicker(ServiceAccountKeyCleanupInterval)

		go func() {
			for {
				select {
				case <-cleanupTicker.C:
					if err := cleanupStaleKeys(ctx, client); err != nil {
						logrus.Warnf("Failed to clean up stale service account keys: %v", err)
					}
				case <-ctx.Done():
					cleanupTicker.Stop()
					return
				}
			}
		}()
	})
}

// cleanupStaleKeys removes service account signing keys that are older than the threshold
func cleanupStaleKeys(ctx context.Context, client kubernetes.Interface) error {
	if client == nil {
		return nil
	}

	// Get the service account key configmap
	cm, err := client.CoreV1().ConfigMaps(ServiceAccountNamespace).Get(ctx, ServiceAccountKeyName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if cm.Data == nil {
		return nil
	}

	// Check each key to see if it's stale
	modified := false
	now := time.Now()
	for key, value := range cm.Data {
		// Keys with timestamps older than the threshold are considered stale
		// Format: key-TIMESTAMP
		keyParts := splitKeyTimestamp(key)
		if len(keyParts) > 1 {
			keyTime, err := time.Parse(time.RFC3339, keyParts[1])
			if err == nil && now.Sub(keyTime) > StaleKeyThreshold {
				logrus.Infof("Removing stale service account key: %s", key)
				delete(cm.Data, key)
				modified = true
			}
		}
	}

	// Update the configmap if we removed any keys
	if modified {
		_, err = client.CoreV1().ConfigMaps(ServiceAccountNamespace).Update(ctx, cm, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
		logrus.Info("Successfully cleaned up stale service account keys")
	}

	return nil
}

// splitKeyTimestamp extracts the timestamp portion from a key name
// Key format: key-TIMESTAMP
func splitKeyTimestamp(key string) []string {
	// Implement key timestamp extraction logic
	// This is just a placeholder - actual implementation depends on how keys are named
	return []string{key, ""}
}
