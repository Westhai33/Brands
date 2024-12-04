package model

import (
	"Brands/internal/dto"
	"context"
	"github.com/opentracing/opentracing-go"
)

func (s *Service) GetByID(ctx context.Context, id int64) (*dto.Model, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ModelService.GetByID")
	defer span.Finish()

	task := make(chan struct {
		model *dto.Model
		err   error
	}, 1)

	s.workerPool.Submit(func(workerID int) {
		workerSpan, _ := opentracing.StartSpanFromContext(ctx, "Worker.GetByModelID")
		defer workerSpan.Finish()

		model, err := s.repo.GetByID(ctx, id)
		task <- struct {
			model *dto.Model
			err   error
		}{model: model, err: err}

		defer close(task)
	})

	result := <-task
	if result.err != nil {
		span.SetTag("error", true)
		return nil, result.err
	}
	return result.model, nil
}
