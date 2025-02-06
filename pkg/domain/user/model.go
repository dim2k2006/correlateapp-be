package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID
	ExternalID string
	FirstName  string
	LastName   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
