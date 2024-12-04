package model

import (
	"Brands/internal/dto"
	"context"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
)

func (s *ModelService) GetByID(ctx context.Context, id uuid.UUID) (*dto.Model, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ModelService.GetByID")
	defer span.Finish()

	model, err := s.repo.GetByID(ctx, id)
	if err != nil {

		return nil, err
	}
	return model, nil
}
