package controller

import (
	"Brands/internal/dto"
	"Brands/internal/service"
	"context"
	"fmt"
)

// CreateBrand создает новый бренд.
func CreateBrand(ctx context.Context, brandService *service.BrandService, name, link, description, logoURL, coverImageURL string, foundedYear int, originCountry string, popularity int, isPremium, isUpcoming bool) (int64, error) {

	if name == "" {
		return 0, fmt.Errorf("название обязательно")
	}

	brand := &dto.Brand{
		Name:          name,
		Link:          link,
		Description:   description,
		LogoURL:       logoURL,
		CoverImageURL: coverImageURL,
		FoundedYear:   foundedYear,
		OriginCountry: originCountry,
		Popularity:    popularity,
		IsPremium:     isPremium,
		IsUpcoming:    isUpcoming,
	}
	brandID, err := brandService.Create(ctx, brand)
	if err != nil {
		return 0, fmt.Errorf("ошибка при создании бренда: %w", err)
	}

	return brandID, nil
}

// GetBrandByID получает бренд по ID.
func GetBrandByID(ctx context.Context, brandService *service.BrandService, brandID int64) (*dto.Brand, error) {

	brand, err := brandService.GetByID(ctx, brandID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения бренда с ID %d: %w", brandID, err)
	}

	return brand, nil
}

// GetAllBrands получает все бренды с фильтрацией и сортировкой.
func GetAllBrands(ctx context.Context, brandService *service.BrandService, filter, sort string) ([]dto.Brand, error) {

	brands, err := brandService.GetAll(ctx, filter, sort)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения всех брендов: %w", err)
	}

	return brands, nil
}

// UpdateBrand обновляет данные бренда с проверкой существования.
func UpdateBrand(ctx context.Context, brandService *service.BrandService, brandID int64, name, link, description, logoURL, coverImageURL string, foundedYear int, originCountry string, popularity int, isPremium, isUpcoming bool) error {

	brand, err := brandService.GetByID(ctx, brandID)
	if err != nil {
		return fmt.Errorf("бренд с ID %d не найден: %w", brandID, err)
	}

	brand.Name = name
	brand.Link = link
	brand.Description = description
	brand.LogoURL = logoURL
	brand.CoverImageURL = coverImageURL
	brand.FoundedYear = foundedYear
	brand.OriginCountry = originCountry
	brand.Popularity = popularity
	brand.IsPremium = isPremium
	brand.IsUpcoming = isUpcoming

	err = brandService.Update(ctx, brand)
	if err != nil {
		return fmt.Errorf("ошибка обновления бренда с ID %d: %w", brandID, err)
	}

	return nil
}

// SoftDeleteBrand мягко удаляет бренд с проверкой существования.
func SoftDeleteBrand(ctx context.Context, brandService *service.BrandService, brandID int64) error {

	_, err := brandService.GetByID(ctx, brandID)
	if err != nil {
		return fmt.Errorf("бренд с ID %d не найден: %w", brandID, err)
	}

	err = brandService.SoftDelete(ctx, brandID)
	if err != nil {
		return fmt.Errorf("ошибка мягкого удаления бренда с ID %d: %w", brandID, err)
	}

	return nil
}

// RestoreBrand восстанавливает мягко удалённый бренд.
func RestoreBrand(ctx context.Context, brandService *service.BrandService, brandID int64) error {

	_, err := brandService.GetByID(ctx, brandID)
	if err != nil {
		return fmt.Errorf("бренд с ID %d не найден: %w", brandID, err)
	}

	err = brandService.Restore(ctx, brandID)
	if err != nil {
		return fmt.Errorf("ошибка восстановления бренда с ID %d: %w", brandID, err)
	}

	return nil
}
