package db

import (
	"log"

	"github.com/A1essandr0/umf-server/internal/models"
	"github.com/A1essandr0/umf-server/internal/repositories"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDBStore(config *models.Config) repositories.DBStore {
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
					log.Fatalf("Couldn't apply or recognize Postgres DB schema: %+v", err)
				}
				log.Println("Postgres migrations applied")
			}
			
			log.Println("Postgres DB initialised")
			return &PStore{
				db: DB,
				DEFAULT_RECORDS_AMOUNT_TO_GET: config.DEFAULT_RECORDS_AMOUNT_TO_GET,
			}
		
		// sqlite
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
				db: DB,
				DEFAULT_RECORDS_AMOUNT_TO_GET: config.DEFAULT_RECORDS_AMOUNT_TO_GET,
			}

		case "sqlite-inmemory":
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
				db: DB,
				DEFAULT_RECORDS_AMOUNT_TO_GET: config.DEFAULT_RECORDS_AMOUNT_TO_GET,
			}

		default:
			log.Println("Using pure mock DB")
			return &MockDB{}
	}
}