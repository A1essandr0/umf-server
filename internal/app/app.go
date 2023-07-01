package app

import (
	"errors"
	"log"
	"net/http"

	"github.com/A1essandr0/umf-server/internal/models"
	"github.com/A1essandr0/umf-server/internal/repositories"
	"github.com/A1essandr0/umf-server/internal/router"
)

type ctxKey struct{}

var (
	Routes []router.RoutePattern
	KVClient repositories.KeyValueStore
	DB repositories.DBStore
	Config models.Config
)

func Run(
	conf models.Config,
	routes []router.RoutePattern,
	db repositories.DBStore,
	kvClient repositories.KeyValueStore,
) {
	DB = db
	KVClient = kvClient
	Routes = routes
	Config = conf

	var httpStartError error
	server := &http.Server{
		Addr: Config.WEB_PORT,
		Handler: corsMiddleWare(http.HandlerFunc(serve)),
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