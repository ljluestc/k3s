package control

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

var (
	removedNodes    = make(map[string]string)
	removedNodesRWL sync.RWMutex
	keyCleanupOnce  sync.Once
)

// RegisterRemovedNode registers a node as removed for token cleanup
func RegisterRemovedNode(nodeName string, nodeID string) {
	removedNodesRWL.Lock()
	defer removedNodesRWL.Unlock()
	removedNodes[nodeName] = nodeID
}

// IsNodeRemoved checks if a node is registered as removed
func IsNodeRemoved(nodeName string) (string, bool) {
	removedNodesRWL.RLock()
	defer removedNodesRWL.RUnlock()
	nodeID, ok := removedNodes[nodeName]
	return nodeID, ok
}

// EnsureKeyCleanupStarted starts the key cleanup task once
func EnsureKeyCleanupStarted(ctx context.Context, clientset kubernetes.Interface) {
	keyCleanupOnce.Do(func() {
		go runKeyCleanup(ctx, clientset)
	})
}

// runKeyCleanup periodically cleans up keys for removed nodes
func runKeyCleanup(ctx context.Context, clientset kubernetes.Interface) {
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			cleanupNodeKeys(ctx, clientset)
		}
	}
}

// cleanupNodeKeys removes authentication keys for nodes that have been removed from the cluster
func cleanupNodeKeys(ctx context.Context, clientset kubernetes.Interface) {
	removedNodesRWL.RLock()
	nodesToClean := make(map[string]string)
	for name, id := range removedNodes {
		nodesToClean[name] = id
	}
	removedNodesRWL.RUnlock()

	if len(nodesToClean) == 0 {
		return
	}

	logrus.Infof("Running cleanup for %d removed nodes", len(nodesToClean))
	
	// Clean up node-specific secrets
	for nodeName, nodeID := range nodesToClean {
		// Here we would clean up any authentication tokens or secrets related to the node
		logrus.Infof("Cleaning up authentication data for removed node %s (ID: %s)", nodeName, nodeID)
		
		// Example: delete node password secret
		err := clientset.CoreV1().Secrets("kube-system").Delete(ctx, "node-password-"+nodeName, metav1.DeleteOptions{})
		if err != nil {
			logrus.Warnf("Failed to delete node password for %s: %v", nodeName, err)
		}
	}
}
