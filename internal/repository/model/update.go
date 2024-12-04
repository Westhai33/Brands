package model

import (
	"Brands/internal/dto"
	"Brands/internal/repository/brand"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// Update обновляет данные модели
func (r *ModelRepository) Update(ctx context.Context, model *dto.Model) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ModelRepository.Update")
	defer span.Finish()

	exists, err := r.brandExists(ctx, model.BrandID)
	if err != nil {

		span.LogFields(
			log.Error(err),
			log.Object("model", model),
		)
		r.log.Warn().Interface("model", model).Msgf("Failed to check if brand exists: %s", err.Error())
		return err
	}

	if !exists {
		err = fmt.Errorf(
			"brand with ID %d does not exist: %w",
			model.BrandID,
			brand.ErrBrandNotFound,
		)

		span.LogFields(
			log.Error(err),
			log.Object("model", model),
		)
		r.log.Warn().Interface("model", model).Msg(err.Error())
		return err
	}

	query := `
		UPDATE models 
		SET brand_id = @brand_id, 
		    name = @name, 
		    release_date = @release_date, 
		    is_upcoming = @is_upcoming, 
		    is_limited = @is_limited, 
		    updated_at = NOW()
		WHERE id = @id AND is_deleted = false
	`

	args := pgx.NamedArgs{
		"id":           model.ID,
		"brand_id":     model.BrandID,
		"name":         model.Name,
		"release_date": model.ReleaseDate,
		"is_upcoming":  model.IsUpcoming,
		"is_limited":   model.IsLimited,
	}
	_, err = r.pool.Exec(ctx, query, args)
	if err != nil {
		span.LogFields(log.Error(err))
		r.log.Error().Err(err).Interface("model", model).Msg("Failed to update model")
		return err
	}
	return nil
}
