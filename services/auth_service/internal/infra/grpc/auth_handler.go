package grpc

import (
	"context"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/mykytakuzminov/ridely-svc/services/auth_service/internal/domain"
	authpb "github.com/mykytakuzminov/ridely-svc/shared/proto/auth"
)

type authGrpcHandler struct {
	authpb.UnimplementedAuthServiceServer
	authSvc   domain.AuthService
	validator *validator.Validate
	logger    *zap.SugaredLogger
}

func NewAuthHandler(
	authSvc domain.AuthService,
	validator *validator.Validate,
	logger *zap.SugaredLogger,
) *authGrpcHandler {
	return &authGrpcHandler{
		authSvc:   authSvc,
		validator: validator,
		logger:    logger.With("layer", "handler"),
	}
}

func (h *authGrpcHandler) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	input := &domain.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
	}
	if err := h.validator.Struct(input); err != nil {
		h.logger.Warnw("registration data validation failed", "error", err)
		return nil, err
	}

	tokens, err := h.authSvc.Register(ctx, input)
	if err != nil {
		return nil, err
	}

	return &authpb.RegisterResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
