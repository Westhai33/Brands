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

	brands, err := s.repo.GetAll(ctx, filter, sort)
	if err != nil {
		return nil, err
	}
	return brands, nil
}
