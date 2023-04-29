package app

import (
	"errors"
	"log"
	"net/http"

	"github.com/A1essandr0/umf-server/internal/models"
	"github.com/A1essandr0/umf-server/internal/redisclient"
	"github.com/A1essandr0/umf-server/internal/router"
	"gorm.io/gorm"
)

type ctxKey struct{}

var (
	Routes []router.RoutePattern
	RedisClient *redisclient.RedisClient
	DB *gorm.DB
	Config models.Config
)

func Run(routes []router.RoutePattern, conf models.Config, db *gorm.DB, redisClient *redisclient.RedisClient) {
	DB = db
	RedisClient = redisClient
	Routes = routes
	Config = conf

	var httpStartError error
	server := &http.Server{
		Addr: Config.WEB_PORT,
		Handler: corsMiddleWare(http.HandlerFunc(serve)),
	}
	if Config.USE_TLS == true {
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