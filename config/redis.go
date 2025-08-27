// internal/config/redis.go

package config

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var (
	Ctx         = context.Background()
	RedisClient *redis.Client
)

func ConnectRedis() {

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     AppConfig.Redis.Host,
		Password: AppConfig.Redis.Password,
		DB:       AppConfig.Redis.DB,
	})
	log.Println("✅ Connected to Redis")
}
