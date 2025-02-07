package schemas

import (
	"time"

	"github.com/dim2k2006/correlateapp-be/pkg/domain/measurement"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CreateMeasurementRequest struct {
	ParameterID uuid.UUID   `json:"parameterId" validate:"required,uuid4"`
	Notes       string      `json:"notes,omitempty" validate:"omitempty"`
	Value       interface{} `json:"value" validate:"required,min=1"`
}

func getMeasurementRequestValidator() *validator.Validate {
	return validator.New()
}

func (r *CreateMeasurementRequest) Validate() error {
	return getMeasurementRequestValidator().Struct(r)
}

type MeasurementResponse struct {
	ID          uuid.UUID            `json:"id"`
	Type        measurement.DataType `json:"type"`
	UserID      uuid.UUID            `json:"userId"`
	ParameterID uuid.UUID            `json:"parameterId"`
	Timestamp   time.Time            `json:"timestamp"`
	Notes       string               `json:"notes,omitempty"`
	CreatedAt   time.Time            `json:"createdAt"`
	UpdatedAt   time.Time            `json:"updatedAt"`
}

func NewMeasurementResponse(m measurement.Measurement) MeasurementResponse {
	return MeasurementResponse{
		ID:          m.GetID(),
		Type:        m.GetType(),
		UserID:      m.GetUserID(),
		ParameterID: m.GetParameterID(),
		Timestamp:   m.GetTimestamp(),
		Notes:       m.GetNotes(),
		CreatedAt:   m.GetCreatedAt(),
		UpdatedAt:   m.GetUpdatedAt(),
	}
}
