package domain

import "context"

type RegisterInput struct {
	Email    string `json:"email"    validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=12,max=72"`
}

type AuthService interface {
	Register(ctx context.Context, input *RegisterInput) (*TokenPair, error)
}
