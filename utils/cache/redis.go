package cache

import (
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/studio-senkou/lentera-cendekia-be/utils/app"
)

var RedisClient *redis.Client

func InitRedis() error {
	host := app.GetEnv("REDIS_HOST", "localhost")
	port := app.GetEnv("REDIS_PORT", "6379")
	password := app.GetEnv("REDIS_PASSWORD", "")
	dbStr := app.GetEnv("REDIS_DB", "0")

	db, err := strconv.Atoi(dbStr)
	if err != nil {
		return fmt.Errorf("invalid REDIS_DB value: %w", err)
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	return nil
}

func CloseRedis() {
	if RedisClient != nil {
		RedisClient.Close()
	}
}
