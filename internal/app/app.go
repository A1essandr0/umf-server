package app

import (
	"github.com/A1essandr0/umf-server/internal/models"
	"github.com/A1essandr0/umf-server/internal/repositories"
	"github.com/A1essandr0/umf-server/internal/webserver"
)



type AppContainer struct {
	Config *models.Config
	KVStore repositories.KeyValueStore
	DB repositories.DBStore
	Server webserver.WebServer
}

func (a *AppContainer) Run() {
	

	a.Server.Run()
}