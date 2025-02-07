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
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        string
	Description string
	DataType    DataType
	Unit        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
