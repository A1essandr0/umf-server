package webserver

import (
	"github.com/A1essandr0/umf-server/internal/models"
)


func NewServer(config *models.Config) WebServer {
	server_type := "std"
	switch server_type {
	default:
		return &StdWebServer{Config: config}
	}
	
}

type WebServer interface {
	Run()
}