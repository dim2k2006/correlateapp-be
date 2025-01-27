package parameter

import (
	"time"

	"github.com/google/uuid"
)

type DataType string

const (
	DataTypeFloat DataType = "float"
	// DataTypeBoolean DataType = "boolean"
	// DataTypeCategory DataType = "category"
	// Future data types can be added here.
)

type Parameter struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	DataType    DataType  `json:"data_type"`
	Unit        string    `json:"unit,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
