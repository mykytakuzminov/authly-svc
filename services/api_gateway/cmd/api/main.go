package main

import (
	"context"
	"errors"
	"fmt"

	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/mykytakuzminov/ridely-svc/docs"
	"github.com/mykytakuzminov/ridely-svc/services/api_gateway/internal/http"
	"github.com/mykytakuzminov/ridely-svc/services/api_gateway/internal/infra/grpc"
	"github.com/mykytakuzminov/ridely-svc/shared/server"
)

// @title           Ridely
// @version         1.0
//
// @BasePath        /api/v1
//
// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
// @description                Enter your Bearer token in the format: Bearer <token>
func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()

	app, err := newApp(logger)
	if err != nil {
		logger.Fatalw("failed to initialize app", "error", err)
	}

	os.Exit(app.Run())
}

type app struct {
	httpServer *server.HTTPServer
	authClient *grpc.AuthServiceClient
	logger     *zap.SugaredLogger
}

func newApp(logger *zap.SugaredLogger) (*app, error) {
	a := &app{logger: logger}

	authClient, err := grpc.NewAuthServiceClient()
	if err != nil {
		return nil, fmt.Errorf("auth client init: %w", err)
	}
	a.authClient = authClient

	router := a.initRouter()

	httpServer, err := server.NewHTTPServer(router, logger)
	if err != nil {
		a.Close()
		return nil, fmt.Errorf("http server init: %w", err)
	}
	a.httpServer = httpServer

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
	defer a.Close()

	errCh := make(chan error, 1)
	go func() {
		errCh <- a.httpServer.Run()
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

	if err := a.httpServer.GracefulStop(ctx); err != nil {
		exitCode = 1
	}

	return exitCode
}

func (a *app) Close() {
	if a.authClient != nil {
		if err := a.authClient.Close(); err != nil {
			a.logger.Errorw("failed to close auth client connection", "error", err)
		}
	}
}

func (a *app) initRouter() chi.Router {
	authHandler := http.NewAuthHTTPHandler(a.authClient.Client, validator.New(), a.logger)

	router := chi.NewRouter()

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	router.Use(http.TimeoutMiddleware)
	router.Use(http.TraceMiddleware)

	router.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)
	})

	return router
}
