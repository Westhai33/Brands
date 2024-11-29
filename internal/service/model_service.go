package service

import (
	"Brands/internal/dto"
	"Brands/internal/repository"
	"context"
	"fmt"
)

// ModelServiceInterface интерфейс для работы с моделями
type ModelServiceInterface interface {
	Create(ctx context.Context, model *dto.Model) (int64, error)
	GetByID(ctx context.Context, id int64) (*dto.Model, error)
	Update(ctx context.Context, model *dto.Model) error
	SoftDelete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	GetAll(ctx context.Context, filter, sort, order string) ([]dto.Model, error)
}

// ModelService представляет слой сервиса для работы с моделями
type ModelService struct {
	repo *repository.ModelRepository
}

// NewModelService создает новый экземпляр ModelService
func NewModelService(repo *repository.ModelRepository) *ModelService {
	return &ModelService{repo: repo}
}

// Create создает новую модель
func (s *ModelService) Create(ctx context.Context, model *dto.Model) (int64, error) {
	// Проверка входных данных (можно добавить дополнительные проверки)
	if model.Name == "" {
		return 0, fmt.Errorf("name is required")
	}

	// Вызов репозитория для создания модели
	return s.repo.Create(ctx, model)
}

// GetByID получает модель по ID
func (s *ModelService) GetByID(ctx context.Context, id int64) (*dto.Model, error) {
	// Вызов репозитория для получения модели
	model, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if model == nil {
		return nil, fmt.Errorf("model with ID %d not found", id)
	}
	return model, nil
}

// Update обновляет данные модели
func (s *ModelService) Update(ctx context.Context, model *dto.Model) error {
	// Проверка входных данных (можно добавить дополнительные проверки)
	if model.Name == "" {
		return fmt.Errorf("name is required")
	}

	// Вызов репозитория для обновления модели
	return s.repo.Update(ctx, model)
}

// SoftDelete мягко удаляет модель
func (s *ModelService) SoftDelete(ctx context.Context, id int64) error {
	// Вызов репозитория для мягкого удаления модели
	return s.repo.SoftDelete(ctx, id)
}

// Restore восстанавливает мягко удалённую модель
func (s *ModelService) Restore(ctx context.Context, id int64) error {
	// Вызов репозитория для восстановления модели
	return s.repo.Restore(ctx, id)
}

// GetAll получает все модели с фильтрацией и сортировкой
func (s *ModelService) GetAll(ctx context.Context, filter, sort, order string) ([]dto.Model, error) {
	// Вызов репозитория для получения всех моделей
	models, err := s.repo.GetAll(ctx, filter, sort, order)
	if err != nil {
		return nil, err
	}
	return models, nil
}
