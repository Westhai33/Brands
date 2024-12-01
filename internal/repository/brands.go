package repository

import (
	"Brands/internal/dto"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type BrandRepository struct {
	ctx  context.Context
	pool *pgxpool.Pool
	log  zerolog.Logger
}

func NewBrandRepository(ctx context.Context, pool *pgxpool.Pool, logger zerolog.Logger) (*BrandRepository, error) {
	return &BrandRepository{ctx: ctx, pool: pool, log: logger}, nil
}

// Create создает новый бренд
func (r *BrandRepository) Create(ctx context.Context, brand *dto.Brand) (int64, error) {
	r.log.Info().Str("operation", "Create").Str("brand_name", brand.Name).Msg("Creating a new brand")

	query := `
		INSERT INTO public.brands (name, link, description, logo_url, cover_image_url, founded_year, origin_country, popularity, is_premium, is_upcoming, created_at, updated_at, is_deleted)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id`
	now := time.Now()
	err := r.pool.QueryRow(ctx, query,
		brand.Name, brand.Link, brand.Description, brand.LogoURL, brand.CoverImageURL,
		brand.FoundedYear, brand.OriginCountry, brand.Popularity, brand.IsPremium,
		brand.IsUpcoming, now, now, false,
	).Scan(&brand.ID)
	if err != nil {
		r.log.Error().Err(err).Msg("Failed to create brand")
		return 0, fmt.Errorf("unable to insert brand: %w", err)
	}

	r.log.Info().Int64("brand_id", brand.ID).Msg("Brand created successfully")
	return brand.ID, nil
}

// GetByID получает бренд по ID
func (r *BrandRepository) GetByID(ctx context.Context, id int64) (*dto.Brand, error) {
	r.log.Info().Str("operation", "GetByID").Int64("brand_id", id).Msg("Fetching brand by ID")

	brand := &dto.Brand{}
	row := r.pool.QueryRow(ctx, `
        SELECT id, name, link, description, logo_url, cover_image_url, founded_year, 
               origin_country, popularity, is_premium, is_upcoming, is_deleted, created_at, updated_at
        FROM brands
        WHERE id = $1
    `, id)

	err := row.Scan(&brand.ID, &brand.Name, &brand.Link, &brand.Description, &brand.LogoURL, &brand.CoverImageURL,
		&brand.FoundedYear, &brand.OriginCountry, &brand.Popularity, &brand.IsPremium, &brand.IsUpcoming,
		&brand.IsDeleted, &brand.CreatedAt, &brand.UpdatedAt)
	if err != nil {
		r.log.Error().Err(err).Int64("brand_id", id).Msg("Failed to fetch brand by ID")
		return nil, err
	}

	if brand.IsDeleted {
		err := fmt.Errorf("brand has been soft-deleted")
		r.log.Warn().Int64("brand_id", id).Msg("Brand is soft-deleted")
		return nil, err
	}

	r.log.Info().Interface("brand", brand).Msg("Brand fetched successfully")
	return brand, nil
}

// Update обновляет данные бренда
func (r *BrandRepository) Update(ctx context.Context, brand *dto.Brand) error {
	r.log.Info().Str("operation", "Update").Int64("brand_id", brand.ID).Msg("Updating brand")

	query := `
		UPDATE brands SET name = $2, link = $3, description = $4, logo_url = $5, cover_image_url = $6, 
		                   founded_year = $7, origin_country = $8, popularity = $9, is_premium = $10, 
		                   is_upcoming = $11, updated_at = $12
		WHERE id = $1 AND is_deleted = false`
	_, err := r.pool.Exec(ctx, query, brand.ID, brand.Name, brand.Link, brand.Description, brand.LogoURL,
		brand.CoverImageURL, brand.FoundedYear, brand.OriginCountry, brand.Popularity, brand.IsPremium,
		brand.IsUpcoming, time.Now())
	if err != nil {
		r.log.Error().Err(err).Int64("brand_id", brand.ID).Msg("Failed to update brand")
		return err
	}

	r.log.Info().Int64("brand_id", brand.ID).Msg("Brand updated successfully")
	return nil
}

// SoftDelete мягко удаляет бренд
func (r *BrandRepository) SoftDelete(ctx context.Context, id int64) error {
	r.log.Info().Str("operation", "SoftDelete").Int64("brand_id", id).Msg("Soft deleting brand")

	query := `UPDATE brands SET is_deleted = true, updated_at = $2 WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id, time.Now())
	if err != nil {
		r.log.Error().Err(err).Int64("brand_id", id).Msg("Failed to soft delete brand")
		return err
	}

	r.log.Info().Int64("brand_id", id).Msg("Brand soft deleted successfully")
	return nil
}

// Restore восстанавливает мягко удалённый бренд
func (r *BrandRepository) Restore(ctx context.Context, id int64) error {
	r.log.Info().Str("operation", "Restore").Int64("brand_id", id).Msg("Restoring brand")

	query := `UPDATE brands SET is_deleted = false, updated_at = $2 WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id, time.Now())
	if err != nil {
		r.log.Error().Err(err).Int64("brand_id", id).Msg("Failed to restore brand")
		return err
	}

	r.log.Info().Int64("brand_id", id).Msg("Brand restored successfully")
	return nil
}

