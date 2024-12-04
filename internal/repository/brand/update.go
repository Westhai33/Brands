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

func (r *BrandRepository) Update(ctx context.Context, brand *dto.Brand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BrandRepository.Update")
	defer span.Finish()

	query := `
        UPDATE brands SET 
            name = @name, 
            link = @link, 
            description = @description, 
            logo_url = @logo_url, 
            cover_image_url = @cover_image_url, 
            founded_year = @founded_year, 
            origin_country = @origin_country, 
            popularity = @popularity, 
            is_premium = @is_premium, 
            is_upcoming = @is_upcoming, 
            updated_at = NOW()
        WHERE id = @id AND is_deleted = false
    `
	args := pgx.NamedArgs{
		"id":              brand.ID,
		"name":            brand.Name,
		"link":            brand.Link,
		"description":     brand.Description,
		"logo_url":        brand.LogoURL,
		"cover_image_url": brand.CoverImageURL,
		"founded_year":    brand.FoundedYear,
		"origin_country":  brand.OriginCountry,
		"popularity":      brand.Popularity,
		"is_premium":      brand.IsPremium,
		"is_upcoming":     brand.IsUpcoming,
	}

	cmdTag, err := r.pool.Exec(ctx, query, args)

	if err != nil {
		span.LogFields(log.Error(err))
		r.log.Error().
			Err(err).
			Interface("brand", brand).
			Msg("Failed to update brand")

		if errors.Is(err, pgx.ErrNoRows) {
			return ErrBrandNotFound
		}
		return fmt.Errorf("unable to update brand: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		span.LogFields(log.Error(ErrBrandNotFound))
		r.log.Warn().
			Interface("brand", brand).
			Msg("No brand found to update")
		return ErrBrandNotFound
	}
	return nil
}
