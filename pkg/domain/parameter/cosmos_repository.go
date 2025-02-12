package parameter

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
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
