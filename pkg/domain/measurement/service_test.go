package measurement_test

import (
	"context"
	"testing"
	"time"

	"github.com/dim2k2006/correlateapp-be/pkg/domain/measurement"
	"github.com/dim2k2006/correlateapp-be/pkg/domain/parameter"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	require.NoError(t, err)
	assert.NotNil(t, createdParam)

	measurementInput := measurement.CreateMeasurementInput{
		Type:        measurement.MeasurementTypeFloat,
		UserID:      uuid.UUID{},
		ParameterID: createdParam.ID,
		Notes:       "Test Notes",
		Value:       25.5,
	}

	createdMeasurement, err := measurementService.CreateMeasurement(context.Background(), measurementInput)

	require.NoError(t, err)
	assert.NotNil(t, createdMeasurement)

	floatMeas, ok := createdMeasurement.(*measurement.FloatMeasurement)
	assert.True(t, ok)
	assert.InEpsilon(t, 25.5, floatMeas.Value, 0.0001)
	assert.Equal(t, measurement.MeasurementTypeFloat, floatMeas.Type)
	assert.Equal(t, createdParam.ID, floatMeas.ParameterID)
	assert.Equal(t, createdParam.UserID, floatMeas.UserID)
}

func TestCreateMeasurement_InvalidValueType_ForFloatMeasurement(t *testing.T) {
	parameterRepository := parameter.NewInMemoryRepository()
	parameterService := parameter.NewService(parameterRepository)

	measurementRepository := measurement.NewInMemoryRepository()
	measurementService := measurement.NewService(measurementRepository, parameterService)

	parameterInput := parameter.CreateParameterInput{
		UserID:      uuid.New(),
		Name:        "Temperature",
		Description: "Ambient temperature in Celsius",
		DataType:    parameter.DataTypeFloat,
		Unit:        "Celsius",
	}

	createdParam, err := parameterService.CreateParameter(context.Background(), parameterInput)
	require.NoError(t, err)
	require.NotNil(t, createdParam)

	measurementInput := measurement.CreateMeasurementInput{
		Type:        measurement.MeasurementTypeFloat,
		UserID:      createdParam.UserID,
		ParameterID: createdParam.ID,
		Notes:       "Invalid value type",
		Value:       "not a float", // Invalid value type for float measurement
		Timestamp:   time.Now().UTC(),
	}

	createdMeasurement, err := measurementService.CreateMeasurement(context.Background(), measurementInput)
	require.Error(t, err)
	assert.Nil(t, createdMeasurement)
	assert.Contains(t, err.Error(), "invalid value type for float measurement")
}

func TestCreateMeasurement_ParameterNotFound_ForFloatMeasurement(t *testing.T) {
	parameterRepository := parameter.NewInMemoryRepository()
	parameterService := parameter.NewService(parameterRepository)

	measurementRepository := measurement.NewInMemoryRepository()
	measurementService := measurement.NewService(measurementRepository, parameterService)

	// Use a random ParameterID that does not exist in the repository
	nonExistentParamID := uuid.New()

	measurementInput := measurement.CreateMeasurementInput{
		Type:        measurement.MeasurementTypeFloat,
		UserID:      uuid.New(),
		ParameterID: nonExistentParamID,
		Notes:       "Parameter does not exist",
		Value:       25.0,
		Timestamp:   time.Now().UTC(),
	}

	createdMeasurement, err := measurementService.CreateMeasurement(context.Background(), measurementInput)

	require.Error(t, err)
	assert.Nil(t, createdMeasurement)
	assert.Contains(t, err.Error(), "parameter not found")
}

// https://chatgpt.com/c/678a9722-bd3c-800d-b4cb-bb5d2566a59d?model=o1-mini
