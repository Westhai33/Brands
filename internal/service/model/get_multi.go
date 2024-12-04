package model

import (
	"Brands/internal/dto"
	"Brands/pkg/zerohook"
	"context"
	"github.com/opentracing/opentracing-go"
)

// GetAll получает все модели с фильтрацией и сортировкой
func (s *Service) GetAll(
	ctx context.Context,
	filter map[string]interface{},
	sort string,
) ([]dto.Model, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ModelService.GetAll")
	defer span.Finish()
	span.SetTag("filter", filter)
	span.SetTag("sort", sort)

	zerohook.Logger.Info().
		Interface("filter", filter).
		Str("sort", sort).
		Msg("Starting GetAll operation")

	models, err := s.repo.GetAll(ctx, filter, sort)

	if err != nil {
		return nil, err
	}
	return models, nil
}
