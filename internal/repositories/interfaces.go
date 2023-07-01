package repositories

import (
	"log"

	"github.com/A1essandr0/umf-server/internal/models"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type KeyValueStore interface {
	CreateKVStoreRecord(key, value string) error
	GetKVStoreRecord(key string) (string, error)
}

type DBStore interface {
	CreateClickEvent(link, value, ip string)
	CreateNewLinkEvent(link, url, ip string)
	GetNewLinkEvents(ip string) []models.NewLinkEvent
}


func NewKVStore(config models.Config) (KeyValueStore, error) {
	switch config.KVSTORE_TYPE {

		default:
			client := redis.NewClient(&redis.Options{
				Addr:     config.REDIS_ADDR, 
				Password: config.REDIS_PWD,
				DB:       config.REDIS_DB_NUM,
			})
			if err := client.Ping(client.Context()).Err(); err != nil {
				return nil, err
			}
			log.Printf("Got key/value store instance up on %s", config.REDIS_ADDR)
			return &RedisClient{client, config.DEFAULT_TTL}, nil
	}
}


func NewDBStore(config models.Config) (DBStore, error) {
	switch config.DBSTORE_TYPE {
		
		default:
			DB, err := gorm.Open(postgres.New(postgres.Config{
				DSN: config.DB_DSN,
				PreferSimpleProtocol: false,
			}), &gorm.Config{
				// Logger: logger.Default.LogMode(logger.Info),
			})
			if err != nil {
				return nil, err
			}

			if config.APPLY_MIGRATIONS {
				DB.AutoMigrate(models.ModelsToAutoMigrate...)
			}
			
			log.Println("Postgres DB initialised")
			return &PStore{
				DB: DB,
				DEFAULT_RECORDS_AMOUNT_TO_GET: config.DEFAULT_RECORDS_AMOUNT_TO_GET,
			}, nil
	}
}