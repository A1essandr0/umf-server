package main

import (
	"github.com/A1essandr0/umf-server/internal/app"
	"github.com/A1essandr0/umf-server/internal/config"
	"github.com/A1essandr0/umf-server/internal/repositories"
	"github.com/A1essandr0/umf-server/internal/server"
	"github.com/A1essandr0/umf-server/internal/webserver"
)

func main() {
	config := config.Init("config")

	kvStore := repositories.NewKVStore(config)	
	dbStore := repositories.NewDBStore(config)
	webServer := webserver.NewServer(config)

	application := &app.AppContainer{
		KVStore: kvStore,
		DB: dbStore,
		Config: config,
		Server: webServer,
	}
	application.Run()



	// TO BE DEPRECATED
	server.Run(config, dbStore, kvStore)
}
