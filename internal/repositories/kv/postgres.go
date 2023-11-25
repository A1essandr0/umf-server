package kv

import "gorm.io/gorm"


type PostgresKVClient struct {
	db *gorm.DB
}


func (c *PostgresKVClient) CreateKVStoreRecord(key, value string) error {


	return nil
}

func (c *PostgresKVClient) GetKVStoreRecord(key string) (string, error) {


	return "", nil
}

