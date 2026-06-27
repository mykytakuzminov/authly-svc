package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mykytakuzminov/ridely-svc/shared/env"
	"github.com/mykytakuzminov/ridely-svc/shared/errors"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresConfig() (*PostgresConfig, error) {
	cfg := &PostgresConfig{}
	cfg.load()
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *PostgresConfig) load() {
	c.Host = env.GetString("POSTGRES_HOST", "localhost")
	c.Port = env.GetString("POSTGRES_PORT", "5432")
	c.User = env.GetString("POSTGRES_USER", "")
	c.Password = env.GetString("POSTGRES_PASSWORD", "")
	c.DBName = env.GetString("POSTGRES_DB", "")
	c.SSLMode = env.GetString("POSTGRES_SSLMODE", "disable")
}

func (c *PostgresConfig) validate() error {
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

func (c *PostgresConfig) getDSN() string {
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

func NewPostgresClient(ctx context.Context, cfg *PostgresConfig) (*pgxpool.Pool, error) {
	url := cfg.getDSN()

	client, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}
	if err = client.Ping(ctx); err != nil {
		return nil, err
	}

	return client, nil
}
