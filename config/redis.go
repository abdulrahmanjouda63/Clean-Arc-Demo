package config

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"temp/global"
)

// InitRedis initializes Redis client with support for Sentinel
func InitRedis(cfg *Config) error {
	if !cfg.Redis.Enabled {
		global.Logger.Info("Redis is disabled")
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use Sentinel if configured
	if cfg.Redis.UseSentinel {
		global.Logger.Info("Initializing Redis with Sentinel",
			zap.Strings("sentinel_addrs", cfg.Redis.SentinelAddrs),
			zap.String("master", cfg.Redis.SentinelMaster),
		)

		global.Redis = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:       cfg.Redis.SentinelMaster,
			SentinelAddrs:    cfg.Redis.SentinelAddrs,
			SentinelPassword: cfg.Redis.SentinelPassword,
			Password:         cfg.Redis.Password,
			DB:               cfg.Redis.DB,
			DialTimeout:      5 * time.Second,
			ReadTimeout:      3 * time.Second,
			WriteTimeout:     3 * time.Second,
			PoolSize:         10,
			MinIdleConns:     5,
		})
	} else {
		global.Logger.Info("Initializing Redis standalone",
			zap.String("addr", cfg.Redis.Addr),
		)

		global.Redis = redis.NewClient(&redis.Options{
			Addr:         cfg.Redis.Addr,
			Password:     cfg.Redis.Password,
			DB:           cfg.Redis.DB,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
			PoolSize:     10,
			MinIdleConns: 5,
		})
	}

	// Ping to verify connection
	if err := global.Redis.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to ping redis: %w", err)
	}

	global.Logger.Info("Redis initialized successfully")
	return nil
}

// CloseRedis gracefully closes Redis connection
func CloseRedis() error {
	if global.Redis != nil {
		global.Logger.Info("Closing Redis connection")
		return global.Redis.Close()
	}
	return nil
}
