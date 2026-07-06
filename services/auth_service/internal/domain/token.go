package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	Token  string
	UserID uuid.UUID
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
	GenerateAccessToken(ctx context.Context, userID uuid.UUID, role string) (string, error)
	GenerateRefreshToken(ctx context.Context, userID uuid.UUID, role string) (string, error)
	Parse(ctx context.Context, token string) (*Claims, error)
	GetRefreshTTL() time.Duration
}

type TokenRepository interface {
	Set(ctx context.Context, token *RefreshToken) error
	Get(ctx context.Context, token string) (string, error)
	Delete(ctx context.Context, token string) error
}
