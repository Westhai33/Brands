package repository

import (
	"Brands/internal/dto"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type BrandRepository struct {
	pool *pgxpool.Pool
}

func NewBrandRepository(pool *pgxpool.Pool) *BrandRepository {
	return &BrandRepository{pool: pool}
}

// Create создает новый бренд
func (r *BrandRepository) Create(ctx context.Context, brand *dto.Brand) (int64, error) {
	query := `
		INSERT INTO brands (name, link, description, logo_url, cover_image_url, founded_year, origin_country, popularity, is_premium, is_upcoming, created_at, updated_at, is_deleted)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id`
	now := time.Now()
	err := r.pool.QueryRow(ctx, query,
		brand.Name, brand.Link, brand.Description, brand.LogoURL, brand.CoverImageURL,
		brand.FoundedYear, brand.OriginCountry, brand.Popularity, brand.IsPremium,
		brand.IsUpcoming, now, now, false,
	).Scan(&brand.ID)
	return brand.ID, err
}

// GetByID получает бренд по ID
func (r *BrandRepository) GetByID(ctx context.Context, id int64) (*dto.Brand, error) {
	query := `
		SELECT id, name, link, description, logo_url, cover_image_url, founded_year, origin_country, popularity, is_premium, is_upcoming, is_deleted, created_at, updated_at
		FROM brands WHERE id = $1 AND is_deleted = false`
	brand := dto.Brand{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&brand.ID, &brand.Name, &brand.Link, &brand.Description, &brand.LogoURL,
		&brand.CoverImageURL, &brand.FoundedYear, &brand.OriginCountry, &brand.Popularity,
		&brand.IsPremium, &brand.IsUpcoming, &brand.IsDeleted, &brand.CreatedAt, &brand.UpdatedAt,
	)
	return &brand, err
}

// Update обновляет данные бренда
func (r *BrandRepository) Update(ctx context.Context, brand *dto.Brand) error {
	query := `
		UPDATE brands SET name = $2, link = $3, description = $4, logo_url = $5, cover_image_url = $6, 
		                   founded_year = $7, origin_country = $8, popularity = $9, is_premium = $10, 
		                   is_upcoming = $11, updated_at = $12
		WHERE id = $1 AND is_deleted = false`
	_, err := r.pool.Exec(ctx, query, brand.ID, brand.Name, brand.Link, brand.Description, brand.LogoURL,
		brand.CoverImageURL, brand.FoundedYear, brand.OriginCountry, brand.Popularity, brand.IsPremium,
		brand.IsUpcoming, time.Now())
	return err
}

// SoftDelete мягко удаляет бренд
func (r *BrandRepository) SoftDelete(ctx context.Context, id int64) error {
	query := `
		UPDATE brands SET is_deleted = true, updated_at = $2 WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id, time.Now())
	return err
}

// Restore восстанавливает мягко удалённый бренд
func (r *BrandRepository) Restore(ctx context.Context, id int64) error {
	query := `
		UPDATE brands SET is_deleted = false, updated_at = $2 WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id, time.Now())
	return err
}

// GetAll получает все бренды с фильтрацией и сортировкой
func (r *BrandRepository) GetAll(ctx context.Context, filter string, sort string) ([]dto.Brand, error) {
	query := `
		SELECT id, name, link, description, logo_url, cover_image_url, founded_year, origin_country, popularity, is_premium, is_upcoming, is_deleted, created_at, updated_at
		FROM brands
		WHERE is_deleted = false AND ($1 = '' OR name ILIKE '%' || $1 || '%')
		ORDER BY 
		    CASE WHEN $2 = 'popularity' THEN popularity END DESC,
		    CASE WHEN $2 = 'founded_year' THEN founded_year END DESC,
		    name ASC`
	rows, err := r.pool.Query(ctx, query, filter, sort)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var brands []dto.Brand
	for rows.Next() {
		var brand dto.Brand
		err := rows.Scan(
			&brand.ID, &brand.Name, &brand.Link, &brand.Description, &brand.LogoURL,
			&brand.CoverImageURL, &brand.FoundedYear, &brand.OriginCountry, &brand.Popularity,
			&brand.IsPremium, &brand.IsUpcoming, &brand.IsDeleted, &brand.CreatedAt, &brand.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		brands = append(brands, brand)
	}
	return brands, nil
}
