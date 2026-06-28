package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/mykytakuzminov/ridely-svc/shared/env"
	"github.com/mykytakuzminov/ridely-svc/shared/errors"
)

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func NewRedisConfig() (*RedisConfig, error) {
	cfg := &RedisConfig{}
	cfg.load()
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *RedisConfig) load() {
	c.Host = env.GetString("REDIS_HOST", "localhost")
	c.Port = env.GetString("REDIS_PORT", "6379")
	c.Password = env.GetString("REDIS_PASSWORD", "")
	c.DB = env.GetInt("REDIS_DB", 0)
}

func (c *RedisConfig) validate() error {
	if c.Password == "" {
		return fmt.Errorf("REDIS_PASSWORD %w", errors.ErrIsRequired)
	}
	return nil
}

func (c *RedisConfig) getAddr() string {
	return c.Host + ":" + c.Port
}

func NewRedisClient(ctx context.Context, cfg *RedisConfig) (*redis.Client, error) {
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
