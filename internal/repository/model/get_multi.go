package model

import (
	"Brands/internal/dto"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

func (r *ModelRepository) GetAll(ctx context.Context) ([]dto.Model, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ModelRepository.GetAll")
	defer span.Finish()

	query := `
        SELECT * 
        FROM models 
        WHERE is_deleted = false
        ORDER BY name ASC
    `
	// Выполняем запрос
	rows, err := r.pool.Query(ctx, query)
	if err != nil {

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
	models, err = pgx.CollectRows(rows, pgx.RowToStructByName[dto.Model])
	if err != nil {

		span.LogFields(
			log.Error(err),
			log.String("event", "collect_rows_error"),
		)
		r.log.Error().
			Err(err).
			Str("operation", "GetAll").
			Msg("Failed to collect rows into models")
		return nil, fmt.Errorf("error collecting rows: %w", err)
	}
	return models, nil
}
