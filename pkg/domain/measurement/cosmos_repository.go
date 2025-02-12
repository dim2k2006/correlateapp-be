package measurement

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
)

const (
	databaseName  = "correlateapp"
	containerName = "Measurements"
	partitionKey  = "/parameterId"
)

type CosmosMeasurementRepository struct {
	client    *azcosmos.Client
	container *azcosmos.ContainerClient
}

func NewCosmosMeasurementRepository(connectionString string) (*CosmosMeasurementRepository, error) {
	client, err := azcosmos.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Cosmos DB client for parameter repository: %w", err)
	}

	container, err := client.NewContainer(databaseName, containerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get Cosmos DB container for parameter repository: %w", err)
	}

	return &CosmosMeasurementRepository{
		client:    client,
		container: container,
	}, nil
}

func (r *CosmosMeasurementRepository) CreateMeasurement(
	ctx context.Context,
	measurement Measurement,
) (Measurement, error) {
	measurementJSON, err := json.Marshal(NewCosmosMeasurement(measurement))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal measurement: %w", err)
	}

	pk := azcosmos.NewPartitionKeyString(measurement.GetParameterID().String())

	_, err = r.container.CreateItem(ctx, pk, measurementJSON, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create measurement in Cosmos DB: %w", err)
	}

	return measurement, nil
}

type CosmosMeasurement struct {
	Type        DataType    `json:"type"`
	ID          uuid.UUID   `json:"id"`
	UserID      uuid.UUID   `json:"userId"`
	ParameterID uuid.UUID   `json:"parameterId"`
	Timestamp   time.Time   `json:"timestamp"`
	Notes       string      `json:"notes"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
	Value       interface{} `json:"value"`
}

func NewCosmosMeasurement(m Measurement) *CosmosMeasurement {
	switch m.GetType() {
	case DataTypeFloat:
		return &CosmosMeasurement{
			Type:        m.GetType(),
			ID:          m.GetID(),
			UserID:      m.GetUserID(),
			ParameterID: m.GetParameterID(),
			Timestamp:   m.GetTimestamp(),
			Notes:       m.GetNotes(),
			CreatedAt:   m.GetCreatedAt(),
			UpdatedAt:   m.GetUpdatedAt(),
			Value:       m.(*FloatMeasurement).Value,
		}
	default:
		return nil
	}
}

func NewMeasurement(m *CosmosMeasurement) Measurement {
	switch m.Type {
	case DataTypeFloat:
		return &FloatMeasurement{
			BaseMeasurement: BaseMeasurement{
				Type:        m.Type,
				ID:          m.ID,
				UserID:      m.UserID,
				ParameterID: m.ParameterID,
				Timestamp:   m.Timestamp,
				Notes:       m.Notes,
				CreatedAt:   m.CreatedAt,
				UpdatedAt:   m.UpdatedAt,
			},
			Value: m.Value.(float64),
		}
	default:
		return nil
	}
}
