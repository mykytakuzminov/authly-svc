package domain

import "context"

type AuthService interface {
	Register(ctx context.Context, input *RegisterInput) (*TokenPair, error)
}
