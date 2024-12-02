package repository

import (
	"Brands/internal/dto"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type ModelRepository struct {
	ctx  context.Context
	pool *pgxpool.Pool
	log  zerolog.Logger
}

func NewModelRepository(ctx context.Context, pool *pgxpool.Pool, logger zerolog.Logger) (*ModelRepository, error) {
	return &ModelRepository{ctx: ctx, pool: pool, log: logger}, nil
}

// BrandExists проверяет, существует ли бренд с заданным ID
func (r *ModelRepository) BrandExists(ctx context.Context, brandID int64) (bool, error) {
	r.log.Info().Int64("brand_id", brandID).Msg("Checking if brand exists")
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM brands WHERE id = $1 AND is_deleted = false)`
	err := r.pool.QueryRow(ctx, query, brandID).Scan(&exists)
	if err != nil {
		r.log.Error().Err(err).Int64("brand_id", brandID).Msg("Failed to check if brand exists")
	}
	return exists, err
}

// Create создает новую модель
func (r *ModelRepository) Create(ctx context.Context, model *dto.Model) (int64, error) {
	r.log.Info().Str("operation", "Create").Interface("model", model).Msg("Creating a new model")
	exists, err := r.BrandExists(ctx, model.BrandID)
	if err != nil {
		return 0, err
	}
	if !exists {
		err := fmt.Errorf("brand with ID %d does not exist", model.BrandID)
		r.log.Warn().Int64("brand_id", model.BrandID).Msg(err.Error())
		return 0, err
	}

	query := `
		INSERT INTO models (brand_id, name, release_date, is_upcoming, is_limited, created_at, updated_at, is_deleted)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`
	now := time.Now()
	err = r.pool.QueryRow(ctx, query,
		model.BrandID, model.Name, model.ReleaseDate, model.IsUpcoming,
		model.IsLimited, now, now, false,
	).Scan(&model.ID)
	if err != nil {
		r.log.Error().Err(err).Msg("Failed to create model")
		return 0, err
	}

	r.log.Info().Int64("model_id", model.ID).Msg("Model created successfully")
	return model.ID, nil
}

// GetByID получает модель по ID
func (r *ModelRepository) GetByID(ctx context.Context, id int64) (*dto.Model, error) {
	r.log.Info().Int64("model_id", id).Msg("Fetching model by ID")
	query := `
		SELECT id, brand_id, name, release_date, is_upcoming, is_limited, is_deleted, created_at, updated_at
		FROM models WHERE id = $1 AND is_deleted = false`
	model := dto.Model{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&model.ID, &model.BrandID, &model.Name, &model.ReleaseDate, &model.IsUpcoming,
		&model.IsLimited, &model.IsDeleted, &model.CreatedAt, &model.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			r.log.Warn().Int64("model_id", id).Msg("Model not found")
			return nil, nil
		}
		r.log.Error().Err(err).Int64("model_id", id).Msg("Failed to fetch model by ID")
		return nil, err
	}

	r.log.Info().Interface("model", model).Msg("Model fetched successfully")
	return &model, nil
}

// Update обновляет данные модели
func (r *ModelRepository) Update(ctx context.Context, model *dto.Model) error {
	r.log.Info().Int64("model_id", model.ID).Msg("Updating model")
	exists, err := r.BrandExists(ctx, model.BrandID)
	if err != nil {
		return err
	}
	if !exists {
		err := fmt.Errorf("brand with ID %d does not exist", model.BrandID)
		r.log.Warn().Int64("brand_id", model.BrandID).Msg(err.Error())
		return err
	}

	query := `
		UPDATE models SET brand_id = $2, name = $3, release_date = $4, is_upcoming = $5, 
		                 is_limited = $6, updated_at = $7
		WHERE id = $1 AND is_deleted = false`
	_, err = r.pool.Exec(ctx, query, model.ID, model.BrandID, model.Name, model.ReleaseDate,
		model.IsUpcoming, model.IsLimited, time.Now())
	if err != nil {
		r.log.Error().Err(err).Int64("model_id", model.ID).Msg("Failed to update model")
		return err
	}

	r.log.Info().Int64("model_id", model.ID).Msg("Model updated successfully")
	return nil
}

// SoftDelete мягко удаляет модель
func (r *ModelRepository) SoftDelete(ctx context.Context, id int64) error {
	r.log.Info().Int64("model_id", id).Msg("Soft deleting model")
	query := `
		UPDATE models SET is_deleted = true, updated_at = $2 WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id, time.Now())
	if err != nil {
		r.log.Error().Err(err).Int64("model_id", id).Msg("Failed to soft delete model")
	}
	return err
}

