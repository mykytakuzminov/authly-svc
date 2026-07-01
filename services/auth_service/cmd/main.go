package main

import (
	"context"
	"net"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	grpclib "google.golang.org/grpc"

	"github.com/mykytakuzminov/ridely-svc/services/auth_service/internal/infra/grpc"
	"github.com/mykytakuzminov/ridely-svc/services/auth_service/internal/infra/jwt"
	"github.com/mykytakuzminov/ridely-svc/services/auth_service/internal/infra/repository"
	"github.com/mykytakuzminov/ridely-svc/services/auth_service/internal/service"
	"github.com/mykytakuzminov/ridely-svc/shared/postgres"
	authpb "github.com/mykytakuzminov/ridely-svc/shared/proto/auth"
	"github.com/mykytakuzminov/ridely-svc/shared/redis"
)

func main() {
	l := zap.Must(zap.NewProduction())
	logger := l.Sugar()
	validator := validator.New()

	if err := godotenv.Load(); err != nil {
		logger.Warnw("no .env file, reading from environment")
	}

	pgCfg, err := postgres.NewPostgresConfig()
	if err != nil {
		logger.Fatalw("failed to load postgres config", "error", err)
	}
	rdCfg, err := redis.NewRedisConfig()
	if err != nil {
		logger.Fatalw("failed to load redis config", "error", err)
	}
	jwtCfg, err := jwt.NewJWTConfig()
	if err != nil {
		logger.Fatalw("failed to load jwt config", "error", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pgClient, err := postgres.NewPostgresClient(ctx, pgCfg)
	if err != nil {
		logger.Fatalw("failed to establish postgres connection", "error", err)
	}
	rdClient, err := redis.NewRedisClient(ctx, rdCfg)
	if err != nil {
		logger.Fatalw("failed to establish redis connection", "error", err)
	}

	userRepo := repository.NewUserRepository(pgClient, logger)
	tokenRepo := repository.NewTokenRepository(rdClient, logger)
	jwtManager := jwt.NewJWTManager(jwtCfg)
	authSvc := service.NewAuthService(userRepo, tokenRepo, jwtManager, logger)
	authHandler := grpc.NewAuthHandler(authSvc, validator, logger)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Fatalw("failed to listen", "error", err)
	}

	grpcServer := grpclib.NewServer()
	authpb.RegisterAuthServiceServer(grpcServer, authHandler)
	logger.Infow("starting grpc server", "address", listener.Addr().String())
	if err := grpcServer.Serve(listener); err != nil {
		logger.Fatalw("grpc server failed", "error", err)
	}
}
