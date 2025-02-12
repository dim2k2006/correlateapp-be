package parameter

import (
	"context"
	"encoding/json"
	"errors"
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
	query := "SELECT * FROM parameters p WHERE p.id = @id"
	params := []azcosmos.QueryParameter{
		{Name: "@id", Value: id.String()},
	}

	queryOptions := &azcosmos.QueryOptions{QueryParameters: params}
	pager := r.container.NewQueryItemsPager(query, azcosmos.NewPartitionKey(), queryOptions)

	var cosmosParameter CosmosParameter
	for pager.More() {
		resp, nextPageErr := pager.NextPage(ctx)
		if nextPageErr != nil {
			return nil, fmt.Errorf("query failed: %w", nextPageErr)
		}

		if len(resp.Items) > 0 {
			if err := json.Unmarshal(resp.Items[0], &cosmosParameter); err != nil {
				return nil, fmt.Errorf("failed to unmarshal parameter: %w", err)
			}

			return NewParameter(&cosmosParameter), nil
		}
	}

	return nil, errors.New("parameter not found")
}

func (r *CosmosParameterRepository) ListParametersByUser(ctx context.Context, userID uuid.UUID) ([]*Parameter, error) {
	query := "SELECT * FROM parameters p WHERE p.userId = @userID"
	params := []azcosmos.QueryParameter{
		{Name: "@userID", Value: userID.String()},
	}

	queryOptions := &azcosmos.QueryOptions{QueryParameters: params}
	pager := r.container.NewQueryItemsPager(query, azcosmos.NewPartitionKey(), queryOptions)

	var parameters []*Parameter
	for pager.More() {
		resp, nextPageErr := pager.NextPage(ctx)
		if nextPageErr != nil {
			return nil, fmt.Errorf("query failed: %w", nextPageErr)
		}

		for _, item := range resp.Items {
			var cosmosParameter CosmosParameter
			if err := json.Unmarshal(item, &cosmosParameter); err != nil {
				return nil, fmt.Errorf("failed to unmarshal parameter: %w", err)
			}
			parameters = append(parameters, NewParameter(&cosmosParameter))
		}
	}

	return parameters, nil
}

func (r *CosmosParameterRepository) UpdateParameter(ctx context.Context, param *Parameter) (*Parameter, error) {
	paramJSON, err := json.Marshal(NewCosmosParameter(param))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal parameter: %w", err)
	}

	pk := azcosmos.NewPartitionKeyString(param.UserID.String())

	_, err = r.container.ReplaceItem(ctx, pk, param.ID.String(), paramJSON, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to update parameter in Cosmos DB: %w", err)
	}

	return param, nil
}

func (r *CosmosParameterRepository) DeleteParameter(ctx context.Context, id uuid.UUID) error {
	// First, retrieve the parameter to get its UserID (required for partition key)
	param, err := r.GetParameterByID(ctx, id)
	if err != nil {
		return err
	}

	pk := azcosmos.NewPartitionKeyString(param.UserID.String())

	_, err = r.container.DeleteItem(ctx, pk, id.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to delete parameter from Cosmos DB: %w", err)
	}

	return nil
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
