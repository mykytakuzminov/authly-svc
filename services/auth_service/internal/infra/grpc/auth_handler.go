package grpc

import (
	"context"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mykytakuzminov/ridely-svc/services/auth_service/internal/domain"
	"github.com/mykytakuzminov/ridely-svc/shared/errors"
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
) authpb.AuthServiceServer {
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
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	tokens, err := h.authSvc.Register(ctx, input)
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}

	return &authpb.RegisterResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
