package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	Token  string
	UserID string
	TTL    time.Duration
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type Claims struct {
	UserID uuid.UUID
	Role   string
}

type TokenManager interface {
	GenerateAccessToken(userID uuid.UUID, role string) (string, error)
	GenerateRefreshToken(userID uuid.UUID, role string) (string, error)
	Parse(token string) (*Claims, error)
}

type TokenRepository interface {
	Set(ctx context.Context, token *RefreshToken) error
	Get(ctx context.Context, token string) (string, error)
	Delete(ctx context.Context, token string) error
}
