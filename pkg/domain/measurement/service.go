package measurement

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	CreateMeasurement(ctx context.Context, input CreateMeasurementInput) (Measurement, error)
	ListMeasurementsByUser(ctx context.Context, userID uuid.UUID) ([]Measurement, error)
	ListMeasurementsByParameter(ctx context.Context, parameterID uuid.UUID) ([]Measurement, error)
	DeleteMeasurement(ctx context.Context, id uuid.UUID) error
}

type CreateMeasurementInput struct {
	ParameterID uuid.UUID
	Notes       string
	Value       interface{}
	Timestamp   time.Time
}
