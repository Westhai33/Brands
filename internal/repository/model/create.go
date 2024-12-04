package model

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"

	"Brands/internal/dto"
	"Brands/internal/repository/brand"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// Create создает новую модель
func (r *ModelRepository) Create(ctx context.Context, model *dto.Model) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ModelRepository.Create")
	defer span.Finish()

	exists, err := r.brandExists(ctx, model.BrandID)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(
			log.Error(err),
			log.Object("model", model),
		)
		r.log.Error().Err(err).Interface("model", model).Msg("Failed to check if brand exists")
		return fmt.Errorf("failed to check brand existence: %w", err)
	}
	if !exists {
		err = fmt.Errorf(
			"brand with ID %d does not exist: %w",
			model.BrandID,
			brand.ErrBrandNotFound,
		)
		span.SetTag("error", true)
		span.LogFields(
			log.Error(err),
			log.Object("model", model),
		)
		r.log.Warn().Interface("model", model).Msg(err.Error())
		return err
	}

	query := `
		INSERT INTO models (id, brand_id, name, release_date, is_upcoming, is_limited, created_at, updated_at, is_deleted)
		VALUES (:id, :name, :release_date, :is_upcoming, :is_limited, NOW(), NOW(), false)
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
		span.SetTag("error", true)
		span.LogFields(
			log.Error(err),
			log.Object("brand", model),
		)
		r.log.Error().Interface("model", model).Err(err).Msg("Failed to create model")
		return fmt.Errorf("failed to create model: %w", err)
	}
	return nil
}
