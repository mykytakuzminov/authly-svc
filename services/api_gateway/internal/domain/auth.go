package domain

import authpb "github.com/mykytakuzminov/ridely-svc/shared/proto/auth"

type RegisterInput struct {
	Email    string `json:"email"    validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=12,max=72"`
}

func (i *RegisterInput) ToProto() *authpb.RegisterRequest {
	return &authpb.RegisterRequest{
		Email:    i.Email,
		Password: i.Password,
	}
}

type LoginInput struct {
	Email    string `json:"email"    validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=12,max=72"`
}

func (i *LoginInput) ToProto() *authpb.LoginRequest {
	return &authpb.LoginRequest{
		Email:    i.Email,
		Password: i.Password,
	}
}

type RefreshInput struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (i *RefreshInput) ToProto() *authpb.RefreshRequest {
	return &authpb.RefreshRequest{
		RefreshToken: i.RefreshToken,
	}
}

type LogoutInput struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (i *LogoutInput) ToProto() *authpb.LogoutRequest {
	return &authpb.LogoutRequest{
		RefreshToken: i.RefreshToken,
	}
}

type TokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func TokensResponseFromProto(t *authpb.TokenPair) TokensResponse {
	return TokensResponse{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
	}
}
