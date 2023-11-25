package main

import (
	"github.com/A1essandr0/umf-server/internal/app"
	"github.com/A1essandr0/umf-server/internal/config"
	"github.com/A1essandr0/umf-server/internal/controllers"
	"github.com/A1essandr0/umf-server/internal/repositories/db"
	"github.com/A1essandr0/umf-server/internal/repositories/kv"
	"github.com/A1essandr0/umf-server/internal/webserver"
)

func main() {
	config := config.Init("config")

	kvStore := kv.NewKVStore(config)	
	dbStore := db.NewDBStore(config)

	linksController := controllers.NewLinksController(kvStore, dbStore)
	recordsController := controllers.NewRecordsController(dbStore)
	
	webServer := webserver.NewServer(
		config,
		linksController,
		recordsController,
	)

	application := &app.AppContainer{Server: webServer}
	application.Run()
}
