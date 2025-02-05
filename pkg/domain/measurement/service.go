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
	Type        Type        `json:"type"`
	UserID      uuid.UUID   `json:"userId"`
	ParameterID uuid.UUID   `json:"parameterId"`
	Timestamp   time.Time   `json:"timestamp,omitempty"`
	Notes       string      `json:"notes,omitempty"`
	Value       interface{} `json:"value"`
}
