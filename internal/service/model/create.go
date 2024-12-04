package model

import (
	"Brands/internal/dto"
	"Brands/pkg/zerohook"
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// Create создает новую модель
func (s *ModelService) Create(ctx context.Context, model *dto.Model) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ModelService.Create")
	defer span.Finish()

	if model.Name == "" {
		err := fmt.Errorf("name is required")
		zerohook.Logger.Error().
			Err(err).
			Msg("Validation failed in Create")

		span.LogFields(
			log.String("event", "validation error"),
			log.String("reason", "name is required"),
		)

		return err
	}

	err := s.repo.Create(ctx, model)

	if err != nil {

		return err
	}
	return nil
}
