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

// Create создает новую модель
func (r *ModelRepository) Create(ctx context.Context, model *dto.Model) error {
	// Создаем span для трассировки
	span, ctx := opentracing.StartSpanFromContext(ctx, "ModelRepository.Create")
	defer span.Finish()

	// Проверяем существует ли бренд с указанным ID
	exists, err := r.brandExists(ctx, model.BrandID)
	if err != nil {
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
			brand.ErrBrandNotFound, // Убедитесь, что это ошибка правильно определена в пакете brand
		)

		span.LogFields(
			log.Error(err),
			log.Object("model", model),
		)
		r.log.Warn().Interface("model", model).Msg(err.Error())
		return err
	}

	// SQL-запрос для вставки новой модели
	query := `
		INSERT INTO models (id, brand_id, name, release_date, is_upcoming, is_limited, created_at, updated_at, is_deleted)
		VALUES (:id, :brand_id, :name, :release_date, :is_upcoming, :is_limited, NOW(), NOW(), false)
	`

	// Используем pgx.NamedArgs для передачи параметров запроса
	args := pgx.NamedArgs{
		"id":           model.ID,
		"brand_id":     model.BrandID,
		"name":         model.Name,
		"release_date": model.ReleaseDate,
		"is_upcoming":  model.IsUpcoming,
		"is_limited":   model.IsLimited,
	}

	// Выполняем запрос
	_, err = r.pool.Exec(ctx, query, args)
	if err != nil {
		// Логируем ошибку при создании модели
		span.LogFields(
			log.Error(err),
			log.Object("model", model),
		)
		r.log.Error().Interface("model", model).Err(err).Msg("Failed to create model")

		// Возвращаем ошибку
		return fmt.Errorf("failed to create model: %w", err)
	}

	return nil
}
