package brand

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

// SoftDelete мягко удаляет бренд
func (r *BrandRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BrandRepository.SoftDelete")
	defer span.Finish()

	query := `UPDATE brands SET is_deleted = true, updated_at = NOW() WHERE id = $1`
	cmdTag, err := r.pool.Exec(ctx, query, id)

	if err != nil {
		span.LogFields(log.Error(err))
		if errors.Is(err, pgx.ErrNoRows) {
			r.log.Warn().Str("brand_id", id.String()).Msg("Brand not found")
			return ErrBrandNotFound
		}
		r.log.Error().Err(err).Str("brand_id", id.String()).Msg("Failed to soft delete brand")
		return fmt.Errorf("unable to soft delete brand: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		span.LogFields(log.Error(ErrBrandNotFound))
		r.log.Warn().
			Interface("brand_id", id.String()).
			Msg("No brand found to update")
		return ErrBrandNotFound
	}
	return nil
}

// Restore восстанавливает мягко удалённый бренд
func (r *BrandRepository) Restore(ctx context.Context, id uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BrandRepository.Restore")
	defer span.Finish()

	query := `UPDATE brands SET is_deleted = false, updated_at = NOW() WHERE id = $1`
	cmdTag, err := r.pool.Exec(ctx, query, id)

	if err != nil {
		span.LogFields(log.Error(err))
		if errors.Is(err, pgx.ErrNoRows) {
			r.log.Warn().Str("brand_id", id.String()).Msg("Brand not found")
			return ErrBrandNotFound
		}
		r.log.Error().Err(err).Str("brand_id", id.String()).Msg("Failed to restore brand")
		return fmt.Errorf("unable to restore brand: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		span.LogFields(log.Error(ErrBrandNotFound))
		r.log.Warn().
			Interface("brand_id", id.String()).
			Msg("No brand found to update")
		return ErrBrandNotFound
	}
	return nil
}
