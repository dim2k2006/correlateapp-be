package schemas

import (
	"time"

	"github.com/dim2k2006/correlateapp-be/pkg/domain/measurement"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CreateMeasurementRequest struct {
	Type        measurement.Type `json:"type" validate:"required,oneof=float"`
	UserID      uuid.UUID        `json:"userId" validate:"required,uuid4"`
	ParameterID uuid.UUID        `json:"parameterId" validate:"required,uuid4"`
	Notes       string           `json:"notes,omitempty" validate:"omitempty"`
	Value       interface{}      `json:"value" validate:"required,min=1"`
}

func getMeasurementRequestValidator() *validator.Validate {
	return validator.New()
}

func (r *CreateMeasurementRequest) Validate() error {
	return getMeasurementRequestValidator().Struct(r)
}

type MeasurementResponse struct {
	ID          uuid.UUID        `json:"id"`
	Type        measurement.Type `json:"type"`
	UserID      uuid.UUID        `json:"userId"`
	ParameterID uuid.UUID        `json:"parameterId"`
	Timestamp   time.Time        `json:"timestamp"`
	Notes       string           `json:"notes,omitempty"`
	CreatedAt   time.Time        `json:"createdAt"`
	UpdatedAt   time.Time        `json:"updatedAt"`
}

func NewMeasurementResponse(m *measurement.FloatMeasurement) MeasurementResponse {
	return MeasurementResponse{
		ID:          m.ID,
		Type:        m.Type,
		UserID:      m.UserID,
		ParameterID: m.ParameterID,
		Timestamp:   m.Timestamp,
		Notes:       m.Notes,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