// GetAll получает все бренды с фильтрацией и сортировкой
func (r *BrandRepository) GetAll(ctx context.Context, filter map[string]interface{}, sortBy string) ([]dto.Brand, error) {
	r.log.Info().Str("operation", "GetAll").Interface("filter", filter).Str("sort", sortBy).Msg("Fetching all brands")

	var queryBuilder strings.Builder
	var args []interface{}
	argCounter := 1

	// Начинаем строить запрос с условием is_deleted = false
	queryBuilder.WriteString("SELECT id, name, link, description, logo_url, cover_image_url, founded_year, origin_country, popularity, is_premium, is_upcoming, created_at, updated_at FROM brands WHERE is_deleted = false")

	// Добавляем фильтры в запрос
	if name, ok := filter["name"]; ok && name != "" {
		queryBuilder.WriteString(fmt.Sprintf(" AND name ILIKE $%d", argCounter))
		args = append(args, "%"+name.(string)+"%")
		argCounter++
	}

	if originCountry, ok := filter["origin_country"]; ok && originCountry != "" {
		queryBuilder.WriteString(fmt.Sprintf(" AND origin_country = $%d", argCounter))
		args = append(args, originCountry.(string))
		argCounter++
	}

	if popularity, ok := filter["popularity"]; ok && popularity != 0 {
		queryBuilder.WriteString(fmt.Sprintf(" AND popularity = $%d", argCounter))
		args = append(args, popularity)
		argCounter++
	}

	if isPremium, ok := filter["is_premium"]; ok {
		if isPremiumBool, ok := isPremium.(bool); ok {
			queryBuilder.WriteString(fmt.Sprintf(" AND is_premium = $%d", argCounter))
			args = append(args, isPremiumBool)
			argCounter++
		}
	}

	// Добавляем сортировку
	if sortBy != "" {
		if strings.HasPrefix(sortBy, "-") {
			queryBuilder.WriteString(fmt.Sprintf(" ORDER BY %s DESC", strings.TrimPrefix(sortBy, "-")))
		} else {
			queryBuilder.WriteString(fmt.Sprintf(" ORDER BY %s ASC", sortBy))
		}
	}

	// Выполняем запрос
	rows, err := r.pool.Query(ctx, queryBuilder.String(), args...)
	if err != nil {
		r.log.Error().Err(err).Msg("Failed to execute GetAll query")
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	// Считываем результаты
	var brands []dto.Brand
	for rows.Next() {
		var brand dto.Brand
		if err := rows.Scan(
			&brand.ID, &brand.Name, &brand.Link, &brand.Description, &brand.LogoURL, &brand.CoverImageURL,
			&brand.FoundedYear, &brand.OriginCountry, &brand.Popularity, &brand.IsPremium, &brand.IsUpcoming,
			&brand.CreatedAt, &brand.UpdatedAt,
		); err != nil {
			r.log.Error().Err(err).Msg("Failed to scan row in GetAll")
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		brands = append(brands, brand)
	}

	// Проверяем на ошибки чтения строк
	if rows.Err() != nil {
		r.log.Error().Err(rows.Err()).Msg("Error iterating over rows in GetAll")
		return nil, fmt.Errorf("error iterating over rows: %w", rows.Err())
	}

	r.log.Info().Int("brands_count", len(brands)).Msg("GetAll operation completed")
	return brands, nil
}
