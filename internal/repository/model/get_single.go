package model

import (
	"Brands/internal/dto"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// GetByID получает модель по ID
func (r *ModelRepository) GetByID(ctx context.Context, id int64) (*dto.Model, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ModelRepository.GetByID")
	defer span.Finish()

	query := `
		SELECT id, brand_id, name, release_date, is_upcoming, is_limited, is_deleted, created_at, updated_at
		FROM models WHERE id = $1 AND is_deleted = false`

	model := dto.Model{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&model.ID,
		&model.BrandID,
		&model.Name,
		&model.ReleaseDate,
		&model.IsUpcoming,
		&model.IsLimited,
		&model.IsDeleted,
		&model.CreatedAt,
		&model.UpdatedAt,
	)

	if err != nil {
		span.SetTag("error", true)
		span.LogFields(log.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			r.log.Warn().Int64("model_id", id).Msg("Model not found")
			return nil, ErrModelNotFound
		}

		r.log.Error().Err(err).Int64("model_id", id).Msg("Failed to fetch model by ID")
		return nil, fmt.Errorf("failed to fetch model by ID: %w", err)
	}

	return &model, nil
}
