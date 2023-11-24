package repositories

import (
	"log"

	"github.com/A1essandr0/umf-server/internal/models"
	"github.com/glebarez/sqlite"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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


func NewKVStore(config *models.Config) KeyValueStore {
	switch config.KVSTORE_TYPE {
		case "redis":
			client := redis.NewClient(&redis.Options{
				Addr:     config.REDIS_ADDR, 
				Password: config.REDIS_PWD,
				DB:       config.REDIS_DB_NUM,
			})
			if err := client.Ping(client.Context()).Err(); err != nil {
				log.Fatalf("Failed to initialise key-value store: %+v", err)
			}
			log.Printf("Got key/value store instance up on %s", config.REDIS_ADDR)
			return &RedisClient{client, config.DEFAULT_TTL}
		
		default:
			log.Println("Using inmemory key-value store")
			mapStore := make(map[string]string)
			return &InmemoryKV{store: mapStore}
	}
}


func NewDBStore(config *models.Config) DBStore {
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}
	if config.DB_DEBUG_LOG {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}


	switch config.DBSTORE_TYPE {		
		case "postgres":
			DB, err := gorm.Open(postgres.New(postgres.Config{
				DSN: config.DB_DSN,
				PreferSimpleProtocol: false,
			}), gormConfig)
			if err != nil {
				log.Fatalf("Failed to initialise database: %+v", err)
			}

			if config.APPLY_MIGRATIONS {
                if err = DB.AutoMigrate(models.ModelsToAutoMigrate...); err != nil {
					log.Fatalf("Couldn't apply or recognize Postgres DB schema^ %+v", err)
				}
				log.Println("Postgres migrations applied")
			}
			
			log.Println("Postgres DB initialised")
			return &PStore{
				DB: DB,
				DEFAULT_RECORDS_AMOUNT_TO_GET: config.DEFAULT_RECORDS_AMOUNT_TO_GET,
			}
		
		case "sqlite":
			DB, err := gorm.Open(sqlite.Open(config.DB_FILE), gormConfig)
			if err != nil {
				log.Fatalf("Failed to initialise database: %+v", err)
			}
			if err = DB.AutoMigrate(models.ModelsToAutoMigrate...); err != nil {
				log.Fatalf("Couldn't apply or recognize sqlite DB schema: %+v", err)
			}
			log.Println("Sqlite migrations applied")

			log.Println("Sqlite DB initialised")
			return &MStore{
				DB: DB,
				DEFAULT_RECORDS_AMOUNT_TO_GET: config.DEFAULT_RECORDS_AMOUNT_TO_GET,
			}

		default:
			DB, err := gorm.Open(sqlite.Open(":memory:"), gormConfig)
			if err != nil {
				log.Fatalf("Failed to initialise database: %+v", err)
			}
			if err = DB.AutoMigrate(models.ModelsToAutoMigrate...); err != nil {
				log.Fatalf("Couldn't apply or recognize inmemory DB schema: %+v", err)
			}
			log.Println("Inmemory migrations applied")

			log.Println("Inmemory DB initialised")
			return &MStore{
				DB: DB,
				DEFAULT_RECORDS_AMOUNT_TO_GET: config.DEFAULT_RECORDS_AMOUNT_TO_GET,
			}
	}
}