package model

import (
	"Brands/internal/dto"
	"context"
	"github.com/opentracing/opentracing-go"
)

// GetAll получает все модели с фильтрацией и сортировкой
func (s *ModelService) GetAll(
	ctx context.Context,
	filter map[string]interface{},
	sortBy string,
) ([]dto.Model, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ModelService.GetAll")
	defer span.Finish()

	models, err := s.repo.GetAll(ctx, filter, sortBy)
	if err != nil {
		return nil, err
	}
	return models, nil
}
