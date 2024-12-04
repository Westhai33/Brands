package model

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// brandExists проверяет, существует ли бренд с заданным ID
func (r *ModelRepository) brandExists(ctx context.Context, brandID int64) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ModelRepository.brandExists")
	defer span.Finish()

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM brands WHERE id = $1 AND is_deleted = false)`
	err := r.pool.QueryRow(ctx, query, brandID).Scan(&exists)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(log.Error(err))
		r.log.Error().Err(err).Int64("brand_id", brandID).Msg("Failed to check if brand exists")
		return false, fmt.Errorf("unable to check if brand exists: %w", err)
	}

	return exists, nil
}
