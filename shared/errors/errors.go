package errors

import "errors"

var (
	ErrNotFound   = errors.New("not found")
	ErrIsRequired = errors.New("is required")
)
