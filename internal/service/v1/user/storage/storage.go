package storage

import (
	"context"

	"github.com/cockroachdb/errors"
)

var ErrNotFound = errors.New("not found")

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Age   int    `json:"age"`
}

type UserRepo interface {
	CreateUser(ctx context.Context, user User) (int64, error)
	GetUser(ctx context.Context, userID int64) (*User, error)
	GetUsers(ctx context.Context) ([]User, error)
	UpdateUser(ctx context.Context, user User) error
	DeleteUser(ctx context.Context, userID int64) error
}
