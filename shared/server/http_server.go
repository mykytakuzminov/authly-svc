package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/mykytakuzminov/ridely-svc/shared/env"
)

type httpServerConfig struct {
	Host string
	Port string
}

func newHTTPServerConfig() (*httpServerConfig, error) {
	cfg := &httpServerConfig{}
	cfg.load()
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *httpServerConfig) load() {
	c.Host = env.GetString("HTTP_SERVER_HOST", "0.0.0.0")
	c.Port = env.GetString("HTTP_SERVER_PORT", "8080")
}

func (c *httpServerConfig) validate() error {
	return nil
}

func (c *httpServerConfig) getAddr() string {
	return c.Host + ":" + c.Port
}

type HTTPServer struct {
	server *http.Server
	logger *zap.SugaredLogger
}

func NewHTTPServer(router http.Handler, logger *zap.SugaredLogger) (*HTTPServer, error) {
	cfg, err := newHTTPServerConfig()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	httpServer := &http.Server{
		Addr:    cfg.getAddr(),
		Handler: router,
	}

	return &HTTPServer{
		server: httpServer,
		logger: logger.With("component", "http_server"),
	}, nil
}

func (s *HTTPServer) Run() error {
	s.logger.Infow("starting http server", "address", s.server.Addr)

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *HTTPServer) GracefulStop(ctx context.Context) error {
	s.logger.Infow("shutting down http server")

	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Warnw("graceful shutdown timed out, forcing stop", "error", err)
		return s.server.Close()
	}

	s.logger.Infow("http server stopped gracefully")
	return nil
}
