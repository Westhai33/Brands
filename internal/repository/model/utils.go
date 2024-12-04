package model

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// brandExists проверяет, существует ли бренд с заданным ID
func (r *ModelRepository) brandExists(ctx context.Context, brandID uuid.UUID) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ModelRepository.brandExists")
	defer span.Finish()

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM brands WHERE id = $1 AND is_deleted = false)`
	err := r.pool.QueryRow(ctx, query, brandID).Scan(&exists)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(log.Error(err), log.String("brand_id", brandID.String()))
		r.log.Error().Err(err).Str("brand_id", brandID.String()).Msg("Failed to check if brand exists")
		return false, fmt.Errorf("unable to check if brand (%s) exists: %w", brandID, err)
	}

	return exists, nil
}
