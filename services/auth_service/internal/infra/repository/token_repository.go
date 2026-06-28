package repository

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/mykytakuzminov/ridely-svc/services/auth_service/internal/domain"
)

type tokenRepository struct {
	db     *redis.Client
	logger *zap.SugaredLogger
}

func NewTokenRepository(db *redis.Client, logger *zap.SugaredLogger) domain.TokenRepository {
	return &tokenRepository{
		db:     db,
		logger: logger.With("layer", "repository"),
	}
}

func (r *tokenRepository) Set(ctx context.Context, token *domain.RefreshToken) error {
	if err := r.db.Set(ctx, token.Token, token.UserID.String(), token.TTL).Err(); err != nil {
		r.logger.Errorw("failed to set refresh token", "error", err)
		return err
	}

	r.logger.Infow("refresh token set successfully", "user_id", token.UserID)
	return nil
}

func (r *tokenRepository) Get(ctx context.Context, token string) (string, error) {
	userID, err := r.db.Get(ctx, token).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", domain.ErrTokenNotFound
		}

		r.logger.Errorw("failed to get refresh token", "error", err)
		return "", err
	}
	return userID, nil
}

func (r *tokenRepository) Delete(ctx context.Context, token string) error {
	if err := r.db.Del(ctx, token).Err(); err != nil {
		r.logger.Errorw("failed to delete refresh token", "error", err)
		return err
	}

	r.logger.Infow("refresh token deleted successfully")
	return nil
}
