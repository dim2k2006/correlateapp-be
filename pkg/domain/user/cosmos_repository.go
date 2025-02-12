package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
)

const (
	databaseName  = "correlateapp"
	containerName = "Users"
	partitionKey  = "/id"
)

type CosmosUserRepository struct {
	client    *azcosmos.Client
	container *azcosmos.ContainerClient
}

func NewCosmosUserRepository(connectionString string) (*CosmosUserRepository, error) {
	client, err := azcosmos.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Cosmos DB client for user repository: %w", err)
	}

	container, err := client.NewContainer(databaseName, containerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get Cosmos DB container for user repository: %w", err)
	}

	return &CosmosUserRepository{
		client:    client,
		container: container,
	}, nil
}

func (r *CosmosUserRepository) CreateUser(ctx context.Context, user *User) (*User, error) {
	userJSON, err := json.Marshal(NewCosmosUser(user))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user: %w", err)
	}

	pk := azcosmos.NewPartitionKeyString(user.ID.String())

	_, err = r.container.CreateItem(ctx, pk, userJSON, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user in Cosmos DB: %w", err)
	}

	return user, nil
}

func (r *CosmosUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	pk := azcosmos.NewPartitionKeyString(id.String())

	resp, readItemErr := r.container.ReadItem(ctx, pk, id.String(), nil)
	if readItemErr != nil {
		return nil, ErrUserNotFound
	}

	var cosmosUser *CosmosUser
	if err := json.Unmarshal(resp.Value, &cosmosUser); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	user := NewUser(cosmosUser)

	return user, nil
}

func (r *CosmosUserRepository) GetUserByExternalID(ctx context.Context, externalID string) (*User, error) {
	query := "SELECT * FROM users u WHERE u.externalId = @externalID"
	params := []azcosmos.QueryParameter{
		{Name: "@externalID", Value: externalID},
	}

	queryOptions := &azcosmos.QueryOptions{QueryParameters: params}
	pager := r.container.NewQueryItemsPager(query, azcosmos.NewPartitionKey(), queryOptions)

	var cosmosUser *CosmosUser
	for pager.More() {
		resp, nextPageErr := pager.NextPage(ctx)
		if nextPageErr != nil {
			return nil, fmt.Errorf("query failed: %w", nextPageErr)
		}

		if len(resp.Items) > 0 {
			if err := json.Unmarshal(resp.Items[0], &cosmosUser); err != nil {
				return nil, fmt.Errorf("failed to unmarshal user: %w", err)
			}

			user := NewUser(cosmosUser)

			return user, nil
		}
	}

	return nil, ErrUserNotFound
}

func (r *CosmosUserRepository) UpdateUser(ctx context.Context, user *User) (*User, error) {
	userJSON, err := json.Marshal(NewCosmosUser(user))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user: %w", err)
	}

	pk := azcosmos.NewPartitionKeyString(user.ID.String())

	_, err = r.container.ReplaceItem(ctx, pk, user.ID.String(), userJSON, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to update user in Cosmos DB: %w", err)
	}

	return user, nil
}

func (r *CosmosUserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	pk := azcosmos.NewPartitionKeyString(id.String())

	_, err := r.container.DeleteItem(ctx, pk, id.String(), nil)
	if err != nil {
		var responseErr *azcore.ResponseError
		errors.As(err, &responseErr)
		panic(responseErr)
	}

	return nil
}

type CosmosUser struct {
	ID         uuid.UUID `json:"id"`
	ExternalID string    `json:"externalId"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func NewCosmosUser(u *User) *CosmosUser {
	return &CosmosUser{
		ID:         u.ID,
		ExternalID: u.ExternalID,
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}

func NewUser(cu *CosmosUser) *User {
	return &User{
		ID:         cu.ID,
		ExternalID: cu.ExternalID,
		FirstName:  cu.FirstName,
		LastName:   cu.LastName,
		CreatedAt:  cu.CreatedAt,
		UpdatedAt:  cu.UpdatedAt,
	}
}
