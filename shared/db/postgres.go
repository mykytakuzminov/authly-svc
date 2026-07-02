package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mykytakuzminov/ridely-svc/shared/env"
	"github.com/mykytakuzminov/ridely-svc/shared/errors"
)

type postgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func newPostgresConfig() (*postgresConfig, error) {
	cfg := &postgresConfig{}
	cfg.load()
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *postgresConfig) load() {
	c.Host = env.GetString("POSTGRES_HOST", "localhost")
	c.Port = env.GetString("POSTGRES_PORT", "5432")
	c.User = env.GetString("POSTGRES_USER", "")
	c.Password = env.GetString("POSTGRES_PASSWORD", "")
	c.DBName = env.GetString("POSTGRES_DB", "")
	c.SSLMode = env.GetString("POSTGRES_SSLMODE", "disable")
}

func (c *postgresConfig) validate() error {
	if c.User == "" {
		return fmt.Errorf("POSTGRES_USER %w", errors.ErrIsRequired)
	}
	if c.Password == "" {
		return fmt.Errorf("POSTGRES_PASSWORD %w", errors.ErrIsRequired)
	}
	if c.DBName == "" {
		return fmt.Errorf("POSTGRES_DB %w", errors.ErrIsRequired)
	}
	return nil
}

func (c *postgresConfig) getDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
		c.SSLMode,
	)
}

func newPostgresClient(ctx context.Context, cfg *postgresConfig) (*pgxpool.Pool, error) {
	client, err := pgxpool.New(ctx, cfg.getDSN())
	if err != nil {
		return nil, err
	}
	if err = client.Ping(ctx); err != nil {
		return nil, err
	}

	return client, nil
}

func ConnectPostgres(ctx context.Context, timeout time.Duration) (*pgxpool.Pool, error) {
	cfg, err := newPostgresConfig()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	connCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	client, err := newPostgresClient(connCtx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}
	return client, nil
}
