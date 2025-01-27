package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"id"`
	ExternalID string    `json:"external_id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
