package config

import (
	"os"

	"github.com/Polo44444/harpo/models"
	"gopkg.in/yaml.v3"
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
	}

	// Validate registered storages
	for name, storage := range s.Storages {

		err := storage.Validate(name)
		if err != nil {
			return err
		}
	}

	return nil
}
