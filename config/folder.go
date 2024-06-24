package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/Polo44444/harpo/archiving"
)

type Folder struct {
	Name                string `json:"name" yaml:"name"`
	Path                string `json:"path" yaml:"path"`
	Remove              bool   `json:"remove" yaml:"remove"`
	IgnoreArchiveErrors bool   `json:"ignore_archive_errors" yaml:"ignore_archive_errors"`
	Destination         string `json:"destination" yaml:"destination"`
	Schedule            string `json:"schedule" yaml:"schedule"`
	Archiver            string `json:"archiver" yaml:"archiver"`
	// TODO: Give the ability to add compression level
	Storages  []string `json:"storages" yaml:"storages"`
	Notifiers []string `json:"notifiers" yaml:"notifiers"`
}

// Validate checks if the folder is valid
func (f *Folder) Validate(storages map[string]Storage, notifiers map[string]Notifier) error {

	// Check name
	if strings.TrimSpace(f.Name) == "" {
		return fmt.Errorf("name of path %s is not valid", f.Path)
	}

	// Check path
	fileInfo, err := os.Stat(f.Path)
	if err != nil {
		return fmt.Errorf("unable to check folder path %s existence: %w", f.Path, err)
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("path %s is not a folder", f.Path)
	}

	// Check destination
	if strings.TrimSpace(f.Destination) == "" {
		return fmt.Errorf("destination of path %s is not valid", f.Path)
	}

	// Check schedule
	if strings.TrimSpace(f.Schedule) == "" {
		return fmt.Errorf("schedule of path %s is not valid", f.Path)
	}

	// Check archiver
	if strings.TrimSpace(f.Archiver) == "" ||
		(strings.ToUpper(f.Archiver) != string(archiving.ZipProvider) &&
			strings.ToUpper(f.Archiver) != string(archiving.TarProvider)) {
		return fmt.Errorf("archiver of path %s is not valid", f.Path)
	}

	// Check storages
	for _, storage := range f.Storages {
		if _, ok := storages[storage]; !ok {
			return fmt.Errorf("storage %s of path %s is not valid or have not been declared", storage, f.Path)
		}
	}

	// Check notifiers
	for _, notifier := range f.Notifiers {
		if _, ok := notifiers[notifier]; !ok {
			return fmt.Errorf("notifier %s of path %s is not valid or have not been declared", notifier, f.Path)
		}
	}

	return nil
}
