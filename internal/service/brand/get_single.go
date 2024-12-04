package brand

import (
	"Brands/internal/dto"
	"context"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
)

// GetByID получает бренд по ID
func (s *BrandService) GetByID(ctx context.Context, id uuid.UUID) (*dto.Brand, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BrandService.GetByID")
	defer span.Finish()
	brand, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return brand, nil
}
