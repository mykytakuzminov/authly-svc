package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type RegisterInput struct {
	Email    string `json:"email"    validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=12,max=72"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
