package main

import (
	"log"

	"github.com/A1essandr0/umf-server/internal/app"
	"github.com/A1essandr0/umf-server/internal/config"
	"github.com/A1essandr0/umf-server/internal/repositories"
	"github.com/A1essandr0/umf-server/internal/router"
)

func main() {
	config := config.Init("config")

	routes := []router.RoutePattern{
		router.NewRoute("POST", "/create", app.CreateLink),
		router.NewRoute("GET",  "/records", app.GetRecords),
		router.NewRoute("GET",  "/([a-zA-Z0-9_-]{2,32})", app.GetLink),
	}

	kvStore, err := repositories.NewKVStore(config)
	if err != nil {
		log.Fatalf("Failed to initialize key/value store: %s", err.Error())
	}
	
	dbStore, err := repositories.NewDBStore(config)
	if err != nil {
		log.Fatalf("Failed to initialize DB store: %s", err.Error())
	}

	app.Run(config, routes, dbStore, kvStore)
}
