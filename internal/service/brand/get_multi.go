package brand

import (
	"Brands/internal/dto"
	"context"
	"github.com/opentracing/opentracing-go"
)

// GetAll получает все бренды с фильтрацией и сортировкой
func (s *BrandService) GetAll(ctx context.Context, filter map[string]interface{}, sortBy string) ([]dto.Brand, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BrandService.GetAll")
	defer span.Finish()

	s.log.Info().
		Interface("filter", filter).
		Str("sortBy", sortBy).
		Msg("Fetching all brands with filters and sorting")

	brands, err := s.repo.GetAll(ctx, filter, sortBy)
	if err != nil {
		s.log.Error().
			Err(err).
			Msg("Failed to fetch brands in BrandService.GetAll")
		return nil, err
	}

	s.log.Info().
		Int("brands_count", len(brands)).
		Msg("Successfully fetched brands in BrandService.GetAll")
	return brands, nil
}
