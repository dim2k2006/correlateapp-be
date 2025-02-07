package user

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	CreateUser(ctx context.Context, input CreateUserInput) (*User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByExternalID(ctx context.Context, externalID string) (*User, error)
	UpdateUser(ctx context.Context, input UpdateUserInput) (*User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type CreateUserInput struct {
	ExternalID string
	FirstName  string
	LastName   string
}

type UpdateUserInput struct {
	ID        uuid.UUID
	FirstName *string
	LastName  *string
}
