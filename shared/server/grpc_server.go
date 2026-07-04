package server

import (
	"context"
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/mykytakuzminov/ridely-svc/shared/env"
)

type grpcServerConfig struct {
	Network string
	Host    string
	Port    string
}

func newGRPCServerConfig() (*grpcServerConfig, error) {
	cfg := &grpcServerConfig{}
	cfg.load()
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *grpcServerConfig) load() {
	c.Network = env.GetString("GRPC_SERVER_NETWORK", "tcp")
	c.Host = env.GetString("GRPC_SERVER_HOST", "0.0.0.0")
	c.Port = env.GetString("GRPC_SERVER_PORT", "50051")
}

func (c *grpcServerConfig) validate() error {
	return nil
}

func (c *grpcServerConfig) getAddr() string {
	return c.Host + ":" + c.Port
}

type GRPCServer struct {
	server   *grpc.Server
	listener net.Listener
	logger   *zap.SugaredLogger
}

type RegisterFn func(*grpc.Server)

func NewGRPCServer(logger *zap.SugaredLogger, registerFns ...RegisterFn) (*GRPCServer, error) {
	cfg, err := newGRPCServerConfig()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	listener, err := net.Listen(cfg.Network, cfg.getAddr())
	if err != nil {
		return nil, fmt.Errorf("listen: %w", err)
	}

	grpcServer := grpc.NewServer()

	for _, register := range registerFns {
		register(grpcServer)
	}

	return &GRPCServer{
		server:   grpcServer,
		listener: listener,
		logger:   logger.With("component", "grpc_server"),
	}, nil
}

func (s *GRPCServer) Run() error {
	s.logger.Infow("starting grpc server", "address", s.listener.Addr().String())
	return s.server.Serve(s.listener)
}

func (s *GRPCServer) GracefulStop(ctx context.Context) {
	s.logger.Infow("shutting down grpc server")

	done := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		s.logger.Infow("grpc server stopped gracefully")
	case <-ctx.Done():
		s.logger.Warnw("graceful shutdown timed out, forcing stop", "error", ctx.Err())
		s.server.Stop()
	}
}
