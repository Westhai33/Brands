package model

import (
	"Brands/internal/dto"
	"Brands/internal/repository/brand"
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// Update обновляет данные модели
func (r *Repository) Update(ctx context.Context, model *dto.Model) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ModelRepository.Update")
	defer span.Finish()

	exists, err := r.brandExists(ctx, model.BrandID)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(log.Error(err))
		r.log.Error().Err(err).Int64("brand_id", model.BrandID).Msg("Failed to check if brand exists")
		return err
	}

	if !exists {
		err = fmt.Errorf(
			"brand with ID %d does not exist: %w",
			model.BrandID,
			brand.ErrBrandNotFound,
		)
		span.SetTag("error", true)
		span.LogFields(log.Error(err))
		r.log.Warn().Int64("brand_id", model.BrandID).Msg(err.Error())
		return err
	}

	query := `
		UPDATE models SET brand_id = $2, name = $3, release_date = $4, is_upcoming = $5, 
		                 is_limited = $6, updated_at = NOW()
		WHERE id = $1 AND is_deleted = false`

	_, err = r.pool.Exec(ctx, query,
		model.ID,
		model.BrandID,
		model.Name,
		model.ReleaseDate,
		model.IsUpcoming,
		model.IsLimited,
	)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(log.Error(err))
		r.log.Error().Err(err).Int64("model_id", model.ID).Msg("Failed to update model")
		return err
	}
	return nil
}
