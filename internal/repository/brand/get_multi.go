package brand

import (
	"Brands/internal/dto"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"strings"
)

// GetAll получает все бренды с сортировкой по имени
func (r *BrandRepository) GetAll(ctx context.Context) ([]dto.Brand, error) {
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

func (r *BrandRepository) BrandsFilter(
	ctx context.Context,
	filter map[string]any,
	sortBy string,
) ([]dto.Brand, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BrandRepository.BrandsFilter")
	defer span.Finish()

	var queryBuilder strings.Builder
	var args []any
	argCounter := 1

	queryBuilder.WriteString(`
		SELECT * 
		FROM brands 
		WHERE is_deleted = false
	`)

	for key, value := range filter {
		if value == nil || value == "" {
			continue
		}
		switch key {
		case "name":
			queryBuilder.WriteString(fmt.Sprintf(" AND name ILIKE $%d", argCounter))
			args = append(args, "%"+value.(string)+"%")
		case "origin_country":
			queryBuilder.WriteString(fmt.Sprintf(" AND origin_country = $%d", argCounter))
			args = append(args, value.(string))
		case "popularity":
			queryBuilder.WriteString(fmt.Sprintf(" AND popularity = $%d", argCounter))
			args = append(args, value)
		case "is_premium":
			queryBuilder.WriteString(fmt.Sprintf(" AND is_premium = $%d", argCounter))
			args = append(args, value.(bool))
		case "is_upcoming":
			queryBuilder.WriteString(fmt.Sprintf(" AND is_upcoming = $%d", argCounter))
			args = append(args, value.(bool))
		case "founded_year":
			queryBuilder.WriteString(fmt.Sprintf(" AND founded_year = $%d", argCounter))
			args = append(args, value.(int))
		}
		argCounter++
	}

	if sortBy != "" {
		if strings.HasPrefix(sortBy, "-") {
			queryBuilder.WriteString(fmt.Sprintf(" ORDER BY %s DESC", strings.TrimPrefix(sortBy, "-")))
		} else {
			queryBuilder.WriteString(fmt.Sprintf(" ORDER BY %s ASC", sortBy))
		}
	} else {
		queryBuilder.WriteString(" ORDER BY created_at DESC")
	}

	r.log.Info().Str("query", queryBuilder.String()).Msg("Executing query with sorting")
	query := queryBuilder.String()

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		span.LogFields(
			log.Error(err),
			log.String("query", query),
		)
		r.log.Error().
			Err(err).
			Str("operation", "ModelsFilter").
			Msg("Failed to execute ModelsFilter query")
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var brands []dto.Brand
	brands, err = pgx.CollectRows(rows, pgx.RowToStructByName[dto.Brand])
	if err != nil {
		span.LogFields(
			log.Error(err),
			log.String("event", "collect_rows_error"),
		)
		r.log.Error().
			Err(err).
			Str("operation", "BrandsFilter").
			Msg("Failed to collect rows into brands")
		return nil, fmt.Errorf("error collecting rows: %w", err)
	}
	return brands, nil
}
