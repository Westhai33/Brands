package service

import (
	"Brands/internal/dto"
	"Brands/internal/repository"
	"Brands/pkg/pool"
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
	repo       *repository.ModelRepository
	workerPool *pool.WorkerPool
}

// NewModelService создает новый экземпляр ModelService
func NewModelService(repo *repository.ModelRepository, workerPool *pool.WorkerPool) *ModelService {
	return &ModelService{
		repo:       repo,
		workerPool: workerPool,
	}
}

// Create создает новую модель
func (s *ModelService) Create(ctx context.Context, model *dto.Model) (int64, error) {
	if model.Name == "" {
		return 0, fmt.Errorf("name is required")
	}

	resultChan := make(chan int64, 1)
	errorChan := make(chan error, 1)

	s.workerPool.SubmitTask(func() {
		id, err := s.repo.Create(ctx, model)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- id
	})

	select {
	case err := <-errorChan:
		return 0, err
	case id := <-resultChan:
		return id, nil
	}
}

// GetByID получает модель по ID
func (s *ModelService) GetByID(ctx context.Context, id int64) (*dto.Model, error) {
	resultChan := make(chan *dto.Model, 1)
	errorChan := make(chan error, 1)

	s.workerPool.SubmitTask(func() {
		model, err := s.repo.GetByID(ctx, id)
		if err != nil {
			errorChan <- err
			return
		}
		if model == nil {
			errorChan <- fmt.Errorf("model with ID %d not found", id)
			return
		}
		resultChan <- model
	})

	select {
	case err := <-errorChan:
		return nil, err
	case model := <-resultChan:
		return model, nil
	}
}

// Update обновляет данные модели
func (s *ModelService) Update(ctx context.Context, model *dto.Model) error {
	if model.Name == "" {
		return fmt.Errorf("name is required")
	}

	errorChan := make(chan error, 1)

	s.workerPool.SubmitTask(func() {
		err := s.repo.Update(ctx, model)
		errorChan <- err
	})

	return <-errorChan
}

// SoftDelete мягко удаляет модель
func (s *ModelService) SoftDelete(ctx context.Context, id int64) error {
	errorChan := make(chan error, 1)

	s.workerPool.SubmitTask(func() {
		err := s.repo.SoftDelete(ctx, id)
		errorChan <- err
	})

	return <-errorChan
}

// Restore восстанавливает мягко удалённую модель
func (s *ModelService) Restore(ctx context.Context, id int64) error {
	errorChan := make(chan error, 1)

	s.workerPool.SubmitTask(func() {
		err := s.repo.Restore(ctx, id)
		errorChan <- err
	})

	return <-errorChan
}

// GetAll получает все модели с фильтрацией и сортировкой
func (s *ModelService) GetAll(ctx context.Context, filter map[string]interface{}, sortBy string) ([]dto.Model, error) {
	resultChan := make(chan []dto.Model, 1)
	errorChan := make(chan error, 1)

	s.workerPool.SubmitTask(func() {
		models, err := s.repo.GetAll(ctx, filter, sortBy)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- models
	})

	select {
	case err := <-errorChan:
		return nil, err
	case models := <-resultChan:
		return models, nil
	}
}
