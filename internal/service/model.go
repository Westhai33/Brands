package service

import (
	"Brands/internal/dto"
	"Brands/internal/repository"
	"Brands/pkg/pool"
	"Brands/pkg/zerohook"
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
	GetAll(ctx context.Context, filter map[string]interface{}, sort string) ([]dto.Model, error)
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
	zerohook.Logger.Info().Msg("Starting Create operation")
	if model.Name == "" {
		err := fmt.Errorf("name is required")
		zerohook.Logger.Error().Err(err).Msg("Validation failed in Create")
		return 0, err
	}

	resultChan := make(chan int64, 1)
	errorChan := make(chan error, 1)

	s.workerPool.SubmitTask(func() {
		id, err := s.repo.Create(ctx, model)
		if err != nil {
			zerohook.Logger.Error().Err(err).Msg("Failed to create model in repository")
			errorChan <- err
			return
		}
		zerohook.Logger.Info().Int64("model_id", id).Msg("Model created successfully")
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
	zerohook.Logger.Info().Int64("model_id", id).Msg("Starting GetByID operation")
	resultChan := make(chan *dto.Model, 1)
	errorChan := make(chan error, 1)

	s.workerPool.SubmitTask(func() {
		model, err := s.repo.GetByID(ctx, id)
		if err != nil {
			zerohook.Logger.Error().Err(err).Int64("model_id", id).Msg("Failed to get model in repository")
			errorChan <- err
			return
		}
		zerohook.Logger.Info().Interface("model", model).Msg("Model retrieved successfully")
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
	zerohook.Logger.Info().Interface("model", model).Msg("Starting Update operation")
	if model.Name == "" {
		err := fmt.Errorf("name is required")
		zerohook.Logger.Error().Err(err).Msg("Validation failed in Update")
		return err
	}

	errorChan := make(chan error, 1)

	s.workerPool.SubmitTask(func() {
		err := s.repo.Update(ctx, model)
		if err != nil {
			zerohook.Logger.Error().Err(err).Msg("Failed to update model in repository")
		} else {
			zerohook.Logger.Info().Int64("model_id", model.ID).Msg("Model updated successfully")
		}
		errorChan <- err
	})

	return <-errorChan
}

// SoftDelete мягко удаляет модель
func (s *ModelService) SoftDelete(ctx context.Context, id int64) error {
	zerohook.Logger.Info().Int64("model_id", id).Msg("Starting SoftDelete operation")
	errorChan := make(chan error, 1)

	s.workerPool.SubmitTask(func() {
		err := s.repo.SoftDelete(ctx, id)
		if err != nil {
			zerohook.Logger.Error().Err(err).Int64("model_id", id).Msg("Failed to soft delete model in repository")
		} else {
			zerohook.Logger.Info().Int64("model_id", id).Msg("Model soft deleted successfully")
		}
		errorChan <- err
	})

	return <-errorChan
}

// Restore восстанавливает мягко удалённую модель
func (s *ModelService) Restore(ctx context.Context, id int64) error {
	zerohook.Logger.Info().Int64("model_id", id).Msg("Starting Restore operation")
	errorChan := make(chan error, 1)

	s.workerPool.SubmitTask(func() {
		err := s.repo.Restore(ctx, id)
		if err != nil {
			zerohook.Logger.Error().Err(err).Int64("model_id", id).Msg("Failed to restore model in repository")
		} else {
			zerohook.Logger.Info().Int64("model_id", id).Msg("Model restored successfully")
		}
		errorChan <- err
	})

	return <-errorChan
}

// GetAll получает все модели с фильтрацией и сортировкой
func (s *ModelService) GetAll(ctx context.Context, filter map[string]interface{}, sort string) ([]dto.Model, error) {
	zerohook.Logger.Info().Interface("filter", filter).Str("sort", sort).Msg("Starting GetAll operation")
	resultChan := make(chan []dto.Model, 1)
	errorChan := make(chan error, 1)

	s.workerPool.SubmitTask(func() {
		models, err := s.repo.GetAll(ctx, filter, sort)
		if err != nil {
			zerohook.Logger.Error().Err(err).Msg("Failed to get all models in repository")
			errorChan <- err
			return
		}
		zerohook.Logger.Info().Int("models_count", len(models)).Msg("Models retrieved successfully")
		resultChan <- models
	})

	select {
	case err := <-errorChan:
		return nil, err
	case models := <-resultChan:
		return models, nil
	}
}
