package grpc

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/mykytakuzminov/ridely-svc/shared/env"
	i "github.com/mykytakuzminov/ridely-svc/shared/interceptors"
	authpb "github.com/mykytakuzminov/ridely-svc/shared/proto/auth"
)

type AuthServiceClient struct {
	Client authpb.AuthServiceClient
	conn   *grpc.ClientConn
}

func NewAuthServiceClient() (*AuthServiceClient, error) {
	addr := env.GetString("AUTH_SERVICE_ADDR", "0.0.0.0:50051")

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(i.TraceClientInterceptor),
	)
	if err != nil {
		return nil, fmt.Errorf("connect auth service: %w", err)
	}

	return &AuthServiceClient{
		Client: authpb.NewAuthServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *AuthServiceClient) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}
