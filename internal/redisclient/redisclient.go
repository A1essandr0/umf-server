package redisclient

import (
	"time"

	"github.com/A1essandr0/umf-server/internal/models"
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


func NewRedisClient(config models.Config) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.REDIS_ADDR, 
		Password: config.REDIS_PWD,
		DB:       config.REDIS_DB_NUM,
	 })
	 if err := client.Ping(client.Context()).Err(); err != nil {
		return nil, err
	 }
	 return &RedisClient{client, config.DEFAULT_TTL}, nil
}

