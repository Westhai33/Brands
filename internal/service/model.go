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

	task := make(chan struct {
		id  int64
		err error
	}, 1)

	s.workerPool.Submit(func(workerID int) {
		zerohook.Logger.Info().Int("worker_id", workerID).Msg("Executing Create operation")
		id, err := s.repo.Create(ctx, model)
		task <- struct {
			id  int64
			err error
		}{id: id, err: err}
		close(task)
	})

	result := <-task
	if result.err != nil {
		zerohook.Logger.Error().Err(result.err).Msg("Failed to create model in repository")
		return 0, result.err
	}
	zerohook.Logger.Info().Int64("model_id", result.id).Msg("Model created successfully")
	return result.id, nil
}

// GetByID получает модель по ID
func (s *ModelService) GetByID(ctx context.Context, id int64) (*dto.Model, error) {
	zerohook.Logger.Info().Int64("model_id", id).Msg("Starting GetByID operation")
	task := make(chan struct {
		model *dto.Model
		err   error
	}, 1)

	s.workerPool.Submit(func(workerID int) {
		zerohook.Logger.Info().Int("worker_id", workerID).Msg("Executing GetByID operation")
		model, err := s.repo.GetByID(ctx, id)
		task <- struct {
			model *dto.Model
			err   error
		}{model: model, err: err}
		close(task)
	})

	result := <-task
	if result.err != nil {
		zerohook.Logger.Error().Err(result.err).Int64("model_id", id).Msg("Failed to get model in repository")
		return nil, result.err
	}
	zerohook.Logger.Info().Interface("model", result.model).Msg("Model retrieved successfully")
	return result.model, nil
}

// Update обновляет данные модели
func (s *ModelService) Update(ctx context.Context, model *dto.Model) error {
	zerohook.Logger.Info().Interface("model", model).Msg("Starting Update operation")
	if model.Name == "" {
		err := fmt.Errorf("name is required")
		zerohook.Logger.Error().Err(err).Msg("Validation failed in Update")
		return err
	}

	task := make(chan error, 1)

	s.workerPool.Submit(func(workerID int) {
		zerohook.Logger.Info().Int("worker_id", workerID).Msg("Executing Update operation")
		err := s.repo.Update(ctx, model)
		task <- err
		close(task)
	})

	err := <-task
	if err != nil {
		zerohook.Logger.Error().Err(err).Msg("Failed to update model in repository")
		return err
	}
	zerohook.Logger.Info().Int64("model_id", model.ID).Msg("Model updated successfully")
	return nil
}

// SoftDelete мягко удаляет модель
func (s *ModelService) SoftDelete(ctx context.Context, id int64) error {
	zerohook.Logger.Info().Int64("model_id", id).Msg("Starting SoftDelete operation")
	task := make(chan error, 1)

	s.workerPool.Submit(func(workerID int) {
		zerohook.Logger.Info().Int("worker_id", workerID).Msg("Executing SoftDelete operation")
		err := s.repo.SoftDelete(ctx, id)
		task <- err
		close(task)
	})

	err := <-task
	if err != nil {
		zerohook.Logger.Error().Err(err).Int64("model_id", id).Msg("Failed to soft delete model in repository")
		return err
	}
	zerohook.Logger.Info().Int64("model_id", id).Msg("Model soft deleted successfully")
	return nil
}

// Restore восстанавливает мягко удалённую модель
func (s *ModelService) Restore(ctx context.Context, id int64) error {
	zerohook.Logger.Info().Int64("model_id", id).Msg("Starting Restore operation")
	task := make(chan error, 1)

	s.workerPool.Submit(func(workerID int) {
		zerohook.Logger.Info().Int("worker_id", workerID).Msg("Executing Restore operation")
		err := s.repo.Restore(ctx, id)
		task <- err
		close(task)
	})

	err := <-task
	if err != nil {
		zerohook.Logger.Error().Err(err).Int64("model_id", id).Msg("Failed to restore model in repository")
		return err
	}
	zerohook.Logger.Info().Int64("model_id", id).Msg("Model restored successfully")
	return nil
}

// GetAll получает все модели с фильтрацией и сортировкой
func (s *ModelService) GetAll(ctx context.Context, filter map[string]interface{}, sort string) ([]dto.Model, error) {
	zerohook.Logger.Info().Interface("filter", filter).Str("sort", sort).Msg("Starting GetAll operation")
	task := make(chan struct {
		models []dto.Model
		err    error
	}, 1)

	s.workerPool.Submit(func(workerID int) {
		zerohook.Logger.Info().Int("worker_id", workerID).Msg("Executing GetAll operation")
		models, err := s.repo.GetAll(ctx, filter, sort)
		task <- struct {
			models []dto.Model
			err    error
		}{models: models, err: err}
		close(task)
	})

	result := <-task
	if result.err != nil {
		zerohook.Logger.Error().Err(result.err).Msg("Failed to get all models in repository")
		return nil, result.err
	}
	zerohook.Logger.Info().Int("models_count", len(result.models)).Msg("Models retrieved successfully")
	return result.models, nil
}
