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

	task := make(chan struct {
		models []dto.Model
		err    error
	}, 1)

	s.workerPool.Submit(func(workerID int) {
		workerSpan, _ := opentracing.StartSpanFromContext(ctx, "Worker.GetAllModels")
		defer workerSpan.Finish()
		models, err := s.repo.GetAll(ctx, filter, sort)
		task <- struct {
			models []dto.Model
			err    error
		}{models: models, err: err}

		defer close(task)
	})
	result := <-task
	if result.err != nil {
		return nil, result.err
	}
	return result.models, nil
}
