package model

import (
	"Brands/internal/dto"
	"Brands/pkg/zerohook"
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// Update обновляет данные модели
func (s *ModelService) Update(ctx context.Context, model *dto.Model) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ModelService.Update")
	defer span.Finish()

	if model.Name == "" {
		err := fmt.Errorf("name is required")
		zerohook.Logger.Error().
			Err(err).
			Msg("Validation failed in Update")

		span.LogFields(
			log.String("event", "validation error"),
			log.String("reason", "name is required"),
			log.Int64("model.id", model.ID),
		)
		return err
	}

	err := s.repo.Update(ctx, model)
	if err != nil {

		return err
	}
	return nil
}
