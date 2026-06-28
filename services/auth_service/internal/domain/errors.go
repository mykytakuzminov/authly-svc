package domain

import (
	"fmt"

	"github.com/mykytakuzminov/ridely-svc/shared/errors"
)

var (
	ErrUserNotFound         = fmt.Errorf("user %w", errors.ErrNotFound)
	ErrUserAlreadyExists    = fmt.Errorf("user %w", errors.ErrAlreadyExists)
	ErrTokenNotFound        = fmt.Errorf("token %w", errors.ErrNotFound)
	ErrInvalidSigningMethod = fmt.Errorf("%w signing method", errors.ErrInvalid)
	ErrInvalidToken         = fmt.Errorf("%w token", errors.ErrInvalid)
	ErrInvalidClaim         = fmt.Errorf("%w claim", errors.ErrInvalid)
)
