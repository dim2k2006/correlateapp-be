package measurement

import (
	"time"

	"github.com/google/uuid"
)

type Type string

const (
	MeasurementTypeFloat Type = "float"
	// MeasurementTypeBoolean MeasurementType = "boolean".
)

type BaseMeasurement struct {
	Type        Type      `json:"type"`
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"userId"`
	ParameterID uuid.UUID `json:"parameterId"`
	Timestamp   time.Time `json:"timestamp"`
	Notes       string    `json:"notes,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type FloatMeasurement struct {
	BaseMeasurement
	Value float64 `json:"value"`
}

func (fm *FloatMeasurement) GetID() uuid.UUID {
	return fm.ID
}

func (fm *FloatMeasurement) GetUserID() uuid.UUID {
	return fm.UserID
}

func (fm *FloatMeasurement) GetParameterID() uuid.UUID {
	return fm.ParameterID
}

func (fm *FloatMeasurement) GetType() Type {
	return fm.Type
}

func (fm *FloatMeasurement) GetTimestamp() time.Time {
	return fm.Timestamp
}

func (fm *FloatMeasurement) GetNotes() string {
	return fm.Notes
}

func (fm *FloatMeasurement) SetID(id uuid.UUID) {
	fm.ID = id
}

func (fm *FloatMeasurement) SetCreatedAt(t time.Time) {
	fm.CreatedAt = t
}

func (fm *FloatMeasurement) SetUpdatedAt(t time.Time) {
	fm.UpdatedAt = t
}

type BooleanMeasurement struct {
	BaseMeasurement
	Value bool `json:"value"` // Boolean-specific field
}

func (bm *BooleanMeasurement) GetID() uuid.UUID {
	return bm.ID
}

func (bm *BooleanMeasurement) GetUserID() uuid.UUID {
	return bm.UserID
}

func (bm *BooleanMeasurement) GetParameterID() uuid.UUID {
	return bm.ParameterID
}

func (bm *BooleanMeasurement) GetType() Type {
	return bm.Type
}

func (bm *BooleanMeasurement) GetTimestamp() time.Time {
	return bm.Timestamp
}

func (bm *BooleanMeasurement) GetNotes() string {
	return bm.Notes
}

func (bm *BooleanMeasurement) SetID(id uuid.UUID) {
	bm.ID = id
}

func (bm *BooleanMeasurement) SetCreatedAt(t time.Time) {
	bm.CreatedAt = t
}

func (bm *BooleanMeasurement) SetUpdatedAt(t time.Time) {
	bm.UpdatedAt = t
}

type Measurement interface {
	GetID() uuid.UUID
	GetUserID() uuid.UUID
	GetParameterID() uuid.UUID
	GetType() Type
	GetTimestamp() time.Time
	GetNotes() string
	SetID(uuid.UUID)
	SetCreatedAt(time.Time)
	SetUpdatedAt(time.Time)
}
