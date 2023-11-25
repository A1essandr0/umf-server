package kv

import (
	"log"

	"github.com/A1essandr0/umf-server/internal/models"
	"github.com/A1essandr0/umf-server/internal/repositories"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewKVStore(config *models.Config) repositories.KeyValueStore {

	switch config.KVSTORE_TYPE {		
		case "postgres":
			gormConfig := &gorm.Config{
				Logger: logger.Default.LogMode(logger.Silent),
			}
			if config.DB_DEBUG_LOG {
				gormConfig.Logger = logger.Default.LogMode(logger.Info)
			}
		
			DB, err := gorm.Open(postgres.New(postgres.Config{
				DSN: config.DB_DSN,
				PreferSimpleProtocol: false,
			}), gormConfig)
			if err != nil {
				log.Fatalf("Failed to initialise key-value store via Postgres: %+v", err)
			}

			if config.APPLY_MIGRATIONS {
                if err = DB.AutoMigrate(models.ModelsToAutoMigrate...); err != nil {
					log.Fatalf("Couldn't apply or recognize Postgres schema for KV: %+v", err)
				}
				log.Println("Postgres migrations for KV applied")
			}

			log.Println("Postgres store for KV initialised")
			return &PostgresKVClient{db: DB}


		case "redis":
			client := redis.NewClient(&redis.Options{
				Addr:     config.REDIS_ADDR, 
				Password: config.REDIS_PWD,
				DB:       config.REDIS_DB_NUM,
			})
			if err := client.Ping(client.Context()).Err(); err != nil {
				log.Fatalf("Failed to initialise key-value store: %+v", err)
			}
			log.Printf("Got Redis key-value store instance up on %s", config.REDIS_ADDR)
			return &RedisClient{client, config.DEFAULT_TTL}


		default:
			log.Println("Using inmemory key-value store")
			mapStore := make(map[string]string)
			return &InmemoryKV{store: mapStore}
	}
}
