package cache

import (
	"context"
	"encoding/json"
	"time"
)

func Set(ctx context.Context, key string, value any, expire time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return RedisClient.Set(ctx, key, jsonData, expire).Err()
}

func Get(ctx context.Context, key string, dest any) error {
	value, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(value), dest)
}

func Delete(ctx context.Context, key string) error {
	return RedisClient.Del(ctx, key).Err()
}

func Exists(ctx context.Context, key string) (bool, error) {
	exists, err := RedisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}
