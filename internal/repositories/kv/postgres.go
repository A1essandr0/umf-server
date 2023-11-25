package kv

import (
	"github.com/A1essandr0/umf-server/internal/models"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)


type PostgresKVClient struct {
	db *gorm.DB
}


func (c *PostgresKVClient) CreateKVStoreRecord(key, value string) error {
	// 

	return nil
}

func (c *PostgresKVClient) GetKVStoreRecord(key string) (string, error) {
	var records []models.NewLinkEvent
	err := c.db.Where(&models.NewLinkEvent{Key: key},
		).Order("created_at desc",
		).Find(&records).Error
	if err != nil {
		return "", err
	}
	if len(records) == 0 {
		return "", redis.Nil
	}

	return records[0].Value, nil
}
