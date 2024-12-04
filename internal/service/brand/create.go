package brand

import (
	"Brands/internal/dto"
	"context"
	"github.com/opentracing/opentracing-go"
)

// Create создает новый бренд
func (s *BrandService) Create(ctx context.Context, brand *dto.Brand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BrandService.Create")
	defer span.Finish()
	err := s.repo.Create(ctx, brand)

	if err != nil {

		return err
	}
	return nil
}
