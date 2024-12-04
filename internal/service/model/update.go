package model

import (
	"Brands/internal/dto"
	"context"
	"github.com/opentracing/opentracing-go"
)

// Update обновляет данные модели
func (s *ModelService) Update(ctx context.Context, model *dto.Model) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ModelService.Update")
	defer span.Finish()

	err := s.repo.Update(ctx, model)
	if err != nil {
		return err
	}
	return nil
}
