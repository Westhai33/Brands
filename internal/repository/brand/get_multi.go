package brand

import (
	"Brands/internal/dto"
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"strings"
)

// GetAll получает все бренды с фильтрацией и сортировкой
func (r *Repository) GetAll(
	ctx context.Context,
	filter map[string]interface{},
	sortBy string,
) ([]dto.Brand, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BrandRepository.GetAll")
	defer span.Finish()

	var queryBuilder strings.Builder
	var args []interface{}
	argCounter := 1

	baseQuery := `
        SELECT id, name, link, description, logo_url, cover_image_url, founded_year, 
               origin_country, popularity, is_premium, is_upcoming, created_at, updated_at 
        FROM brands 
        WHERE is_deleted = false
    `
	queryBuilder.WriteString(baseQuery)

	// Фильтрация по имени
	if name, ok := filter["name"]; ok && name != "" {
		queryBuilder.WriteString(fmt.Sprintf(" AND name ILIKE $%d", argCounter))
		args = append(args, "%"+name.(string)+"%")
		argCounter++
	}

	// Фильтрация по стране происхождения
	if originCountry, ok := filter["origin_country"]; ok && originCountry != "" {
		queryBuilder.WriteString(fmt.Sprintf(" AND origin_country = $%d", argCounter))
		args = append(args, originCountry.(string))
		argCounter++
	}

	// Фильтрация по популярности
	if popularity, ok := filter["popularity"]; ok && popularity != 0 {
		queryBuilder.WriteString(fmt.Sprintf(" AND popularity = $%d", argCounter))
		args = append(args, popularity)
		argCounter++
	}

	// Фильтрация по премиальному статусу
	if isPremium, ok := filter["is_premium"]; ok {
		if isPremiumBool, ok := isPremium.(bool); ok {
			queryBuilder.WriteString(fmt.Sprintf(" AND is_premium = $%d", argCounter))
			args = append(args, isPremiumBool)
			argCounter++
		}
	}

	// Сортировка результатов
	if sortBy != "" {
		order := "ASC"
		sortField := sortBy
		if strings.HasPrefix(sortBy, "-") {
			order = "DESC"
			sortField = strings.TrimPrefix(sortBy, "-")
		}
		// Предотвращение SQL-инъекций путем проверки допустимых полей сортировки
		allowedSortFields := map[string]bool{
			"name":           true,
			"origin_country": true,
			"popularity":     true,
			"founded_year":   true,
			"created_at":     true,
			"updated_at":     true,
		}
		if allowedSortFields[sortField] {
			queryBuilder.WriteString(fmt.Sprintf(" ORDER BY %s %s", sortField, order))
			span.SetTag("sort.field", sortField)
			span.SetTag("sort.order", order)
		} else {
			span.LogFields(
				log.String("event", "invalid_sort_field"),
				log.String("sort_by", sortBy),
			)
			r.log.Warn().
				Str("sort_by", sortBy).
				Msg("Invalid sort field provided, ignoring sorting")
		}
	}

	finalQuery := queryBuilder.String()
	rows, err := r.pool.Query(ctx, finalQuery, args...)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(
			log.String("err", err.Error()),
			log.String("query", finalQuery),
		)
		r.log.Error().
			Err(err).
			Msg("Failed to execute GetAll query")
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var brands []dto.Brand

	// Обработка строк результата
	for rows.Next() {
		var brand dto.Brand
		if err = rows.Scan(
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
			&brand.UpdatedAt,
		); err != nil {
			span.SetTag("error", true)
			span.LogFields(log.Error(err))
			r.log.Error().Err(err).Msg("Failed to scan row in GetAll")
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		brands = append(brands, brand)
	}

	// Проверка ошибок после итерации по строкам
	if rows.Err() != nil {
		span.SetTag("error", true)
		span.LogFields(
			log.String("event", "rows_iteration_error"),
			log.String("error", rows.Err().Error()),
		)
		r.log.Error().
			Err(rows.Err()).
			Str("operation", "GetAll").
			Msg("Error iterating over rows in GetAll")
		return nil, fmt.Errorf("error iterating over rows: %w", rows.Err())
	}

	return brands, nil
}
