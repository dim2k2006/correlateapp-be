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
	ExternalID string `json:"externalId" validate:"required"`
	FirstName  string `json:"firstName" validate:"required"`
	LastName   string `json:"lastName" validate:"required"`
}

type UpdateUserInput struct {
	ID        uuid.UUID `json:"id" validate:"required"`
	FirstName string    `json:"firstName,omitempty"`
	LastName  string    `json:"lastName,omitempty"`
}
