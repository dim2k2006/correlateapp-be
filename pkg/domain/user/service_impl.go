package user

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type serviceImpl struct {
	repo Repository
}

var (
	ErrDuplicateExternalID = errors.New("duplicate external ID")
)

func NewService(repo Repository) Service {
	return &serviceImpl{
		repo: repo,
	}
}

func (s *serviceImpl) CreateUser(ctx context.Context, input CreateUserInput) (*User, error) {
	user, err := s.repo.GetUserByExternalID(ctx, input.ExternalID)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return nil, err
	}

	if user != nil {
		return nil, ErrDuplicateExternalID
	}

	newUser := &User{
		ID:         uuid.New(),
		ExternalID: input.ExternalID,
		FirstName:  input.FirstName,
		LastName:   input.LastName,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	createdUser, err := s.repo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *serviceImpl) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *serviceImpl) GetUserByExternalID(ctx context.Context, externalID string) (*User, error) {
	user, err := s.repo.GetUserByExternalID(ctx, externalID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *serviceImpl) UpdateUser(ctx context.Context, input UpdateUserInput) (*User, error) {
	user, err := s.repo.GetUserByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if input.FirstName != nil {
		user.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		user.LastName = *input.LastName
	}

	user.UpdatedAt = time.Now()

	updatedUser, err := s.repo.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (s *serviceImpl) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
