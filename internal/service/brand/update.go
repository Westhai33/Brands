package brand

import (
	"Brands/internal/dto"
	"context"
	"github.com/opentracing/opentracing-go"
)

// Update обновляет данные бренда
func (s *BrandService) Update(ctx context.Context, brand *dto.Brand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BrandService.Update")
	defer span.Finish()
	err := s.repo.Update(ctx, brand)
	if err != nil {
		return err
	}
	return nil
}
