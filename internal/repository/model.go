package repository

import (
	"Brands/internal/dto"
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ModelRepository struct {
	ctx  context.Context
	pool *pgxpool.Pool
	log  zerolog.Logger
}

func NewModelRepository(ctx context.Context, pool *pgxpool.Pool, logger zerolog.Logger) (*ModelRepository, error) {
	c := &ModelRepository{ctx: ctx, pool: pool, log: logger}
	return c, nil
}

// Проверяет, существует ли бренд с заданным ID
func (r *ModelRepository) BrandExists(ctx context.Context, brandID int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM brands WHERE id = $1 AND is_deleted = false)`
	err := r.pool.QueryRow(ctx, query, brandID).Scan(&exists)
	return exists, err
}

// Create создает новую модель
func (r *ModelRepository) Create(ctx context.Context, model *dto.Model) (int64, error) {
	exists, err := r.BrandExists(ctx, model.BrandID)
	if err != nil {
		return 0, err
	}
	if !exists {
		return 0, fmt.Errorf("brand with ID %d does not exist", model.BrandID)
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
	return model.ID, err
}

// GetByID получает модель по ID
func (r *ModelRepository) GetByID(ctx context.Context, id int64) (*dto.Model, error) {
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
			return nil, nil
		}
		return nil, err
	}
	return &model, nil
}

// Update обновляет данные модели
func (r *ModelRepository) Update(ctx context.Context, model *dto.Model) error {
	exists, err := r.BrandExists(ctx, model.BrandID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("brand with ID %d does not exist", model.BrandID)
	}

	query := `
		UPDATE models SET brand_id = $2, name = $3, release_date = $4, is_upcoming = $5, 
		                 is_limited = $6, updated_at = $7
		WHERE id = $1 AND is_deleted = false`
	_, err = r.pool.Exec(ctx, query, model.ID, model.BrandID, model.Name, model.ReleaseDate,
		model.IsUpcoming, model.IsLimited, time.Now())
	return err
}

// SoftDelete мягко удаляет модель
func (r *ModelRepository) SoftDelete(ctx context.Context, id int64) error {
	query := `
		UPDATE models SET is_deleted = true, updated_at = $2 WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id, time.Now())
	return err
}

// Restore восстанавливает мягко удалённую модель
func (r *ModelRepository) Restore(ctx context.Context, id int64) error {
	query := `
		UPDATE models SET is_deleted = false, updated_at = $2 WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id, time.Now())
	return err
}

// GetAll получает все модели с фильтрацией и сортировкой
func (r *ModelRepository) GetAll(ctx context.Context, filter string, sort string, order string) ([]dto.Model, error) {
	baseQuery := `
		SELECT id, brand_id, name, release_date, is_upcoming, is_limited, is_deleted, created_at, updated_at
		FROM models
		WHERE is_deleted = false`

	if filter != "" {
		baseQuery += " AND name ILIKE '%' || $1 || '%'"
	}

	orderClause := " ORDER BY name ASC"

	if sort != "" {
		switch sort {
		case "release_date":
			if order == "desc" {
				orderClause = " ORDER BY release_date DESC"
			} else {
				orderClause = " ORDER BY release_date ASC"
			}
		case "name":
			if order == "desc" {
				orderClause = " ORDER BY name DESC"
			} else {
				orderClause = " ORDER BY name ASC"
			}
		}
	}

	finalQuery := baseQuery + orderClause

	var rows pgx.Rows
	var err error
	if filter != "" {
		rows, err = r.pool.Query(ctx, finalQuery, filter)
	} else {
		rows, err = r.pool.Query(ctx, finalQuery)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []dto.Model
	for rows.Next() {
		var model dto.Model
		err := rows.Scan(
			&model.ID, &model.BrandID, &model.Name, &model.ReleaseDate, &model.IsUpcoming,
			&model.IsLimited, &model.IsDeleted, &model.CreatedAt, &model.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		models = append(models, model)
	}
	return models, nil
}
