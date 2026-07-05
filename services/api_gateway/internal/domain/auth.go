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
