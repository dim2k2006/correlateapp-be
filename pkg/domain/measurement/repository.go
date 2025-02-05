package measurement

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	CreateMeasurement(ctx context.Context, measurement Measurement) (Measurement, error)
	ListMeasurementsByUser(ctx context.Context, userID uuid.UUID) ([]Measurement, error)
	ListMeasurementsByParameter(ctx context.Context, parameterID uuid.UUID) ([]Measurement, error)
	DeleteMeasurement(ctx context.Context, id uuid.UUID) error
}
