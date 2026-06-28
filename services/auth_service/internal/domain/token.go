package domain

import (
	"context"
	"time"
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

type TokenRepository interface {
	Set(ctx context.Context, token *RefreshToken) error
	Get(ctx context.Context, token string) (string, error)
	Delete(ctx context.Context, token string) error
}
