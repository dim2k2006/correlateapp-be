package user

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByExternalID(ctx context.Context, externalID string) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
