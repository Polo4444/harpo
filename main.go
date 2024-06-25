package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Polo44444/harpo/backup"
	"github.com/Polo44444/harpo/config"
	"github.com/Polo44444/harpo/constants"
	"github.com/Polo44444/harpo/utils"
)

func main() {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	// Read config file path  from flags
	var configFilePath *string = flag.String("c", config.DefaultConfigPath, "path to the config file")
	flag.Parse()

	// Load config
	settings, err := config.Load(*configFilePath)
	utils.LogFatalIfErr(err)

	// Validate config
	utils.LogFatalIfErr(settings.Validate())

	log.Println("Config is valid")

	// ─── Load Providers ──────────────────────────────────────────────────
	storages := settings.GetStorageProviders()
	notifiers := settings.GetNotifierProviders()

	// Start backup engine
	bck := backup.NewEngine(settings.Folders, storages, notifiers)
	utils.LogFatalIfErr(bck.BuildJobs())
	bck.Start()

	log.Printf("%s backup engine started\n", constants.AppName)

	<-sigs
	// ─── Graceful Shutdown ──────────────────────────────────────────────

	// Stop providers
	for _, storage := range storages {
		storage.Close(context.Background())
	}

	for _, notifier := range notifiers {
		notifier.Close(context.Background())
	}

	// Stop backup engine
	bck.Stop()
}
