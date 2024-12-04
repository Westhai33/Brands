package model

import (
	"Brands/internal/dto"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// GetByID получает модель по ID
func (r *ModelRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*dto.Model, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ModelRepository.GetByID")
	defer span.Finish()

	query := `
		SELECT *
		FROM models 
		WHERE id = $1 AND is_deleted = false
	`
	row, err := r.pool.Query(ctx, query, id)
	if err != nil {
		span.LogFields(log.Error(err))
		r.log.Error().Err(err).Str("model_id", id.String()).Msg("Failed to fetch model by ID")
		return nil, fmt.Errorf("unable to get model by id: %w", err)
	}

	var models []dto.Model
	models, err = pgx.CollectRows(row, pgx.RowToStructByName[dto.Model])
	if err != nil {
		span.LogFields(log.Error(err))
		r.log.Error().Err(err).Str("model_id", id.String()).Msg("Failed to collect rows")
		return nil, fmt.Errorf("unable to collect rows: %w", err)
	}

	if len(models) == 0 {
		r.log.Warn().Str("brand_id", id.String()).Msg("Brand not found")
		return nil, ErrModelNotFound
	}

	if len(models) > 1 {
		err = fmt.Errorf("multiple brands found with id %s", id.String())
		span.LogFields(log.Error(err))
		r.log.Error().Str("brand_id", id.String()).Msg("Multiple brands found with the same ID")
		return nil, err
	}

	return &models[0], nil
}
