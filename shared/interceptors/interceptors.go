package interceptors

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	c "github.com/mykytakuzminov/ridely-svc/shared/context"
)

func TraceClientInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	traceID := c.GetTraceID(ctx)
	ctx = metadata.AppendToOutgoingContext(ctx, "trace_id", traceID.String())
	return invoker(ctx, method, req, reply, cc, opts...)
}

func TraceServerInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return handler(ctx, req)
	}

	values := md.Get("trace_id")
	if len(values) == 0 {
		return handler(ctx, req)
	}

	traceID, err := uuid.Parse(values[0])
	if err != nil {
		return handler(ctx, req)
	}

	ctx = context.WithValue(ctx, c.TraceIDKey, traceID)
	return handler(ctx, req)
}

func UserContextClientInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	if userID, ok := c.GetUserID(ctx); ok {
		ctx = metadata.AppendToOutgoingContext(ctx, "user_id", userID.String())
	}

	return invoker(ctx, method, req, reply, cc, opts...)
}

func UserContextServerInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return handler(ctx, req)
	}

	if values := md.Get("user_id"); len(values) > 0 {
		if userID, err := uuid.Parse(values[0]); err == nil {
			ctx = context.WithValue(ctx, c.UserIDKey, userID)
		}
	}

	return handler(ctx, req)
}
