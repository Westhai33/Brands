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

func (r *Repository) Update(ctx context.Context, brand *dto.Brand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BrandRepository.Update")
	defer span.Finish()

	query := `
        UPDATE brands SET 
            name = $2, 
            link = $3, 
            description = $4, 
            logo_url = $5, 
            cover_image_url = $6, 
            founded_year = $7, 
            origin_country = $8, 
            popularity = $9, 
            is_premium = $10, 
            is_upcoming = $11, 
            updated_at = NOW()
        WHERE id = $1 AND is_deleted = false
        RETURNING id
    `

	err := r.pool.QueryRow(ctx, query,
		brand.ID,
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
	).Scan()

	if err != nil {
		span.SetTag("error", true)
		span.LogFields(log.Error(err))
		r.log.Error().
			Err(err).
			Int64("brand_id", brand.ID).
			Msg("Failed to update brand")

		if errors.Is(err, pgx.ErrNoRows) {
			return ErrBrandNotFound
		}
		return fmt.Errorf("unable to update brand: %w", err)
	}
	return nil
}
