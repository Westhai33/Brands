package model

import (
	"Brands/internal/dto"
	"context"
	"fmt"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// GetAll получает все модели с фильтрацией и сортировкой
func (r *Repository) GetAll(ctx context.Context, filter map[string]interface{}, sortBy string) ([]dto.Model, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ModelRepository.GetAll")
	defer span.Finish()

	var queryBuilder strings.Builder
	var args []interface{}
	argCounter := 1

	baseQuery := `
        SELECT id, brand_id, name, release_date, is_upcoming, is_limited, is_deleted, created_at, updated_at 
        FROM models 
        WHERE is_deleted = false
    `
	queryBuilder.WriteString(baseQuery)

	// Фильтрация по имени
	if name, ok := filter["name"]; ok {
		if nameStr, ok := name.(string); ok && nameStr != "" {
			queryBuilder.WriteString(fmt.Sprintf(" AND name ILIKE $%d", argCounter))
			args = append(args, "%"+nameStr+"%")
			argCounter++
		}
	}

	// Фильтрация по brand_id
	if brandID, ok := filter["brand_id"]; ok {
		if brandIDInt, ok := brandID.(int64); ok && brandIDInt != 0 {
			queryBuilder.WriteString(fmt.Sprintf(" AND brand_id = $%d", argCounter))
			args = append(args, brandIDInt)
			argCounter++
		}
	}

	// Фильтрация по is_upcoming
	if isUpcoming, ok := filter["is_upcoming"]; ok {
		if isUpcomingBool, ok := isUpcoming.(bool); ok {
			queryBuilder.WriteString(fmt.Sprintf(" AND is_upcoming = $%d", argCounter))
			args = append(args, isUpcomingBool)
			argCounter++
		}
	}

	// Фильтрация по is_limited
	if isLimited, ok := filter["is_limited"]; ok {
		if isLimitedBool, ok := isLimited.(bool); ok {
			queryBuilder.WriteString(fmt.Sprintf(" AND is_limited = $%d", argCounter))
			args = append(args, isLimitedBool)
			argCounter++
		}
	}

	if sortBy != "" {
		order := "ASC"
		sortField := sortBy
		if strings.HasPrefix(sortBy, "-") {
			order = "DESC"
			sortField = strings.TrimPrefix(sortBy, "-")
		}

		allowedSortFields := map[string]bool{
			"name":         true,
			"release_date": true,
			"brand_id":     true,
			"created_at":   true,
			"updated_at":   true,
		}

		if allowedSortFields[sortField] {
			queryBuilder.WriteString(fmt.Sprintf(" ORDER BY %s %s", sortField, order))
		} else {
			span.LogFields(
				log.String("event", "invalid_sort_field"),
				log.String("sort_by", sortBy),
			)
			r.log.Warn().
				Str("sort_by", sortBy).
				Msg("Invalid sort field provided, ignoring sorting")
			// Дефолтная сортировка
			queryBuilder.WriteString(" ORDER BY name ASC")
		}
	} else {
		// Дефолтная сортировка
		queryBuilder.WriteString(" ORDER BY name ASC")
	}

	finalQuery := queryBuilder.String()
	span.SetTag("db.query", finalQuery)
	span.SetTag("db.args_count", len(args))

	rows, err := r.pool.Query(ctx, finalQuery, args...)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(
			log.String("event", "db_query_error"),
			log.String("message", err.Error()),
			log.String("query", finalQuery),
		)
		r.log.Error().
			Err(err).
			Str("operation", "GetAll").
			Msg("Failed to execute GetAll query")
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var models []dto.Model

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
