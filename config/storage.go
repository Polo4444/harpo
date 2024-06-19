package config

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Polo44444/harpo/models"
	"github.com/Polo44444/harpo/storing"
)

type Storage struct {
	Type     string                `json:"type" yaml:"type"`
	Settings models.ProviderConfig `json:"settings" yaml:"settings"`
}

// Validate checks if the storage is valid
func (s *Storage) Validate(name string) error {

	typeUpper := strings.ToUpper(s.Type)

	// Get provider
	var (
		p   storing.Provider
		err error
	)
	switch typeUpper {
	case string(storing.S3Provider):
		p, err = storing.GetProvider(storing.S3Provider, s.Settings)
	default:
		return fmt.Errorf("type %s of storage %s is not valid", s.Type, name)
	}

	if err != nil {
		return fmt.Errorf("unable to get provider for storage %s: %w", name, err)
	}

	// Test provider
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // The test should not take more than 30 seconds
	defer cancel()
	err = p.Test(ctx)
	if err != nil {
		return fmt.Errorf("Storage %s test failed:\n\n%v", name, err)
	}

	return nil
}
