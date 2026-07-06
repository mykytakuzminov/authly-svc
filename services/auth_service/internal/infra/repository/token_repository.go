package repository

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	c "github.com/mykytakuzminov/ridely-svc/shared/context"
	log "github.com/mykytakuzminov/ridely-svc/shared/logging"

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
	traceID := c.GetTraceID(ctx)

	if err := r.db.Set(ctx, token.Token, token.UserID.String(), token.TTL).Err(); err != nil {
		log.Failed(r.logger, traceID, "refresh token saving", err)
		return err
	}

	log.Success(r.logger, traceID, "refresh token saving", "user_id", token.UserID)
	return nil
}

func (r *tokenRepository) Get(ctx context.Context, token string) (string, error) {
	traceID := c.GetTraceID(ctx)

	userID, err := r.db.Get(ctx, token).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", domain.ErrTokenNotFound
		}

		log.Failed(r.logger, traceID, "refresh token retrieval", err)
		return "", err
	}
	return userID, nil
}

func (r *tokenRepository) Delete(ctx context.Context, token string) error {
	traceID := c.GetTraceID(ctx)

	if err := r.db.Del(ctx, token).Err(); err != nil {
		log.Failed(r.logger, traceID, "refresh token deletion", err)
		return err
	}

	log.Success(r.logger, traceID, "refresh token deletion")
	return nil
}
