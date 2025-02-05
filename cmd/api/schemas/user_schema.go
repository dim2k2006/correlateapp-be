package schemas

import "github.com/go-playground/validator/v10"

type CreateUserRequest struct {
	ExternalID string `json:"externalId" validate:"required,uuid4"`
	FirstName  string `json:"firstName" validate:"required,min=2,max=50"`
	LastName   string `json:"lastName" validate:"required,min=2,max=50"`
}

type UpdateUserRequest struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

var validate = validator.New()

func (r *CreateUserRequest) Validate() error {
	return validate.Struct(r)
}

func (r *UpdateUserRequest) Validate() error {
	return validate.Struct(r)
}
