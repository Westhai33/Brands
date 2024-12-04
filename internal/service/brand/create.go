package brand

import (
	"Brands/internal/dto"
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// Create создает новый бренд
func (s *BrandService) Create(ctx context.Context, brand *dto.Brand) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BrandService.Create")
	defer span.Finish()
	s.log.Info().Msg("Starting Create operation")

	if brand.Name == "" {
		err := fmt.Errorf("name is required")
		span.SetTag("error", true)
		span.LogFields(
			log.String("event", "validation_failed"),
			log.Error(err),
		)
		s.log.Error().Err(err).Msg("Validation failed in Create")
		return 0, err
	}

	task := make(chan struct {
		id  int64
		err error
	}, 1)

	s.workerPool.Submit(
		func(workerID int) {
			workerSpan := opentracing.StartSpan("worker.CreateBrand", opentracing.ChildOf(span.Context()))
			defer workerSpan.Finish()

			id, err := s.repo.Create(ctx, brand)
			task <- struct {
				id  int64
				err error
			}{
				id:  id,
				err: err,
			}
			defer close(task)
		},
	)
	result := <-task

	if result.err != nil {
		span.SetTag("error", true)
		return 0, result.err
	}
	return result.id, nil
}
