package domain

import "context"

type CredentialsInput struct {
	Email    string `json:"email"    validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=12,max=72"`
}

type RegisterInput = CredentialsInput
type LoginInput = CredentialsInput

type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshInput = RefreshTokenInput
type LogoutInput = RefreshTokenInput

type ValidateInput struct {
	AccessToken string `json:"access_token" validate:"required"`
}

type AuthService interface {
	Register(ctx context.Context, input *RegisterInput) (*TokenPair, error)
	Login(ctx context.Context, input *LoginInput) (*TokenPair, error)
	Refresh(ctx context.Context, input *RefreshInput) (*TokenPair, error)
	Logout(ctx context.Context, input *LogoutInput) error
	Validate(ctx context.Context, input *ValidateInput) (*Claims, error)
}
