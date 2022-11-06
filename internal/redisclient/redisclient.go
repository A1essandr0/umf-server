package redisclient

import (
	"time"

	"github.com/A1essandr0/umf-server/internal/config"
	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	*redis.Client
}

func (c *RedisClient) CreateKVStoreRecord(key, value string) error {
	return c.Set(c.Context(), key, value,
		// convert nanoseconds to hours
		time.Duration(config.DEFAULT_TTL)*1000*1000*1000*3600,
	).Err()
}

func (c *RedisClient) GetKVStoreRecord(key string) (string, error) { 
	value, err := c.Get(c.Context(), key).Result()
	return value, err
}


func NewRedisClient() (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.REDIS_ADDR, 
		Password: config.REDIS_PWD,
		DB:       config.REDIS_DB_NUM,
	 })
	 if err := client.Ping(client.Context()).Err(); err != nil {
		return nil, err
	 }
	 return &RedisClient{client}, nil
}

