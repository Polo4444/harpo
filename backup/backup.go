package backup

import (
	"context"
	"fmt"
	"time"

	"github.com/Polo44444/harpo/alerting"
	"github.com/Polo44444/harpo/config"
	"github.com/Polo44444/harpo/storing"
	"github.com/go-co-op/gocron/v2"
)

type CtxString string

type Engine struct {
	ctx       context.Context
	cancel    context.CancelFunc
	sch       gocron.Scheduler
	folders   []config.Folder
	storages  map[string]storing.Provider  // TODO: Should be changed to sync.Map
	notifiers map[string]alerting.Provider // TODO: Should be changed to sync.Map
}

// Processes contexts keys
const (
	// Holds the archive *os.File pointer coming from archive process.
	// Please make sure you close the file after you are done with it.
	ArchiveCtxKey CtxString = "archive"
)

const (
	harpoTag             = "harpo"
	harpoBackupTag       = "harpo:backup"
	harpoBackupFolderTag = "harpo:backup:%s"
)

var (
	ProcessTimeout = time.Hour * 2 // TODO: Calculate the timeout based on the folder size
)

// NewEngine creates a new Engine instance
func NewEngine(folders []config.Folder, storages map[string]storing.Provider, notifiers map[string]alerting.Provider) *Engine {

	ctx, cancel := context.WithCancel(context.TODO())

	return &Engine{
		folders:   folders,
		storages:  storages,
		notifiers: notifiers,
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Build Jobs build a list of jobs for the given folders
func (e *Engine) BuildJobs() error {

	var err error

	// Create scheduler
	e.sch, err = gocron.NewScheduler()
	if err != nil {
		return fmt.Errorf("unable to create scheduler: %w", err)
	}

	// Loop trough folders to create cron jobs
	for _, folder := range e.folders {

		// Build folder storages
		folderStorages := map[string]storing.Provider{}
		for _, storageName := range folder.Storages {
			if storage, ok := e.storages[storageName]; ok {
				folderStorages[storageName] = storage
			}
		}

		// Build folder notifiers
		folderNotifiers := map[string]alerting.Provider{}
		for _, notifierName := range folder.Notifiers {
			if notifier, ok := e.notifiers[notifierName]; ok {
				folderNotifiers[notifierName] = notifier
			}
		}

		// Create cron job
		_, err := e.sch.NewJob(
			gocron.CronJob(folder.Schedule, false),
			gocron.NewTask(func(folder config.Folder, stors map[string]storing.Provider, notifs map[string]alerting.Provider) {
				e.ProcessFolder(folder, stors, notifs)
			},
				folder,
				folderStorages,
				folderNotifiers,
			),
			gocron.WithTags(harpoTag, harpoBackupTag, fmt.Sprintf(harpoBackupFolderTag, folder.Name)),
		)
		if err != nil {
			e.RemoveJobs()
			return fmt.Errorf("unable to create cron job of folder %s: %w", folder.Name, err)
		}
	}

	return nil
}

// start starts the backup process in background and returns immediately.
func (e *Engine) Start() {

	if len(e.sch.Jobs()) == 0 {
		return
	}

	e.sch.Start()
}

// Run starts the backup process and blocks until the process is finished.
func (e *Engine) Run() {

	e.Start()
	<-e.ctx.Done() // Wait until the context is done
}

// RemoveJobs removes all scheduled jobs
func (e *Engine) RemoveJobs() {
	e.sch.RemoveByTags(harpoTag)
}

// RemoveJob removes the scheduled job of the given folder
func (e *Engine) RemoveJob(folderName string) {
	e.sch.RemoveByTags(fmt.Sprintf(harpoBackupFolderTag, folderName))
}

// StopJobs stops all scheduled jobs
func (e *Engine) StopJobs() error {
	return e.sch.StopJobs()
}

// ProcessFolder processes the given folder. It archives and uploads the folder to the storages.
func (e *Engine) ProcessFolder(folder config.Folder, storages map[string]storing.Provider, notifiers map[string]alerting.Provider) error {

	// Setup chain
	chain := NewArchiver()
	chain.
		setNext(NewUploader()).
		setNext(nil)

	// Execute chain
	ctx, cancel := context.WithTimeout(e.ctx, ProcessTimeout)
	defer cancel()
	chain.process(ctx, folder, storages, notifiers)

	return nil
}

// Stop cancels the engine context, shutdowns the engine and returns immediately.
func (e *Engine) Stop() error {

	e.cancel()
	return e.sch.Shutdown()
}
