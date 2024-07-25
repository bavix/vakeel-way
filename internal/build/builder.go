package build

import (
	"github.com/bavix/vakeel-way/internal/config"
	"github.com/bavix/vakeel-way/internal/domain/usecases"
)

// Builder is a struct that holds the configuration for building the application.
// It is used to create a new instance of the application.
type Builder struct {
	config config.Config

	checker *usecases.Checker
}

// NewBuilder creates a new instance of the Builder struct.
//
// It reads the configuration from the environment variables and creates a new
// instance of the Builder struct with the configuration.
// If there is an error reading the configuration, it returns the error.
//
// Returns a pointer to the newly created Builder instance and an error if there
// was an error reading the configuration.
//
//nolint:exhaustruct
func NewBuilder(config config.Config) (*Builder, error) {
	// Create a new instance of the Builder struct with the configuration.
	return &Builder{config: config}, nil
}
