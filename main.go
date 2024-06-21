package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Polo44444/harpo/config"
)

func main() {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	// read config file path  from flags
	var configFilePath *string
	configFilePath = flag.String("c", config.DefaultConfigPath, "path to the config file")
	flag.Parse()

	// Load config
	settings, err := config.Load(*configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// Validate config
	err = settings.Validate()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Config is valid")

	<-sigs
	// ─── Graceful Shutdown ──────────────────────────────────────────────
}
