package brand

import (
	"Brands/internal/dto"
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// GetAll получает все бренды с сортировкой по имени
func (r *BrandRepository) GetAll(ctx context.Context, sortBy string) ([]dto.Brand, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BrandRepository.GetAll")
	defer span.Finish()

	// Стартуем с базового SQL-запроса
	query := `
        SELECT id, name, link, description, logo_url, cover_image_url, founded_year, 
               origin_country, popularity, is_premium, is_upcoming, created_at, updated_at 
        FROM brands 
        WHERE is_deleted = false
        ORDER BY name ASC
    `

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(log.Error(err))
		r.log.Error().
			Err(err).
			Msg("Failed to execute GetAll query")
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var brands []dto.Brand

	for rows.Next() {
		var brand dto.Brand
		if err = rows.Scan(
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
			&brand.CreatedAt,
			&brand.UpdatedAt,
		); err != nil {
			span.SetTag("error", true)
			span.LogFields(log.Error(err))
			r.log.Error().Err(err).Msg("Failed to scan row in GetAll")
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		brands = append(brands, brand)
	}

	// Проверка ошибок после итерации по строкам
	if rows.Err() != nil {
		span.SetTag("error", true)
		span.LogFields(
			log.String("event", "rows_iteration_error"),
			log.String("error", rows.Err().Error()),
		)
		r.log.Error().
			Err(rows.Err()).
			Str("operation", "GetAll").
			Msg("Error iterating over rows in GetAll")
		return nil, fmt.Errorf("error iterating over rows: %w", rows.Err())
	}

	return brands, nil
}
