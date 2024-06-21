package config

import (
	"os"

	"github.com/Polo44444/harpo/alerting"
	"github.com/Polo44444/harpo/models"
	"github.com/Polo44444/harpo/storing"
	"gopkg.in/yaml.v3"
)

const (
	// DefaultConfigPath is the default path to the config file
	DefaultConfigPath = "harpo.yml"
)

type Notifier struct {
	Type     string                `json:"type" yaml:"type"`
	Settings models.ProviderConfig `json:"settings" yaml:"settings"`
}

type Settings struct {
	Folders   []Folder            `json:"folders" yaml:"folders"`
	Storages  map[string]Storage  `json:"storages" yaml:"storages"`
	Notifiers map[string]Notifier `json:"notifiers" yaml:"notifiers"`
}

func Load(path string) (*Settings, error) {

	// Open file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Load config from file
	settings := &Settings{}
	err = yaml.NewDecoder(file).Decode(&settings)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

// Validate checks if the settings are valid
func (s *Settings) Validate() error {

	// Validate folders
	for _, folder := range s.Folders {

		err := folder.Validate(s.Storages, s.Notifiers)
		if err != nil {
			return err
		}

		// Validate used storages
		for _, fStorage := range folder.Storages {

			storage, ok := s.Storages[fStorage]
			if !ok {
				return err
			}

			err = storage.Validate(fStorage, folder.Destination)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// GetStorageProviders returns the storage providers
func (s *Settings) GetStorageProviders() map[string]storing.Provider {

	providers := map[string]storing.Provider{}

	for name, storage := range s.Storages {

		var (
			p   storing.Provider
			err error
		)
		switch storage.Type {
		case string(storing.S3Provider):

			p, err = storing.GetProvider(storing.S3Provider, storage.Settings)
			if err != nil {
				continue
			}
		default:
			continue
		}

		providers[name] = p
	}

	return providers
}

// GetNotifierProviders returns the notifier providers
func (s *Settings) GetNotifierProviders() map[string]alerting.Provider {

	providers := map[string]alerting.Provider{}

	for name, notifier := range s.Notifiers {

		var (
			p   alerting.Provider
			err error
		)
		switch notifier.Type {
		case string(alerting.DiscordProvider):

			p, err = alerting.GetProvider(alerting.DiscordProvider, notifier.Settings)
			if err != nil {
				continue
			}
		default:
			continue
		}

		providers[name] = p
	}

	return providers
}
