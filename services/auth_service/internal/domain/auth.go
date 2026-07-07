package domain

import "context"

type RegisterInput struct {
	Email    string `json:"email"    validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=12,max=72"`
}

type LoginInput struct {
	Email    string `json:"email"    validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=12,max=72"`
}

type RefreshInput struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ValidateInput struct {
	AccessToken string `json:"access_token" validate:"required"`
}

type AuthService interface {
	Register(ctx context.Context, input *RegisterInput) (*TokenPair, error)
	Login(ctx context.Context, input *LoginInput) (*TokenPair, error)
	Refresh(ctx context.Context, input *RefreshInput) (*TokenPair, error)
	Validate(ctx context.Context, input *ValidateInput) (*Claims, error)
}
