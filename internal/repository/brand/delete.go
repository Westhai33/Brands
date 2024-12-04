package brand

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

// SoftDelete мягко удаляет бренд
func (r *BrandRepository) SoftDelete(ctx context.Context, id int64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BrandRepository.SoftDelete")
	defer span.Finish()

	query := `UPDATE brands SET is_deleted = true, updated_at = NOW() WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)

	if err != nil {
		span.SetTag("error", true)
		span.LogFields(log.Error(err))
		r.log.Error().Err(err).Int64("brand_id", id).Msg("Failed to soft delete brand")
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrBrandNotFound
		}
		return fmt.Errorf("unable to soft delete brand: %w", err)
	}
	return nil
}

// Restore восстанавливает мягко удалённый бренд
func (r *BrandRepository) Restore(ctx context.Context, id int64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BrandRepository.Restore")
	defer span.Finish()

	query := `UPDATE brands SET is_deleted = false, updated_at = NOW() WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)

	if err != nil {
		span.SetTag("error", true)
		span.LogFields(log.Error(err))
		r.log.Error().Err(err).Int64("brand_id", id).Msg("Failed to restore brand")

		if errors.Is(err, pgx.ErrNoRows) {
			return ErrBrandNotFound
		}
		return fmt.Errorf("unable to restore brand: %w", err)
	}
	return nil
}
