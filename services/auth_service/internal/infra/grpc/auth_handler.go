package grpc

import (
	"context"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mykytakuzminov/ridely-svc/services/auth_service/internal/domain"
	c "github.com/mykytakuzminov/ridely-svc/shared/context"
	"github.com/mykytakuzminov/ridely-svc/shared/errors"
	log "github.com/mykytakuzminov/ridely-svc/shared/logging"
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
	traceID := c.GetTraceID(ctx)

	input := &domain.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
	}
	if err := h.validator.Struct(input); err != nil {
		log.FailedValidateData(h.logger, traceID, err)
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

func (h *authGrpcHandler) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	traceID := c.GetTraceID(ctx)

	input := &domain.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}
	if err := h.validator.Struct(input); err != nil {
		log.FailedValidateData(h.logger, traceID, err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	tokens, err := h.authSvc.Login(ctx, input)
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}

	return &authpb.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
