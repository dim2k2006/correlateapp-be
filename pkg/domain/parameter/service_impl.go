package parameter

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ServiceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &ServiceImpl{
		repo: repo,
	}
}

func (s *ServiceImpl) CreateParameter(ctx context.Context, input CreateParameterInput) (*Parameter, error) {
	parameter := &Parameter{
		ID:          uuid.New(),
		UserID:      input.UserID,
		Name:        input.Name,
		Description: input.Description,
		DataType:    input.DataType,
		Unit:        input.Unit,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createdParameter, err := s.repo.CreateParameter(ctx, parameter)
	if err != nil {
		return nil, err
	}

	return createdParameter, nil
}

func (s *ServiceImpl) GetParameterByID(ctx context.Context, id uuid.UUID) (*Parameter, error) {
	parameter, err := s.repo.GetParameterByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return parameter, nil
}

func (s *ServiceImpl) ListParametersByUser(ctx context.Context, userID uuid.UUID) ([]*Parameter, error) {
	parameters, err := s.repo.ListParametersByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return parameters, nil
}

func (s *ServiceImpl) UpdateParameter(ctx context.Context, input UpdateParameterInput) (*Parameter, error) {
	parameter, err := s.repo.GetParameterByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		parameter.Name = *input.Name
	}
	if input.Description != nil {
		parameter.Description = *input.Description
	}
	if input.Unit != nil {
		parameter.Unit = *input.Unit
	}

	parameter.UpdatedAt = time.Now()

	updatedParameter, err := s.repo.UpdateParameter(ctx, parameter)
	if err != nil {
		return nil, err
	}

	return updatedParameter, nil
}

func (s *ServiceImpl) DeleteParameter(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteParameter(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
