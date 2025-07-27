// Package v1 contains mock implementations for the core v1 controllers
package v1

import (
	"k8s.io/client-go/rest"
)

// Interface for the controller factory
type Interface interface {
	// GetConfig returns the rest config
	GetConfig() *rest.Config
}

// Factory is a mock implementation of the controller factory
type Factory struct {
	config *rest.Config
}

// New creates a new factory
func New(config *rest.Config) *Factory {
	return &Factory{
		config: config,
	}
}

// GetConfig returns the rest config
func (f *Factory) GetConfig() *rest.Config {
	return f.config
}

// Core returns the core controller
func (f *Factory) Core() Interface {
	return f
}

// V1 returns the v1 controller
func (f *Factory) V1() Interface {
	return f
}

// Secret returns the secret controller
func (f *Factory) Secret() SecretController {
	return &mockSecretController{}
}

// SecretController provides access to secret resources
type SecretController interface {
	// Add functionality as needed
}

// mockSecretController is a mock implementation of SecretController
type mockSecretController struct {}
