package model

import (
	"Brands/internal/dto"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"strings"
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
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		span.LogFields(log.Error(err))
		r.log.Error().
			Err(err).
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

func (r *ModelRepository) ModelsFilter(ctx context.Context, filter map[string]any, sortBy string) ([]dto.Model, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ModelRepository.ModelsFilter")
	defer span.Finish()

	var queryBuilder strings.Builder
	var args []interface{}
	argCounter := 1

	queryBuilder.WriteString(`
        SELECT * 
        FROM models 
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
		case "brand_id":
			queryBuilder.WriteString(fmt.Sprintf(" AND brand_id = $%d", argCounter))
			args = append(args, value.(string))
		case "is_limited":
			queryBuilder.WriteString(fmt.Sprintf(" AND is_limited = $%d", argCounter))
			args = append(args, value.(bool))
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
		queryBuilder.WriteString(" ORDER BY name ASC")
	}

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

	var models []dto.Model
	models, err = pgx.CollectRows(rows, pgx.RowToStructByName[dto.Model])
	if err != nil {
		span.LogFields(
			log.Error(err),
			log.String("event", "collect_rows_error"),
		)
		r.log.Error().
			Err(err).
			Str("operation", "ModelsFilter").
			Msg("Failed to collect rows into models")
		return nil, fmt.Errorf("error collecting rows: %w", err)
	}

	return models, nil
}
