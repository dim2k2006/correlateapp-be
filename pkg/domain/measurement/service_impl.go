package measurement

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dim2k2006/correlateapp-be/pkg/domain/parameter"
	"github.com/google/uuid"
)

type ServiceImpl struct {
	repo             Repository
	parameterService parameter.Service
}

func NewService(repo Repository, parameterService parameter.Service) Service {
	return &ServiceImpl{
		repo:             repo,
		parameterService: parameterService,
	}
}

func (s *ServiceImpl) CreateMeasurement(ctx context.Context, input CreateMeasurementInput) (Measurement, error) {
	measurementParameter, err := s.parameterService.GetParameterByID(ctx, input.ParameterID)
	if err != nil {
		return nil, err
	}

	switch input.Type {
	case MeasurementTypeFloat:
		v, ok := input.Value.(float64)
		if !ok {
			return nil, errors.New("invalid value type for float measurement")
		}
		measurement := &FloatMeasurement{
			BaseMeasurement: BaseMeasurement{
				Type:        MeasurementTypeFloat,
				ID:          uuid.New(),
				UserID:      input.UserID,
				ParameterID: measurementParameter.ID,
				Timestamp:   input.Timestamp,
				Notes:       input.Notes,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			Value: v,
		}
		return s.repo.CreateMeasurement(ctx, measurement)
	// case MeasurementTypeBoolean:
	//	// parse `input.Value` as bool
	//	b, ok := input.Value.(bool)
	//	if !ok {
	//		return nil, errors.New("expected boolean value for measurement")
	//	}
	//	meas := BooleanMeasurement{
	//		BaseMeasurement: BaseMeasurement{
	//			Type:        MeasurementTypeBoolean,
	//			ID:          uuid.New(),
	//			UserID:      input.UserID,
	//			ParameterID: input.ParameterID,
	//			Timestamp:   input.Timestamp,
	//			Notes:       input.Notes,
	//			CreatedAt:   time.Now(),
	//			UpdatedAt:   time.Now(),
	//		},
	//		Value: b,
	//	}
	//	return s.repo.CreateMeasurement(ctx, meas)
	default:
		return nil, fmt.Errorf("unsupported measurement type: %s", input.Type)
	}
}

func (s *ServiceImpl) ListMeasurementsByUser(ctx context.Context, userID uuid.UUID) ([]Measurement, error) {
	return s.repo.ListMeasurementsByUser(ctx, userID)
}

func (s *ServiceImpl) ListMeasurementsByParameter(ctx context.Context, parameterID uuid.UUID) ([]Measurement, error) {
	return s.repo.ListMeasurementsByParameter(ctx, parameterID)
}

func (s *ServiceImpl) DeleteMeasurement(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteMeasurement(ctx, id)
}
