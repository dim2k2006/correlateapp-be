package parameter

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	CreateParameter(ctx context.Context, input CreateParameterInput) (*Parameter, error)
	GetParameterByID(ctx context.Context, id uuid.UUID) (*Parameter, error)
	ListParametersByUser(ctx context.Context, userID uuid.UUID) ([]*Parameter, error)
	UpdateParameter(ctx context.Context, input UpdateParameterInput) (*Parameter, error)
	DeleteParameter(ctx context.Context, id uuid.UUID) error
}

type CreateParameterInput struct {
	UserID      uuid.UUID `json:"user_id" validate:"required"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description,omitempty"`
	DataType    DataType  `json:"data_type" validate:"required,oneof=float"`
	Unit        string    `json:"unit,omitempty"`
}

type UpdateParameterInput struct {
	Name        string   `json:"name,omitempty" validate:"omitempty"`
	Description string   `json:"description,omitempty"`
	DataType    DataType `json:"data_type,omitempty" validate:"omitempty,oneof=float"`
	Unit        string   `json:"unit,omitempty"`
}
