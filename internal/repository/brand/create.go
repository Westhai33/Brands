package brand

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"

	"Brands/internal/dto"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// Create создает новый бренд
func (r *BrandRepository) Create(
	ctx context.Context,
	brand *dto.Brand,
) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BrandRepository.Create")
	defer span.Finish()

	query := `
		INSERT INTO brands (id, name, link, description, logo_url, cover_image_url, founded_year, origin_country, popularity, is_premium, is_upcoming, created_at, updated_at, is_deleted)
		VALUES (@id, @name, @link, @description, @logo_url, @cover_image_url, @founded_year, @origin_country, @popularity, @is_premium, @is_upcoming, NOW(), NOW(), false)
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
	_, err := r.pool.Exec(ctx, query, args)
	if err != nil {

		span.LogFields(log.Error(err), log.Object("brand", brand))
		r.log.Error().Interface("brand", brand).Err(err).Msg("Failed to create brand")
		return fmt.Errorf("unable to create brand: %w", err)
	}

	return nil
}
