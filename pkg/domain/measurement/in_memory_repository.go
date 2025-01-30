package measurement

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
)

type InMemoryRepository struct {
	mu           sync.RWMutex
	measurements map[uuid.UUID]Measurement
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		measurements: make(map[uuid.UUID]Measurement),
	}
}

var (
	ErrMeasurementNotFound = errors.New("measurement not found")
	ErrInvalidMeasurement  = errors.New("invalid measurement")
)

func (repo *InMemoryRepository) CreateMeasurement(_ context.Context, measurement Measurement) (Measurement, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.measurements[measurement.GetID()] = measurement

	return measurement, nil
}

func (repo *InMemoryRepository) ListMeasurementsByUser(_ context.Context, userID uuid.UUID) ([]Measurement, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	var results []Measurement
	for _, measurement := range repo.measurements {
		if measurement.GetUserID() == userID {
			results = append(results, measurement)
		}
	}

	return results, nil
}

func (repo *InMemoryRepository) ListMeasurementsByParameter(_ context.Context, parameterID uuid.UUID) ([]Measurement, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	var results []Measurement
	for _, measurement := range repo.measurements {
		if measurement.GetParameterID() == parameterID {
			results = append(results, measurement)
		}
	}

	return results, nil
}

func (repo *InMemoryRepository) DeleteMeasurement(_ context.Context, id uuid.UUID) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, ok := repo.measurements[id]; !ok {
		return ErrMeasurementNotFound
	}

	delete(repo.measurements, id)

	return nil
}