// Restore восстанавливает мягко удалённую модель
func (r *ModelRepository) Restore(ctx context.Context, id int64) error {
	r.log.Info().Int64("model_id", id).Msg("Restoring model")
	query := `
		UPDATE models SET is_deleted = false, updated_at = $2 WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id, time.Now())
	if err != nil {
		r.log.Error().Err(err).Int64("model_id", id).Msg("Failed to restore model")
	}
	return err
}

// GetAll получает все модели с фильтрацией и сортировкой
func (r *ModelRepository) GetAll(ctx context.Context, filter map[string]interface{}, sortBy string) ([]dto.Model, error) {
	r.log.Info().Interface("filter", filter).Str("sort", sortBy).Msg("Fetching all models")

	var queryBuilder strings.Builder
	var args []interface{}
	argCounter := 1

	queryBuilder.WriteString("SELECT id, brand_id, name, release_date, is_upcoming, is_limited, is_deleted, created_at, updated_at FROM models WHERE is_deleted = false")

	if name, ok := filter["name"]; ok && name != "" {
		queryBuilder.WriteString(fmt.Sprintf(" AND name ILIKE $%d", argCounter))
		args = append(args, "%"+name.(string)+"%")
		argCounter++
	}

	if brandID, ok := filter["brand_id"]; ok && brandID != 0 {
		queryBuilder.WriteString(fmt.Sprintf(" AND brand_id = $%d", argCounter))
		args = append(args, brandID)
		argCounter++
	}

	if isUpcoming, ok := filter["is_upcoming"]; ok {
		if isUpcomingBool, ok := isUpcoming.(bool); ok {
			queryBuilder.WriteString(fmt.Sprintf(" AND is_upcoming = $%d", argCounter))
			args = append(args, isUpcomingBool)
			argCounter++
		}
	}

	if isLimited, ok := filter["is_limited"]; ok {
		if isLimitedBool, ok := isLimited.(bool); ok {
			queryBuilder.WriteString(fmt.Sprintf(" AND is_limited = $%d", argCounter))
			args = append(args, isLimitedBool)
			argCounter++
		}
	}

	if sortBy != "" {
		if strings.HasPrefix(sortBy, "-") {
			queryBuilder.WriteString(fmt.Sprintf(" ORDER BY %s DESC", strings.TrimPrefix(sortBy, "-")))
		} else {
			queryBuilder.WriteString(fmt.Sprintf(" ORDER BY %s ASC", sortBy))
		}
	} else {
		queryBuilder.WriteString(" ORDER BY name ASC")
	}

	rows, err := r.pool.Query(ctx, queryBuilder.String(), args...)
	if err != nil {
		r.log.Error().Err(err).Msg("Failed to execute GetAll query")
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var models []dto.Model
	for rows.Next() {
		var model dto.Model
		if err := rows.Scan(
			&model.ID, &model.BrandID, &model.Name, &model.ReleaseDate, &model.IsUpcoming, &model.IsLimited,
			&model.IsDeleted, &model.CreatedAt, &model.UpdatedAt,
		); err != nil {
			r.log.Error().Err(err).Msg("Failed to scan row in GetAll")
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		models = append(models, model)
	}

	if rows.Err() != nil {
		r.log.Error().Err(rows.Err()).Msg("Error iterating over rows in GetAll")
		return nil, fmt.Errorf("error iterating over rows: %w", rows.Err())
	}

	r.log.Info().Int("model_count", len(models)).Msg("GetAll models completed successfully")
	return models, nil
}
