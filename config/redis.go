package config

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"temp/global"
)

func InitRedis(cfg *Config) error {
	if !cfg.Redis.Enabled {
		return nil // Redis is disabled, return without error
	}

	global.Redis = redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx := context.Background()
	if err := global.Redis.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to ping redis: %w", err)
	}

	return nil
}
