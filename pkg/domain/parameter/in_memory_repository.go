package parameter

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
)

type InMemoryRepository struct {
	mu         sync.RWMutex
	parameters map[uuid.UUID]*Parameter
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		parameters: make(map[uuid.UUID]*Parameter),
	}
}

var (
	ErrParameterNotFound = errors.New("parameter not found")
)

func (r *InMemoryRepository) CreateParameter(_ context.Context, parameter *Parameter) (*Parameter, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.parameters[parameter.ID] = parameter

	return parameter, nil
}

func (r *InMemoryRepository) GetParameterByID(_ context.Context, id uuid.UUID) (*Parameter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	parameter, ok := r.parameters[id]
	if !ok {
		return nil, ErrParameterNotFound
	}

	return parameter, nil
}

func (r *InMemoryRepository) ListParametersByUser(_ context.Context, userID uuid.UUID) ([]*Parameter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var parameters []*Parameter
	for _, parameter := range r.parameters {
		if parameter.UserID == userID {
			parameters = append(parameters, parameter)
		}
	}

	return parameters, nil
}

func (r *InMemoryRepository) UpdateParameter(_ context.Context, param *Parameter) (*Parameter, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.parameters[param.ID]; !ok {
		return nil, ErrParameterNotFound
	}

	r.parameters[param.ID] = param

	return param, nil
}

func (r *InMemoryRepository) DeleteParameter(_ context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.parameters[id]; !ok {
		return ErrParameterNotFound
	}

	delete(r.parameters, id)

	return nil
}
