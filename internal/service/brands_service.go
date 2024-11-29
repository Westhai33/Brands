package service

import (
	"Brands/internal/dto"
	"Brands/internal/repository"
	"context"
	"fmt"
)

type BrandServiceInterface interface {
	Create(ctx context.Context, brand *dto.Brand) (int64, error)
	GetByID(ctx context.Context, id int64) (*dto.Brand, error)
	Update(ctx context.Context, brand *dto.Brand) error
	SoftDelete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	GetAll(ctx context.Context, filter, sort string) ([]dto.Brand, error)
}

// BrandService представляет слой сервиса для работы с брендами
type BrandService struct {
	repo *repository.BrandRepository
}

// NewBrandService создает новый экземпляр BrandService
func NewBrandService(repo *repository.BrandRepository) *BrandService {
	return &BrandService{repo: repo}
}

// Create создает новый бренд
func (s *BrandService) Create(ctx context.Context, brand *dto.Brand) (int64, error) {
	// Проверка входных данных (можно добавить дополнительные проверки)
	if brand.Name == "" {
		return 0, fmt.Errorf("name is required")
	}

	// Вызов репозитория для создания бренда
	return s.repo.Create(ctx, brand)
}

// GetByID получает бренд по ID
func (s *BrandService) GetByID(ctx context.Context, id int64) (*dto.Brand, error) {
	// Вызов репозитория для получения бренда
	brand, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return brand, nil
}

// Update обновляет данные бренда
func (s *BrandService) Update(ctx context.Context, brand *dto.Brand) error {
	// Проверка входных данных (можно добавить дополнительные проверки)
	if brand.Name == "" {
		return fmt.Errorf("name is required")
	}

	// Вызов репозитория для обновления бренда
	return s.repo.Update(ctx, brand)
}

// SoftDelete мягко удаляет бренд
func (s *BrandService) SoftDelete(ctx context.Context, id int64) error {
	// Вызов репозитория для мягкого удаления бренда
	return s.repo.SoftDelete(ctx, id)
}

// Restore восстанавливает мягко удалённый бренд
func (s *BrandService) Restore(ctx context.Context, id int64) error {
	// Вызов репозитория для восстановления бренда
	return s.repo.Restore(ctx, id)
}

// GetAll получает все бренды с фильтрацией и сортировкой
func (s *BrandService) GetAll(ctx context.Context, filter, sort string) ([]dto.Brand, error) {
	// Вызов репозитория для получения всех брендов
	brands, err := s.repo.GetAll(ctx, filter, sort)
	if err != nil {
		return nil, err
	}
	return brands, nil
}
