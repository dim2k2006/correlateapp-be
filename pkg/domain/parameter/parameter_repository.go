package parameter

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	CreateParameter(ctx context.Context, parameter *Parameter) (*Parameter, error)
	GetParameterByID(ctx context.Context, id uuid.UUID) (*Parameter, error)
	ListParametersByUser(ctx context.Context, userID uuid.UUID) ([]*Parameter, error)
	UpdateParameter(ctx context.Context, param *Parameter) (*Parameter, error)
	DeleteParameter(ctx context.Context, id uuid.UUID) error
}
