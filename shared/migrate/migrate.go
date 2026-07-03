package migrate

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	pgxmigrate "github.com/golang-migrate/migrate/v4/database/pgx/v5"
)

func Run(pool *pgxpool.Pool, migrationsFS embed.FS, logger *zap.SugaredLogger) error {
	src, err := iofs.New(migrationsFS, ".")
	if err != nil {
		return fmt.Errorf("load embedded migrations: %w", err)
	}

	db := stdlib.OpenDBFromPool(pool)
	defer func() {
		if err := db.Close(); err != nil {
			logger.Warnw("failed to close migration db wrapper", "error", err)
		}
	}()

	driver, err := pgxmigrate.WithInstance(db, &pgxmigrate.Config{})
	if err != nil {
		return fmt.Errorf("init driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", src, "pgx5", driver)
	if err != nil {
		return fmt.Errorf("init migrate: %w", err)
	}
	defer func() {
		if srcErr, dbErr := m.Close(); srcErr != nil || dbErr != nil {
			logger.Warnw("failed to close migration",
				"source_error", srcErr,
				"database_error", dbErr,
			)
		}
	}()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("apply migrations: %w", err)
	}

	return nil
}
