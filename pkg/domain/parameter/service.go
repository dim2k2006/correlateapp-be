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
	UserID      uuid.UUID
	Name        string
	Description string
	DataType    DataType
	Unit        string
}

type UpdateParameterInput struct {
	ID          uuid.UUID
	Name        *string
	Description *string
	Unit        *string
}
