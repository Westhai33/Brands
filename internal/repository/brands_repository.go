package repository

import (
	"Brands/internal/dto"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
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
		return nil, err
	}

	if brand.IsDeleted {
		return nil, fmt.Errorf("brand has been soft-deleted") // Возвращаем ошибку, если бренд удалён
	}

	return brand, nil
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
func (r *BrandRepository) GetAll(ctx context.Context, filter map[string]interface{}, sortBy string) ([]dto.Brand, error) {
	var queryBuilder strings.Builder
	var args []interface{}
	argCounter := 1

	// Начинаем строить запрос
	queryBuilder.WriteString("SELECT id, name, link, description, logo_url, cover_image_url, founded_year, origin_country, popularity, is_premium, is_upcoming, created_at, updated_at FROM brands WHERE 1=1")

	// Логируем фильтры
	fmt.Printf("Filter params: %v\n", filter)

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

	// Логируем сгенерированный запрос и аргументы
	fmt.Printf("Generated query: %s\n", queryBuilder.String())
	fmt.Printf("Query args: %v\n", args)

	// Выполняем запрос
	rows, err := r.pool.Query(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	// Преобразуем строки в бренды
	var brands []dto.Brand
	for rows.Next() {
		var brand dto.Brand
		// Исключаем is_deleted из метода rows.Scan
		if err := rows.Scan(
			&brand.ID,
			&brand.Name,
			&brand.Link,
			&brand.Description,
			&brand.LogoURL,
			&brand.CoverImageURL,
			&brand.FoundedYear,
			&brand.OriginCountry,
			&brand.Popularity,
			&brand.IsPremium,
			&brand.IsUpcoming,
			&brand.CreatedAt,
			&brand.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		brands = append(brands, brand)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	// Если бренды не найдены, возвращаем пустой массив
	if len(brands) == 0 {
		fmt.Println("No brands found with the given filters")
	}

	return brands, nil
}
