package brand

import (
	"context"
	"fmt"

	"Brands/internal/dto"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// Create создает новый бренд
func (r *Repository) Create(
	ctx context.Context,
	brand *dto.Brand,
) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BrandRepository.Create")
	defer span.Finish()

	query := `
		INSERT INTO brands (name, link, description, logo_url, cover_image_url, founded_year, origin_country, popularity, is_premium, is_upcoming, created_at, updated_at, is_deleted)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW(), false)
		RETURNING id`

	err := r.pool.QueryRow(
		ctx, query,
		brand.Name,
		brand.Link,
		brand.Description,
		brand.LogoURL,
		brand.CoverImageURL,
		brand.FoundedYear,
		brand.OriginCountry,
		brand.Popularity,
		brand.IsPremium,
		brand.IsUpcoming,
	).Scan(&brand.ID)

	if err != nil {
		span.SetTag("error", true)
		span.LogFields(log.Error(err))
		return 0, fmt.Errorf("unable to create brand: %w", err)
	}

	return brand.ID, nil
}
