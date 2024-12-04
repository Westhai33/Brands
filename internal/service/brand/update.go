package brand

import (
	"Brands/internal/dto"
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// Update обновляет данные бренда
func (s *Service) Update(ctx context.Context, brand *dto.Brand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BrandService.Update")
	defer span.Finish()

	if brand.Name == "" {
		err := fmt.Errorf("name is required")
		s.log.Error().
			Err(err).
			Int64("brand_id", brand.ID).
			Msg("Validation failed in Update")

		// Добавляем ошибку в спан
		span.SetTag("error", true)
		span.LogFields(
			log.String("event", "validation error"),
			log.String("reason", "name is required"),
			log.Int64("brand.id", brand.ID),
		)

		return err
	}

	task := make(chan error, 1)

	s.workerPool.Submit(func(workerID int) {
		workerSpan, _ := opentracing.StartSpanFromContext(ctx, "Worker.UpdateBrand")
		defer workerSpan.Finish()
		err := s.repo.Update(ctx, brand)
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
