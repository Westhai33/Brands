package brand

import (
	"Brands/internal/dto"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// GetAll получает все бренды с сортировкой по имени
func (r *BrandRepository) GetAll(ctx context.Context, sortBy string) ([]dto.Brand, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BrandRepository.GetAll")
	defer span.Finish()

	// Стартуем с базового SQL-запроса
	query := `
        SELECT * 
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
	brands, err = pgx.CollectRows(rows, pgx.RowToStructByName[dto.Brand])
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(
			log.Error(err),
			log.String("event", "collect_rows_error"),
		)
		r.log.Error().
			Err(err).
			Str("operation", "GetAll").
			Msg("Failed to collect rows into brands")
		return nil, fmt.Errorf("error collecting rows: %w", err)
	}
	return brands, nil
}
