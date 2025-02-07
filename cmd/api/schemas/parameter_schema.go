package schemas

import (
	"time"

	"github.com/dim2k2006/correlateapp-be/pkg/domain/parameter"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CreateParameterRequest struct {
	UserID      uuid.UUID          `json:"userId" validate:"required,uuid4"`
	Name        string             `json:"name" validate:"required,min=2,max=100"`
	Description string             `json:"description,omitempty"`
	DataType    parameter.DataType `json:"dataType" validate:"required,oneof=float"`
	Unit        string             `json:"unit,omitempty"`
}

type UpdateParameterRequest struct {
	ID          uuid.UUID `json:"id" validate:"required,uuid4"`
	Name        *string   `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Description *string   `json:"description,omitempty" validate:"omitempty"`
	Unit        *string   `json:"unit,omitempty" validate:"omitempty"`
}

func getParameterRequestValidator() *validator.Validate {
	return validator.New()
}

func (r *CreateParameterRequest) Validate() error {
	return getParameterRequestValidator().Struct(r)
}

func (r *UpdateParameterRequest) Validate() error {
	return getParameterRequestValidator().Struct(r)
}

type ParameterResponse struct {
	ID          uuid.UUID          `json:"id"`
	UserID      uuid.UUID          `json:"userId"`
	Name        string             `json:"name"`
	Description string             `json:"description,omitempty"`
	DataType    parameter.DataType `json:"dataType"`
	Unit        string             `json:"unit,omitempty"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
}

func NewParameterResponse(p *parameter.Parameter) ParameterResponse {
	return ParameterResponse{
		ID:          p.ID,
		UserID:      p.UserID,
		Name:        p.Name,
		Description: p.Description,
		DataType:    p.DataType,
		Unit:        p.Unit,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
