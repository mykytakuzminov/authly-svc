package domain

import (
	"fmt"

	"github.com/mykytakuzminov/ridely-svc/shared/errors"
)

var (
	ErrUserNotFound         = fmt.Errorf("user %w", errors.ErrNotFound)
	ErrUserAlreadyExists    = fmt.Errorf("user %w", errors.ErrAlreadyExists)
	ErrTokenNotFound        = fmt.Errorf("token %w", errors.ErrNotFound)
	ErrInvalidSigningMethod = fmt.Errorf("%w signing method", errors.ErrUnauthenticated)
	ErrInvalidToken         = fmt.Errorf("%w token", errors.ErrUnauthenticated)
	ErrInvalidClaim         = fmt.Errorf("%w claim", errors.ErrUnauthenticated)
	ErrInvalidCredentials   = fmt.Errorf("invalid credentials: %w", errors.ErrUnauthenticated)
)
