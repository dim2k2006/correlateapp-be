package user_test

import (
	"context"
	"testing"

	"github.com/dim2k2006/correlateapp-be/pkg/domain/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_CreateUser(t *testing.T) {
	repo := user.NewInMemoryRepository()
	svc := user.NewService(repo)

	input := user.CreateUserInput{
		ExternalID: "b6541d6a-7987-42ce-b124-018667a76bd5",
		FirstName:  "John",
		LastName:   "Doe",
	}

	createdUser, err := svc.CreateUser(context.Background(), input)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, createdUser.ID)
	assert.Equal(t, input.ExternalID, createdUser.ExternalID)
	assert.Equal(t, input.FirstName, createdUser.FirstName)
	assert.Equal(t, input.LastName, createdUser.LastName)
}

func TestService_CreateUser_DuplicateExternalID(t *testing.T) {
	repo := user.NewInMemoryRepository()
	svc := user.NewService(repo)

	input := user.CreateUserInput{
		ExternalID: "b6541d6a-7987-42ce-b124-018667a76bd5",
		FirstName:  "John",
		LastName:   "Doe",
	}

	_, err := svc.CreateUser(context.Background(), input)
	require.NoError(t, err)

	_, err = svc.CreateUser(context.Background(), input)
	require.Error(t, err)
}
