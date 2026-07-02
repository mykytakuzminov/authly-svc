package grpc

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sharedErrors "github.com/mykytakuzminov/ridely-svc/shared/errors"
)

func toGRPCError(err error) error {
	switch {
	case errors.Is(err, sharedErrors.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, sharedErrors.ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, sharedErrors.ErrInvalid):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, sharedErrors.ErrIsRequired):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
