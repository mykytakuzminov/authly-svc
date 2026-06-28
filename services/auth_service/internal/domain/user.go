package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
	Email    string
	Password string
}

type DB interface {
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) pgx.Rows
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type UserService interface {
	Register(ctx context.Context, input *RegisterInput) (*TokenPair, error)
}
