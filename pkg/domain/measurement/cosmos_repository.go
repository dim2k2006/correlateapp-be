package measurement

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
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
