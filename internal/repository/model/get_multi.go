package model

import (
	"Brands/internal/dto"
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

func (r *ModelRepository) GetAll(ctx context.Context) ([]dto.Model, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ModelRepository.GetAll")
	defer span.Finish()

	query := `
        SELECT id, brand_id, name, release_date, is_upcoming, is_limited, is_deleted, created_at, updated_at 
        FROM models 
        WHERE is_deleted = false
        ORDER BY name ASC
    `

	// Выполняем запрос
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(
			log.Error(err),
			log.String("query", query),
		)
		r.log.Error().
			Err(err).
			Str("operation", "GetAll").
			Msg("Failed to execute GetAll query")
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var models []dto.Model

	// Обрабатываем строки результата
	for rows.Next() {
		var model dto.Model
		if err := rows.Scan(
			&model.ID,
			&model.BrandID,
			&model.Name,
			&model.ReleaseDate,
			&model.IsUpcoming,
			&model.IsLimited,
			&model.IsDeleted,
			&model.CreatedAt,
			&model.UpdatedAt,
		); err != nil {
			span.SetTag("error", true)
			span.LogFields(
				log.String("event", "row_scan_error"),
				log.String("message", err.Error()),
			)
			r.log.Error().
				Err(err).
				Str("operation", "GetAll").
				Msg("Failed to scan row in GetAll")
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		models = append(models, model)
	}

	if rows.Err() != nil {
		span.SetTag("error", true)
		span.LogFields(
			log.String("event", "rows_iteration_error"),
			log.String("message", rows.Err().Error()),
		)
		r.log.Error().
			Err(rows.Err()).
			Msg("Error iterating over rows in GetAll")
		return nil, fmt.Errorf("error iterating over rows: %w", rows.Err())
	}

	return models, nil
}
