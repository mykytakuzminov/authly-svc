package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	grpclib "google.golang.org/grpc"

	"github.com/mykytakuzminov/ridely-svc/services/auth_service/internal/infra/grpc"
	"github.com/mykytakuzminov/ridely-svc/services/auth_service/internal/infra/jwt"
	"github.com/mykytakuzminov/ridely-svc/services/auth_service/internal/infra/migrations"
	"github.com/mykytakuzminov/ridely-svc/services/auth_service/internal/infra/repository"
	"github.com/mykytakuzminov/ridely-svc/services/auth_service/internal/service"
	"github.com/mykytakuzminov/ridely-svc/shared/db"
	i "github.com/mykytakuzminov/ridely-svc/shared/interceptors"
	"github.com/mykytakuzminov/ridely-svc/shared/migrate"
	authpb "github.com/mykytakuzminov/ridely-svc/shared/proto/auth"
	"github.com/mykytakuzminov/ridely-svc/shared/server"
)

func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	app, err := newApp(ctx, logger)
	if err != nil {
		logger.Fatalw("failed to initialize app", "error", err)
	}

	os.Exit(app.Run())
}

type app struct {
	pgClient   *pgxpool.Pool
	rdClient   *redis.Client
	grpcServer *server.GRPCServer
	logger     *zap.SugaredLogger
}

func newApp(ctx context.Context, logger *zap.SugaredLogger) (*app, error) {
	a := &app{logger: logger}

	if err := godotenv.Load(); err != nil {
		logger.Warnw("no .env file, reading from environment")
	}

	pgClient, err := db.ConnectPostgres(ctx, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("postgres init: %w", err)
	}
	a.pgClient = pgClient

	if err := migrate.Run(pgClient, migrations.FS, logger); err != nil {
		a.close()
		return nil, fmt.Errorf("migrate: %w", err)
	}

	rdClient, err := db.ConnectRedis(ctx, 5*time.Second)
	if err != nil {
		a.close()
		return nil, fmt.Errorf("redis init: %w", err)
	}
	a.rdClient = rdClient

	jwtManager, err := jwt.NewManager(logger)
	if err != nil {
		a.close()
		return nil, fmt.Errorf("jwt manager init: %w", err)
	}

	userRepo := repository.NewUserRepository(pgClient, logger)
	tokenRepo := repository.NewTokenRepository(rdClient, logger)
	authSvc := service.NewAuthService(userRepo, tokenRepo, jwtManager, logger)
	authHandler := grpc.NewAuthHandler(authSvc, validator.New(), logger)

	grpcServer, err := server.NewGRPCServer(logger,
		server.WithRegisterFn(func(s *grpclib.Server) {
			authpb.RegisterAuthServiceServer(s, authHandler)
		}),
		server.WithServerOption(grpclib.ChainUnaryInterceptor(
			i.TraceServerInterceptor,
			i.UserContextServerInterceptor,
		)),
	)
	if err != nil {
		a.close()
		return nil, fmt.Errorf("grpc server init: %w", err)
	}
	a.grpcServer = grpcServer

	return a, nil
}

func (a *app) Run() int {
	defer func() {
		if err := a.logger.Sync(); err != nil &&
			!errors.Is(err, syscall.ENOTTY) &&
			!errors.Is(err, syscall.EINVAL) {
			a.logger.Errorw("failed to sync logger", "error", err)
		}
	}()
	defer a.close()

	errCh := make(chan error, 1)
	go func() {
		errCh <- a.grpcServer.Run()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	exitCode := 0
	select {
	case err := <-errCh:
		a.logger.Errorw("unexpected error occurred", "error", err)
		exitCode = 1
	case sig := <-quit:
		a.logger.Infow("shutdown signal received", "signal", sig.String())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	a.grpcServer.GracefulStop(ctx)

	return exitCode
}

func (a *app) close() {
	if a.pgClient != nil {
		a.pgClient.Close()
	}
	if a.rdClient != nil {
		if err := a.rdClient.Close(); err != nil {
			a.logger.Errorw("failed to close redis connection", "error", err)
		}
	}
}
