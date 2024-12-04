package brand

import (
	"Brands/internal/dto"
	"context"
	"github.com/opentracing/opentracing-go"
)

// GetAll получает все бренды с фильтрацией и сортировкой
func (s *BrandService) GetAll(ctx context.Context) ([]dto.Brand, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BrandService.GetAll")
	defer span.Finish()

	brands, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return brands, nil
}
