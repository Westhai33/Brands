package model

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

// SoftDelete мягко удаляет модель
func (r *ModelRepository) SoftDelete(ctx context.Context, id int64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ModelRepository.SoftDelete")
	defer span.Finish()

	query := `UPDATE models SET is_deleted = true, updated_at = $2 WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id, time.Now())
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(log.Error(err))
		r.log.Error().Err(err).Int64("model_id", id).Msg("Failed to soft delete model")
		return fmt.Errorf("failed to soft delete model with id %d: %w", id, err)
	}
	r.log.Info().Int64("model_id", id).Msg("Model soft deleted successfully")
	return nil
}

// Restore восстанавливает мягко удалённую модель
func (r *ModelRepository) Restore(ctx context.Context, id int64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ModelRepository.Restore")
	defer span.Finish()

	query := `UPDATE models SET is_deleted = false, updated_at = $2 WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id, time.Now())
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(log.Error(err))
		r.log.Error().Err(err).Int64("model_id", id).Msg("Failed to restore model")
		return fmt.Errorf("failed to restore model with id %d: %w", id, err)
	}

	r.log.Info().Int64("model_id", id).Msg("Model restored successfully")
	return nil
}
