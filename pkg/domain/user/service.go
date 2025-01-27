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
	ExternalID string `json:"external_id" validate:"required"`
	FirstName  string `json:"first_name" validate:"required"`
	LastName   string `json:"last_name" validate:"required"`
}

type UpdateUserInput struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}
