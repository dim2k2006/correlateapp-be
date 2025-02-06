package schemas

import (
	"time"

	"github.com/dim2k2006/correlateapp-be/pkg/domain/user"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CreateUserRequest struct {
	ExternalID string `json:"externalId" validate:"required,uuid4"`
	FirstName  string `json:"firstName" validate:"required,min=2,max=50"`
	LastName   string `json:"lastName" validate:"required,min=2,max=50"`
}

type UpdateUserRequest struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

func getValidator() *validator.Validate {
	return validator.New()
}

func (r *CreateUserRequest) Validate() error {
	return getValidator().Struct(r)
}

func (r *UpdateUserRequest) Validate() error {
	return getValidator().Struct(r)
}

type UserResponse struct {
	ID         uuid.UUID `json:"id"`
	ExternalID string    `json:"externalId"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func NewUserResponse(u *user.User) UserResponse {
	return UserResponse{
		ID:         u.ID,
		ExternalID: u.ExternalID,
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}
