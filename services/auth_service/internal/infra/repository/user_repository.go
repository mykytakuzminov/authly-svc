package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

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
	query := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, role, created_at, updated_at
	`

	row := r.db.QueryRow(ctx, query, user.Email, user.PasswordHash)
	if err := r.scanUser(row, user); err != nil {
		return err
	}

	r.logger.Infow("user created", "user_id", user.ID)
	return nil
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `
		SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, domain.ErrUserNotFound
		}
		r.logger.Errorw("failed to scan existence", "error", err)
		return false, err
	}

	return exists, nil
}

func (r *userRepository) scanUser(row pgx.Row, user *domain.User) error {
	if err := row.Scan(
		&user.ID,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrUserNotFound
		}
		r.logger.Errorw("failed to scan user", "error", err)
		return err
	}
	return nil
}
