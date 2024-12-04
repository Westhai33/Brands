package brand

import (
	"Brands/internal/dto"
	"context"
	"fmt"
	"strings"
)

// GetAll получает все бренды с сортировкой по имени
func (r *BrandRepository) GetAll(ctx context.Context, filter map[string]interface{}, sortBy string) ([]dto.Brand, error) {
	r.log.Info().Str("operation", "GetAll").Interface("filter", filter).Str("sort", sortBy).Msg("Fetching all brands")

	var queryBuilder strings.Builder
	var args []interface{}
	argCounter := 1

	queryBuilder.WriteString("SELECT id, name, link, description, logo_url, cover_image_url, founded_year, origin_country, popularity, is_premium, is_upcoming, created_at, updated_at FROM brands WHERE is_deleted = false")

	for key, value := range filter {
		if value == nil || value == "" {
			continue
		}
		switch key {
		case "name":
			queryBuilder.WriteString(fmt.Sprintf(" AND name ILIKE $%d", argCounter))
			args = append(args, "%"+value.(string)+"%")
		case "origin_country":
			queryBuilder.WriteString(fmt.Sprintf(" AND origin_country = $%d", argCounter))
			args = append(args, value.(string))
		case "popularity":
			queryBuilder.WriteString(fmt.Sprintf(" AND popularity = $%d", argCounter))
			args = append(args, value)
		case "is_premium":
			queryBuilder.WriteString(fmt.Sprintf(" AND is_premium = $%d", argCounter))
			args = append(args, value.(bool))
		case "is_upcoming":
			queryBuilder.WriteString(fmt.Sprintf(" AND is_upcoming = $%d", argCounter))
			args = append(args, value.(bool))
		case "founded_year":
			queryBuilder.WriteString(fmt.Sprintf(" AND founded_year = $%d", argCounter))
			args = append(args, value.(int))
		}
		argCounter++
	}

	if sortBy != "" {
		if strings.HasPrefix(sortBy, "-") {
			queryBuilder.WriteString(fmt.Sprintf(" ORDER BY %s DESC", strings.TrimPrefix(sortBy, "-")))
		} else {
			queryBuilder.WriteString(fmt.Sprintf(" ORDER BY %s ASC", sortBy))
		}
	} else {
		queryBuilder.WriteString(" ORDER BY created_at DESC")
	}

	r.log.Info().Str("query", queryBuilder.String()).Msg("Executing query with sorting")

	rows, err := r.pool.Query(ctx, queryBuilder.String(), args...)
	if err != nil {
		r.log.Error().Err(err).Msg("Failed to execute GetAll query")
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

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

	if rows.Err() != nil {
		r.log.Error().Err(rows.Err()).Msg("Error iterating over rows in GetAll")
		return nil, fmt.Errorf("error iterating over rows: %w", rows.Err())
	}

	r.log.Info().Int("brands_count", len(brands)).Msg("GetAll operation completed")
	return brands, nil
}
