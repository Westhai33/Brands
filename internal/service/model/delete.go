package model

import (
	"context"
	"github.com/opentracing/opentracing-go"
)

// SoftDelete мягко удаляет модель
func (s *Service) SoftDelete(ctx context.Context, id int64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ModelService.SoftDelete")
	defer span.Finish()

	task := make(chan error, 1)

	s.workerPool.Submit(func(workerID int) {
		workerSpan, _ := opentracing.StartSpanFromContext(ctx, "Worker.SoftDeleteModel")
		defer workerSpan.Finish()
		err := s.repo.SoftDelete(ctx, id)
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

// Restore восстанавливает мягко удалённую модель
func (s *Service) Restore(ctx context.Context, id int64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ModelService.Restore")
	defer span.Finish()
	task := make(chan error, 1)

	s.workerPool.Submit(func(workerID int) {
		workerSpan, _ := opentracing.StartSpanFromContext(ctx, "Worker.RestoreModel")
		defer workerSpan.Finish()

		err := s.repo.Restore(ctx, id)
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
