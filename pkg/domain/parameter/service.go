package parameter

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	CreateParameter(ctx context.Context, input *CreateParameterInput) (*Parameter, error)
	GetParameterByID(ctx context.Context, id uuid.UUID) (*Parameter, error)
	ListParametersByUser(ctx context.Context, userID uuid.UUID) ([]*Parameter, error)
	UpdateParameter(ctx context.Context, param *Parameter) (*Parameter, error)
	DeleteParameter(ctx context.Context, id uuid.UUID) error
}

type CreateParameterInput struct {
	UserID      uuid.UUID `json:"user_id" validate:"required"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description,omitempty"`
	DataType    DataType  `json:"data_type" validate:"required,oneof=float boolean category"`
	Unit        string    `json:"unit,omitempty"`
}
