package backup

import (
	"github.com/Polo44444/harpo/alerting"
	"github.com/Polo44444/harpo/config"
	"github.com/Polo44444/harpo/storing"
)

type CtxString string

type Engine struct {
	folders   []config.Folder
	storages  map[string]storing.Provider
	notifiers map[string]alerting.Provider
}

// Processes contexts keys
const (
	// Holds the archive *os.File pointer coming from archive process.
	// Please make sure you close the file after you are done with it.
	ArchiveCtxKey CtxString = "archive"
)

// NewEngine creates a new Engine instance
func NewEngine(folders []config.Folder, storages map[string]storing.Provider, notifiers map[string]alerting.Provider) *Engine {
	return &Engine{
		folders:   folders,
		storages:  storages,
		notifiers: notifiers,
	}
}

// Run starts the backup process
func (e *Engine) Run() error {

	return nil
}

// Stop stops the backup process and closes the providers. Make sure to create new providers if you want to start the backup process again.
func (e *Engine) Stop() error {
	return nil
}
