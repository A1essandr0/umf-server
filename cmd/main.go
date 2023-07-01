package main

import (
	"log"

	"github.com/A1essandr0/umf-server/internal/config"
	"github.com/A1essandr0/umf-server/internal/repositories"
	"github.com/A1essandr0/umf-server/internal/server"
)

func main() {
	config := config.Init("config")

	kvStore, err := repositories.NewKVStore(config)
	if err != nil {
		log.Fatalf("Failed to initialize key/value store: %s", err.Error())
	}
	
	dbStore, err := repositories.NewDBStore(config)
	if err != nil {
		log.Fatalf("Failed to initialize DB store: %s", err.Error())
	}

	server.Run(config, dbStore, kvStore)
}
