package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"

	c "github.com/mykytakuzminov/ridely-svc/shared/context"
	log "github.com/mykytakuzminov/ridely-svc/shared/logging"

	"github.com/mykytakuzminov/ridely-svc/services/auth_service/internal/domain"
)

type userRepository struct {
	db     domain.DB
	logger *zap.SugaredLogger
}

func NewUserRepository(db domain.DB, logger *zap.SugaredLogger) domain.UserRepository {
	return &userRepository{
		db:     db,
		logger: logger.With("layer", "repository"),
	}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	traceID := c.GetTraceID(ctx)

	query := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, role, created_at, updated_at
	`

	row := r.db.QueryRow(ctx, query, user.Email, user.PasswordHash)
	if err := r.scanUser(row, traceID,
		&user.ID,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return err
	}

	log.Success(r.logger, traceID, "user creation", "user_id", user.ID)
	return nil
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	traceID := c.GetTraceID(ctx)

	query := `
		SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		log.Failed(r.logger, traceID, "existence scanning", err)
		return false, err
	}

	return exists, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	traceID := c.GetTraceID(ctx)

	query := `
		SELECT id, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	row := r.db.QueryRow(ctx, query, email)
	user := &domain.User{}
	if err := r.scanUser(row, traceID,
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) scanUser(row pgx.Row, traceID uuid.UUID, dest ...any) error {
	if err := row.Scan(dest...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrUserNotFound
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrUserAlreadyExists
		}

		log.Failed(r.logger, traceID, "user scanning", err)
		return err
	}
	return nil
}
