package brand

import (
	"Brands/internal/dto"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// GetByID получает бренд по ID
func (r *BrandRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*dto.Brand, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BrandRepository.GetByID")
	defer span.Finish()

	query := `
        SELECT id, name, link, description, logo_url, cover_image_url, founded_year, 
               origin_country, popularity, is_premium, is_upcoming, is_deleted, created_at, updated_at
        FROM brands
        WHERE id = $1 and is_deleted = false
    `
	row, err := r.pool.Query(ctx, query, id)
	if err != nil {
		span.LogFields(log.Error(err))
		r.log.Error().Err(err).Str("brand_id", id.String()).Msg("Failed to fetch brand by ID")
		return nil, fmt.Errorf("unable to get brand by id: %w", err)
	}

	var brands []dto.Brand
	brands, err = pgx.CollectRows(row, pgx.RowToStructByName[dto.Brand])
	if err != nil {
		span.LogFields(log.Error(err))
		r.log.Error().Err(err).Str("brand_id", id.String()).Msg("Failed to collect rows")
		return nil, fmt.Errorf("unable to collect rows: %w", err)
	}

	if len(brands) == 0 {
		r.log.Warn().Str("brand_id", id.String()).Msg("Brand not found")
		return nil, ErrBrandNotFound
	}

	if len(brands) > 1 {
		err = fmt.Errorf("multiple brands found with id %s", id.String())
		span.LogFields(log.Error(err))
		r.log.Error().Str("brand_id", id.String()).Msg("Multiple brands found with the same ID")
		return nil, err
	}
	return &brands[0], nil
}
