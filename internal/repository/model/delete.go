package model

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

// SoftDelete мягко удаляет модель
func (r *ModelRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ModelRepository.SoftDelete")
	defer span.Finish()

	query := `UPDATE models SET is_deleted = true, updated_at = NOW() WHERE id = $1`
	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		span.LogFields(log.Error(err))
		if errors.Is(err, pgx.ErrNoRows) {
			r.log.Warn().Str("model_id", id.String()).Msg("Model not found")
			return ErrModelNotFound
		}
		r.log.Error().Err(err).Str("model_id", id.String()).Msg("Failed to soft delete model")
		return fmt.Errorf("failed to soft delete model with id %d: %w", id, err)
	}
	if cmdTag.RowsAffected() == 0 {
		span.LogFields(log.Error(ErrModelNotFound))
		r.log.Warn().
			Interface("model_id", id.String()).
			Msg("No model found to update")
		return ErrModelNotFound
	}
	return nil
}

// Restore восстанавливает мягко удалённую модель
func (r *ModelRepository) Restore(ctx context.Context, id uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ModelRepository.Restore")
	defer span.Finish()

	query := `UPDATE models SET is_deleted = false, updated_at = NOW() WHERE id = $1`
	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		span.LogFields(log.Error(err))
		if errors.Is(err, pgx.ErrNoRows) {
			r.log.Warn().Str("model_id", id.String()).Msg("Model not found")
			return ErrModelNotFound
		}
		r.log.Error().Err(err).Str("model_id", id.String()).Msg("Failed to restore model")
		return fmt.Errorf("failed to restore model with id %d: %w", id, err)
	}
	if cmdTag.RowsAffected() == 0 {
		span.LogFields(log.Error(ErrModelNotFound))
		r.log.Warn().
			Interface("model_id", id.String()).
			Msg("No model found to update")
		return ErrModelNotFound
	}
	return nil
}
