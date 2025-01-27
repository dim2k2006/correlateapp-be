package measurement

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	CreateMeasurement(ctx context.Context, measurement *Measurement) (*Measurement, error)
	GetMeasurementByID(ctx context.Context, id uuid.UUID) (*Measurement, error)
	ListMeasurementsByUser(ctx context.Context, userID uuid.UUID) ([]*Measurement, error)
	ListMeasurementsByParameter(ctx context.Context, parameterID uuid.UUID) ([]*Measurement, error)
	UpdateMeasurement(ctx context.Context, measurement *Measurement) (*Measurement, error)
	DeleteMeasurement(ctx context.Context, id uuid.UUID) error
}
