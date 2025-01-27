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
	UserID      uuid.UUID `json:"userId"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	DataType    DataType  `json:"dataType"`
	Unit        string    `json:"unit,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
