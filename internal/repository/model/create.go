package model

import (
	"context"
	"fmt"

	"Brands/internal/dto"
	"Brands/internal/repository/brand"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// Create создает новую модель
func (r *ModelRepository) Create(ctx context.Context, model *dto.Model) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ModelRepository.Create")
	defer span.Finish()

	exists, err := r.brandExists(ctx, model.BrandID)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(log.Error(err))
		r.log.Error().Err(err).Int64("brand_id", model.BrandID).Msg("Failed to check if brand exists")
		return 0, fmt.Errorf("failed to check brand existence: %w", err)
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
		return 0, err
	}

	query := `
		INSERT INTO models (brand_id, name, release_date, is_upcoming, is_limited, created_at, updated_at, is_deleted)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW(), false)
		RETURNING id`
	err = r.pool.QueryRow(ctx, query,
		model.BrandID,
		model.Name,
		model.ReleaseDate,
		model.IsUpcoming,
		model.IsLimited,
	).Scan(&model.ID)

	// Обрабатываем ошибку запроса
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(log.Error(err))
		r.log.Error().Err(err).Msg("Failed to create model")
		return 0, fmt.Errorf("failed to create model: %w", err)
	}
	return model.ID, nil
}
