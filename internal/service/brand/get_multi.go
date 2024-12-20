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
	s.log.Info().
		Int("brands_count", len(brands)).
		Msg("Successfully fetched brands in BrandService.GetAll")
	return brands, nil
}

func (s *BrandService) BrandsFilter(
	ctx context.Context,
	filter map[string]any,
	sortBy string,
) ([]dto.Brand, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BrandService.GetAll")
	defer span.Finish()

	s.log.Info().
		Interface("filter", filter).
		Str("sortBy", sortBy).
		Msg("Fetching all brands with filters and sorting")

	brands, err := s.repo.BrandsFilter(ctx, filter, sortBy)

	if err != nil {
		return nil, err
	}
	s.log.Info().
		Int("brands_count", len(brands)).
		Msg("Successfully fetched brands in BrandService.BrandsFilter")
	return brands, nil
}
