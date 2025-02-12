package parameter

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
	containerName = "Parameters"
	partitionKey  = "/userId"
)

type CosmosParameterRepository struct {
	client    *azcosmos.Client
	container *azcosmos.ContainerClient
}

func NewCosmosParameterRepository(connectionString string) (*CosmosParameterRepository, error) {
	client, err := azcosmos.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Cosmos DB client for parameter repository: %w", err)
	}

	container, err := client.NewContainer(databaseName, containerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get Cosmos DB container for parameter repository: %w", err)
	}

	return &CosmosParameterRepository{
		client:    client,
		container: container,
	}, nil
}

func (r *CosmosParameterRepository) CreateParameter(ctx context.Context, parameter *Parameter) (*Parameter, error) {
	parameterJSON, err := json.Marshal(NewCosmosParameter(parameter))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal parameter: %w", err)
	}

	pk := azcosmos.NewPartitionKeyString(parameter.ID.String())

	_, err = r.container.CreateItem(ctx, pk, parameterJSON, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create parameter in Cosmos DB: %w", err)
	}

	return parameter, nil
}

func (r *CosmosParameterRepository) GetParameterByID(ctx context.Context, id uuid.UUID) (*Parameter, error) {
	pk := azcosmos.NewPartitionKeyString(id.String())

	resp, readItemErr := r.container.ReadItem(ctx, pk, id.String(), nil)
	if readItemErr != nil {
		return nil, ErrParameterNotFound
	}

	var cosmosParameter *CosmosParameter
	if err := json.Unmarshal(resp.Value, &cosmosParameter); err != nil {
		return nil, fmt.Errorf("failed to unmarshal parameter: %w", err)
	}

	parameter := NewParameter(cosmosParameter)

	return parameter, nil
}

type CosmosParameter struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"userId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	DataType    DataType  `json:"dataType"`
	Unit        string    `json:"unit"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func NewCosmosParameter(parameter *Parameter) *CosmosParameter {
	return &CosmosParameter{
		ID:          parameter.ID,
		UserID:      parameter.UserID,
		Name:        parameter.Name,
		Description: parameter.Description,
		DataType:    parameter.DataType,
		Unit:        parameter.Unit,
		CreatedAt:   parameter.CreatedAt,
		UpdatedAt:   parameter.UpdatedAt,
	}
}

func NewParameter(cosmosParameter *CosmosParameter) *Parameter {
	return &Parameter{
		ID:          cosmosParameter.ID,
		UserID:      cosmosParameter.UserID,
		Name:        cosmosParameter.Name,
		Description: cosmosParameter.Description,
		DataType:    cosmosParameter.DataType,
		Unit:        cosmosParameter.Unit,
		CreatedAt:   cosmosParameter.CreatedAt,
		UpdatedAt:   cosmosParameter.UpdatedAt,
	}
}
