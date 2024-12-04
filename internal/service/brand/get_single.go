package brand

import (
	"Brands/internal/dto"
	"context"
	"github.com/opentracing/opentracing-go"
)

// GetByID получает бренд по ID
func (s *BrandService) GetByID(ctx context.Context, id int64) (*dto.Brand, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BrandService.GetByID")
	defer span.Finish()

	task := make(chan struct {
		brand *dto.Brand
		err   error
	}, 1)

	s.workerPool.Submit(
		func(workerID int) {
			workerSpan, _ := opentracing.StartSpanFromContext(ctx, "Worker.GetBrandByID")
			defer workerSpan.Finish()

			brand, err := s.repo.GetByID(ctx, id)
			task <- struct {
				brand *dto.Brand
				err   error
			}{
				brand: brand,
				err:   err,
			}

			defer close(task)
		})

	result := <-task

	if result.err != nil {
		span.SetTag("error", true)
		return nil, result.err
	}
	return result.brand, nil
}
