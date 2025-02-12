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

func (r *CosmosMeasurementRepository) ListMeasurementsByUser(
	ctx context.Context,
	userID uuid.UUID,
) ([]Measurement, error) {
	query := "SELECT * FROM measurements m WHERE m.userId = @userID"
	params := []azcosmos.QueryParameter{
		{Name: "@userID", Value: userID.String()},
	}

	queryOptions := &azcosmos.QueryOptions{QueryParameters: params}
	pager := r.container.NewQueryItemsPager(query, azcosmos.NewPartitionKey(), queryOptions)

	measurements := []Measurement{}
	for pager.More() {
		resp, nextPageErr := pager.NextPage(ctx)
		if nextPageErr != nil {
			return nil, fmt.Errorf("query failed: %w", nextPageErr)
		}

		for _, item := range resp.Items {
			var cosmosMeasurement CosmosMeasurement
			if err := json.Unmarshal(item, &cosmosMeasurement); err != nil {
				return nil, fmt.Errorf("failed to unmarshal measurement: %w", err)
			}
			measurements = append(measurements, NewMeasurement(&cosmosMeasurement))
		}
	}

	return measurements, nil
}

//func (r *CosmosMeasurementRepository) ListMeasurementsByParameter(ctx context.Context, parameterID uuid.UUID) ([]Measurement, error) {
//	query := "SELECT * FROM measurements m WHERE m.parameterId = @parameterID"
//	params := []azcosmos.QueryParameter{
//		{Name: "@parameterID", Value: parameterID.String()},
//	}
//
//	queryOptions := &azcosmos.QueryOptions{Parameters: params}
//	pager := r.container.NewQueryItemsPager(query, queryOptions)
//
//	measurements := []Measurement{}
//	for pager.More() {
//		resp, err := pager.NextPage(ctx)
//		if err != nil {
//			return nil, fmt.Errorf("query failed: %w", err)
//		}
//
//		for _, item := range resp.Items {
//			var measurement Measurement
//			if err := json.Unmarshal(item, &measurement); err != nil {
//				return nil, fmt.Errorf("failed to unmarshal measurement: %w", err)
//			}
//			measurements = append(measurements, measurement)
//		}
//	}
//
//	return measurements, nil
//}

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
