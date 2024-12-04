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
func (s *Service) Update(ctx context.Context, model *dto.Model) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Service.Update")
	defer span.Finish()

	if model.Name == "" {
		err := fmt.Errorf("name is required")
		zerohook.Logger.Error().
			Err(err).
			Msg("Validation failed in Update")

		span.SetTag("error", true)
		span.LogFields(
			log.String("event", "validation error"),
			log.String("reason", "name is required"),
			log.Int64("model.id", model.ID),
		)
		return err
	}

	task := make(chan error, 1)

	s.workerPool.Submit(func(workerID int) {
		workerSpan, _ := opentracing.StartSpanFromContext(ctx, "Repository.Update")
		defer workerSpan.Finish()

		err := s.repo.Update(ctx, model)
		task <- err
		defer close(task)
	})

	err := <-task
	if err != nil {
		span.SetTag("error", true)
		return err
	}
	return nil
}
