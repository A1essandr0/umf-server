package main

import (
	"log"

	"github.com/A1essandr0/umf-server/internal/app"
	"github.com/A1essandr0/umf-server/internal/config"
	"github.com/A1essandr0/umf-server/internal/models"
	"github.com/A1essandr0/umf-server/internal/redisclient"
	"github.com/A1essandr0/umf-server/internal/router"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	config := config.Init("config")

	routes := []router.RoutePattern{
		router.NewRoute("POST", "/create", app.CreateLink),
		router.NewRoute("GET",  "/records", app.GetRecords),
		router.NewRoute("GET",  "/([a-zA-Z0-9_-]{2,32})", app.GetLink),
	}

	redisClient, redisError := redisclient.NewRedisClient(config)
	if redisError != nil {
		log.Fatalf("Failed to connect to redis: %s", redisError.Error())
	}
	log.Printf("Got redis instance on %s", config.REDIS_ADDR)

	DB, dbError := gorm.Open(postgres.New(postgres.Config{
		DSN: config.DB_DSN,
		PreferSimpleProtocol: false,
	}), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
	})
	if dbError != nil {
		log.Fatalf("Failed to connect to DB: %s", dbError.Error())
	}
	DB.AutoMigrate(&models.NewLinkEvent{}, &models.ClickEvent{})
	log.Println("DB initialised")

	app.Run(routes, config, DB, redisClient)
}
