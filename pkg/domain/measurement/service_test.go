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

func TestCreateMeasurement_UnsupportedType_ForFloatMeasurement(t *testing.T) {
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
		Type:        "category", // unsupported measurement type
		UserID:      createdParam.UserID,
		ParameterID: createdParam.ID,
		Notes:       "Unsupported measurement type",
		Value:       "high",
		Timestamp:   time.Now().UTC(),
	}

	createdMeasurement, err := measurementService.CreateMeasurement(context.Background(), measurementInput)
	require.Error(t, err)
	assert.Nil(t, createdMeasurement)
	assert.Contains(t, err.Error(), "unsupported measurement type")
}

// func TestUpdateMeasurement_Failure_TypeMismatch_ForFloatMeasurement(t *testing.T) {
//	parameterRepository := parameter.NewInMemoryRepository()
//	parameterService := parameter.NewService(parameterRepository)
//
//	measurementRepository := measurement.NewInMemoryRepository()
//	measurementService := measurement.NewService(measurementRepository, parameterService)
//
//	// Create a Parameter with DataTypeFloat.
//	paramInput := parameter.CreateParameterInput{
//		UserID:      uuid.New(),
//		Name:        "Temperature",
//		Description: "Ambient temperature in Celsius",
//		DataType:    parameter.DataTypeFloat,
//		Unit:        "Celsius",
//	}
//	createdParam, err := parameterService.CreateParameter(context.Background(), paramInput)
//	require.NoError(t, err)
//	require.NotNil(t, createdParam)
//
//	// Create a FloatMeasurement successfully.
//	measurementInput := measurement.CreateMeasurementInput{
//		Type:        measurement.MeasurementTypeFloat,
//		UserID:      createdParam.UserID,
//		ParameterID: createdParam.ID,
//		Value:       25.5,
//		Notes:       "Initial float measurement",
//		Timestamp:   time.Now().UTC(),
//	}
//	createdMeasurement, err := measurementService.CreateMeasurement(context.Background(), measurementInput)
//	require.NoError(t, err)
//	require.NotNil(t, createdMeasurement)
//
//	// Assert that the created measurement is a FloatMeasurement.
//	floatMeas, ok := createdMeasurement.(*measurement.FloatMeasurement)
//	require.True(t, ok, "created measurement should be of type FloatMeasurement")
//
//	// Prepare UpdateMeasurementInput with mismatched value type: using a boolean instead of a float.
//	updateInput := measurement.UpdateMeasurementInput{
//		ID:    floatMeas.ID,
//		Value: true, // Mismatched value: expected a float64, provided a bool.
//	}
//
//	// Act: Attempt to update the measurement.
//	updatedMeasurement, err := measurementService.UpdateMeasurement(context.Background(), updateInput)
//
//	// Assert: Expect an error due to type mismatch.
//	require.Error(t, err)
//	assert.Nil(t, updatedMeasurement)
//	assert.Contains(t, err.Error(), "invalid value type for float measurement")
//}

//func TestUpdateMeasurement_Success_WithConsistentTypes_ForFloatMeasurement(t *testing.T) {
//	// Arrange
//	paramRepo := parameter.NewInMemoryRepository()
//	paramService := parameter.NewService(paramRepo)
//
//	measRepo := measurement.NewInMemoryRepository()
//	measService := measurement.NewService(measRepo, paramService)
//
//	// Create a Parameter with DataTypeFloat
//	paramInput := parameter.CreateParameterInput{
//		UserID:      uuid.New(),
//		Name:        "Temperature",
//		Description: "Ambient temperature in Celsius",
//		DataType:    parameter.DataTypeFloat,
//		Unit:        "Celsius",
//	}
//
//	createdParam, err := paramService.CreateParameter(context.Background(), paramInput)
//	require.NoError(t, err)
//	require.NotNil(t, createdParam)
//
//	// Create an initial FloatMeasurement
//	createMeasInput := measurement.CreateMeasurementInput{
//		Type:        measurement.MeasurementTypeFloat,
//		UserID:      createdParam.UserID,
//		ParameterID: createdParam.ID,
//		Value:       25.5,
//		Notes:       "Initial measurement",
//		Timestamp:   time.Now().UTC(),
//	}
//
//	createdMeas, err := measService.CreateMeasurement(context.Background(), createMeasInput)
//	require.NoError(t, err)
//	require.NotNil(t, createdMeas)
//
//	floatMeas, ok := createdMeas.(*measurement.FloatMeasurement)
//	require.True(t, ok, "expected created measurement to be of type FloatMeasurement")
//
//	// Prepare update input with a new valid float value
//	newValue := 26.0
//	updateInput := UpdateMeasurementInput{
//		ID:    floatMeas.ID,
//		Value: newValue,
//		Notes: "Updated measurement",
//	}
//
//	// Act: Update the measurement
//	updatedMeas, err := measService.UpdateMeasurement(context.Background(), updateInput)
//	require.NoError(t, err)
//	require.NotNil(t, updatedMeas)
//
//	updatedFloatMeas, ok := updatedMeas.(FloatMeasurement)
//	require.True(t, ok, "expected updated measurement to be of type FloatMeasurement")
//
//	// Assert: Verify that the measurement has been updated
//	assert.InEpsilon(t, newValue, updatedFloatMeas.Value, 0.0001, "updated value should match new value")
//	assert.Equal(t, floatMeas.ID, updatedFloatMeas.ID, "measurement ID should remain unchanged")
//	assert.Equal(t, floatMeas.UserID, updatedFloatMeas.UserID, "userID should remain unchanged")
//	assert.Equal(t, floatMeas.ParameterID, updatedFloatMeas.ParameterID, "parameterID should remain unchanged")
//	assert.Equal(t, "Updated measurement", updatedFloatMeas.Notes, "notes should be updated")
//	assert.True(t, updatedFloatMeas.UpdatedAt.After(updatedFloatMeas.CreatedAt), "updatedAt should be after createdAt")
//}

// https://chatgpt.com/c/678a9722-bd3c-800d-b4cb-bb5d2566a59d?model=o1-mini
