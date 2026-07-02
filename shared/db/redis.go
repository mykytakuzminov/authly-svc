package db

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/mykytakuzminov/ridely-svc/shared/env"
	"github.com/mykytakuzminov/ridely-svc/shared/errors"
)

type redisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func newRedisConfig() (*redisConfig, error) {
	cfg := &redisConfig{}
	cfg.load()
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *redisConfig) load() {
	c.Host = env.GetString("REDIS_HOST", "localhost")
	c.Port = env.GetString("REDIS_PORT", "6379")
	c.Password = env.GetString("REDIS_PASSWORD", "")
	c.DB = env.GetInt("REDIS_DB", 0)
}

func (c *redisConfig) validate() error {
	if c.Password == "" {
		return fmt.Errorf("REDIS_PASSWORD %w", errors.ErrIsRequired)
	}
	return nil
}

func (c *redisConfig) getAddr() string {
	return c.Host + ":" + c.Port
}

func newRedisClient(ctx context.Context, cfg *redisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.getAddr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}

func ConnectRedis(ctx context.Context, timeout time.Duration) (*redis.Client, error) {
	cfg, err := newRedisConfig()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	connCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	client, err := newRedisClient(connCtx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	return client, nil
}
