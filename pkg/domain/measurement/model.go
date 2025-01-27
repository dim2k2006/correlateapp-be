package measurement

import (
	"encoding/json"
	"errors"
	"fmt"
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

type BooleanMeasurement struct {
	BaseMeasurement
	Value bool `json:"value"` // Boolean-specific field
}

type Measurement interface{}

type Wrapper struct {
	Measurement Measurement `json:"-"`
}

func (mw Wrapper) MarshalJSON() ([]byte, error) {
	switch m := mw.Measurement.(type) {
	case FloatMeasurement:
		return json.Marshal(m)
	case BooleanMeasurement:
		return json.Marshal(m)
	default:
		return nil, fmt.Errorf("unsupported measurement type: %T", m)
	}
}

func (mw *Wrapper) UnmarshalJSON(data []byte) error {
	var base struct {
		Type Type `json:"type"`
	}
	if err := json.Unmarshal(data, &base); err != nil {
		return err
	}

	switch base.Type {
	case MeasurementTypeFloat:
		var fm FloatMeasurement
		if err := json.Unmarshal(data, &fm); err != nil {
			return err
		}
		mw.Measurement = fm
	// case MeasurementTypeBoolean:
	//	var bm BooleanMeasurement
	//	if err := json.Unmarshal(data, &bm); err != nil {
	//		return err
	//	}
	//	mw.Measurement = bm
	default:
		return errors.New("unknown measurement type: " + string(base.Type))
	}
	return nil
}
