package redis

import (
	"Go-REST-API-Portfolio/config"
	"github.com/go-redis/redis"
	"time"
)

func NewRedisClient(c *config.Config) *redis.Client {
	host := c.Redis.RedisAddr
	if host == "" {
		host = ":6379"
	}

	return redis.NewClient(&redis.Options{
		Addr:         host,
		MinIdleConns: c.Redis.MinIdleConns,
		PoolSize:     c.Redis.PoolSize,
		PoolTimeout:  time.Duration(c.Redis.PoolTimeout) * time.Second,
		Password:     c.Redis.Password,
		DB:           c.Redis.DB,
	})
}