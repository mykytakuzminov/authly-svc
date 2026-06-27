package domain

import (
	"fmt"

	"github.com/mykytakuzminov/ridely-svc/shared/errors"
)

var (
	ErrUserNotFound = fmt.Errorf("user %w", errors.ErrNotFound)
)
