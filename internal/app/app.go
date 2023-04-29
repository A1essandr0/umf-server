package app

import (
	"errors"
	"log"
	"net/http"

	"github.com/A1essandr0/umf-server/internal/config"
	"github.com/A1essandr0/umf-server/internal/redisclient"
	"github.com/A1essandr0/umf-server/internal/router"
	"gorm.io/gorm"
)

type ctxKey struct{}

var (
	Routes []router.RoutePattern
	RedisClient *redisclient.RedisClient
	DB *gorm.DB
)

func Run(routes []router.RoutePattern, db *gorm.DB, redisClient *redisclient.RedisClient) {
	DB = db
	RedisClient = redisClient
	Routes = routes

	var httpStartError error
	server := &http.Server{
		Addr: config.WEB_PORT,
		Handler: corsMiddleWare(http.HandlerFunc(serve)),
	}
	if config.USE_TLS == "true" {
		log.Printf("Starting https server on %s, mode: %s", config.WEB_PORT, config.DEVELOPMENT_MODE)
		httpStartError = server.ListenAndServeTLS(config.CERT_FILE, config.CERT_KEY_FILE)
	} else {
		log.Printf("Starting http server on %s", config.WEB_PORT)
		httpStartError = server.ListenAndServe()
	}

	if errors.Is(httpStartError, http.ErrServerClosed) {
		log.Println("... server stopped")
	} else if httpStartError != nil {
		log.Printf("Error starting server: %s", httpStartError)
	}
}