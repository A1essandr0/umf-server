package server

import (
	"errors"
	"log"
	"net/http"

	"github.com/A1essandr0/umf-server/internal/models"
	"github.com/A1essandr0/umf-server/internal/repositories"
)

var (
	KVClient repositories.KeyValueStore
	DB repositories.DBStore
	Config models.Config
)

func Run(conf models.Config, db repositories.DBStore, kvClient repositories.KeyValueStore) {
	Config = conf
	DB = db
	KVClient = kvClient

	var httpStartError error
	server := &http.Server{
		Addr: Config.WEB_PORT,
		Handler: corsMiddleWare(
			http.HandlerFunc(
				ServeMux,
			),
		),
	}

	if Config.USE_TLS {
		log.Printf("Starting https server on %s, mode: %s", Config.WEB_PORT, Config.DEVELOPMENT_MODE)
		httpStartError = server.ListenAndServeTLS(Config.CERT_FILE, Config.CERT_KEY_FILE)
	} else {
		log.Printf("Starting http server on %s", Config.WEB_PORT)
		httpStartError = server.ListenAndServe()
	}

	if errors.Is(httpStartError, http.ErrServerClosed) {
		log.Println("... server stopped")
	} else if httpStartError != nil {
		log.Printf("Error starting server: %s", httpStartError)
	}
}