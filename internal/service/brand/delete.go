package brand

import (
	"context"
	"github.com/opentracing/opentracing-go"
)

// SoftDelete мягко удаляет бренд
func (s *Service) SoftDelete(ctx context.Context, id int64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BrandService.SoftDelete")
	defer span.Finish()

	task := make(chan error, 1)
	s.workerPool.Submit(func(workerID int) {
		workerSpan, _ := opentracing.StartSpanFromContext(ctx, "Worker.BrandSoftDelete")
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

// Restore восстанавливает мягко удалённый бренд
func (s *Service) Restore(ctx context.Context, id int64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BrandService.Restore")
	defer span.Finish()

	task := make(chan error, 1)
	s.workerPool.Submit(func(workerID int) {
		workerSpan, _ := opentracing.StartSpanFromContext(ctx, "Worker.RestoreBrand")
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
