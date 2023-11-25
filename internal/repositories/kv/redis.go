package kv

import (
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	*redis.Client
	TTL int
}


func (c *RedisClient) CreateKVStoreRecord(key, value string) error {
	return c.Set(c.Context(), key, value,
		// convert nanoseconds to hours
		time.Duration(c.TTL)*1000*1000*1000*3600,
	).Err()
}

func (c *RedisClient) GetKVStoreRecord(key string) (string, error) { 
	value, err := c.Get(c.Context(), key).Result()
	return value, err
}
