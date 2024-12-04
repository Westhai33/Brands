package brand

import (
	"Brands/internal/dto"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

// GetByID получает бренд по ID
func (r *Repository) GetByID(
	ctx context.Context,
	id int64,
) (*dto.Brand, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BrandRepository.GetByID")
	defer span.Finish()

	query := `
        SELECT id, name, link, description, logo_url, cover_image_url, founded_year, 
               origin_country, popularity, is_premium, is_upcoming, is_deleted, created_at, updated_at
        FROM brands
        WHERE id = $1 and is_deleted = false
    `

	row := r.pool.QueryRow(ctx, query, id)

	brand := dto.Brand{}
	err := row.Scan(
		&brand.ID,
		&brand.Name,
		&brand.Link,
		&brand.Description,
		&brand.LogoURL,
		&brand.CoverImageURL,
		&brand.FoundedYear,
		&brand.OriginCountry,
		&brand.Popularity,
		&brand.IsPremium,
		&brand.IsUpcoming,
		&brand.IsDeleted,
		&brand.CreatedAt,
		&brand.UpdatedAt,
	)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(log.Error(err))
		r.log.Error().Err(err).Int64("brand_id", id).Msg("Failed to fetch brand by ID")

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrBrandNotFound
		}
		return nil, fmt.Errorf("unable to get brand by id: %w", err)
	}
	return &brand, nil
}
