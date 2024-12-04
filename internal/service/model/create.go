package model

import (
	"context"
	"fmt"

	"Brands/internal/dto"
	"Brands/pkg/zerohook"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// Create создает новую модель
func (s *Service) Create(ctx context.Context, model *dto.Model) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ModelService.Create")
	defer span.Finish()

	if model.Name == "" {
		err := fmt.Errorf("name is required")
		zerohook.Logger.Error().
			Err(err).
			Msg("Validation failed in Create")

		span.SetTag("error", true)
		span.LogFields(
			log.String("event", "validation error"),
			log.String("reason", "name is required"),
		)

		return 0, err
	}

	id, err := s.repo.Create(ctx, model)

	if err != nil {
		span.SetTag("error", true)
		return 0, err
	}
	return id, nil
}
