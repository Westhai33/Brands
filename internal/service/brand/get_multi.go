package brand

import (
	"Brands/internal/dto"
	"Brands/pkg/zerohook"
	"context"
	"github.com/opentracing/opentracing-go"
)

// GetAll получает все бренды с фильтрацией и сортировкой
func (s *Service) GetAll(ctx context.Context, filter map[string]interface{}, sort string) ([]dto.Brand, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BrandService.GetAll")
	defer span.Finish()

	zerohook.Logger.Info().
		Interface("filter", filter).
		Str("sort", sort).
		Msg("Starting GetAll operation")

	task := make(chan struct {
		brands []dto.Brand
		err    error
	}, 1)

	s.workerPool.Submit(func(workerID int) {
		workerSpan, _ := opentracing.StartSpanFromContext(ctx, "Worker.GetAllBrands")
		defer workerSpan.Finish()

		brands, err := s.repo.GetAll(ctx, filter, sort)
		task <- struct {
			brands []dto.Brand
			err    error
		}{
			brands: brands,
			err:    err,
		}
		defer close(task)
	})

	result := <-task
	if result.err != nil {
		return nil, result.err
	}
	return result.brands, nil
}
