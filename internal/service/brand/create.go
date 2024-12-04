package brand

import (
	"Brands/internal/dto"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// Create создает новый бренд
func (s *BrandService) Create(ctx context.Context, brand *dto.Brand) (uuid.UUID, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BrandService.Create")
	defer span.Finish()
	if brand.Name == "" {
		err := fmt.Errorf("name is required")
		span.SetTag("error", true)
		span.LogFields(
			log.String("event", "validation_failed"),
			log.Error(err),
		)
		s.log.Error().Err(err).Msg("Validation failed in Create")
		return uuid.Nil, err
	}

	id, err := s.repo.Create(ctx, brand)

	if err != nil {
		span.SetTag("error", true)
		return uuid.Nil, err
	}
	return id, nil
}
