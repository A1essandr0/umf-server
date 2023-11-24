package app

import (
	"github.com/A1essandr0/umf-server/internal/webserver"
)

type AppContainer struct {
	Server webserver.WebServer
}

func (a *AppContainer) Run() {
	a.Server.Run()
}