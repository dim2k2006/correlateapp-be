package measurement_test

import (
	"context"
	"testing"

	"github.com/dim2k2006/correlateapp-be/pkg/domain/measurement"
	"github.com/dim2k2006/correlateapp-be/pkg/domain/parameter"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateMeasurement_Success(t *testing.T) {
	parameterRepository := parameter.NewInMemoryRepository()
	parameterService := parameter.NewService(parameterRepository)

	measurementRepository := measurement.NewInMemoryRepository()
	measurementService := measurement.NewService(measurementRepository, parameterService)

	parameterInput := parameter.CreateParameterInput{
		UserID:      uuid.UUID{},
		Name:        "Test Parameter",
		Description: "Test Description",
		DataType:    parameter.DataTypeFloat,
		Unit:        "Test Unit",
	}

	createdParam, err := parameterService.CreateParameter(context.Background(), parameterInput)

	assert.NoError(t, err)
	assert.NotNil(t, createdParam)

	measurementInput := measurement.CreateMeasurementInput{
		Type:        measurement.MeasurementTypeFloat,
		UserID:      uuid.UUID{},
		ParameterID: createdParam.ID,
		Notes:       "Test Notes",
		Value:       25.5,
	}

	createdMeasurement, err := measurementService.CreateMeasurement(context.Background(), measurementInput)

	assert.NoError(t, err)
	assert.NotNil(t, createdMeasurement)

	floatMeas, ok := createdMeasurement.(*measurement.FloatMeasurement)
	assert.True(t, ok)
	assert.Equal(t, 25.5, floatMeas.Value)
	assert.Equal(t, measurement.MeasurementTypeFloat, floatMeas.Type)
	assert.Equal(t, createdParam.ID, floatMeas.ParameterID)
	assert.Equal(t, createdParam.UserID, floatMeas.UserID)
}

// https://chatgpt.com/c/678a9722-bd3c-800d-b4cb-bb5d2566a59d?model=o1-mini
