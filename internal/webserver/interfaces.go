package webserver

import (
	"github.com/A1essandr0/umf-server/internal/controllers"
	"github.com/A1essandr0/umf-server/internal/models"
)

func NewServer(
	config *models.Config, 
	linksController *controllers.LinksController,
	recordsController *controllers.RecordsController,
) WebServer {
	server_type := "std"
	switch server_type {
	default:
		return &StdWebServer{
			Config: config,
			LinksController: linksController,
			RecordsController: recordsController,
		}
	}
}

type WebServer interface {
	Run()
}