package service

import (
	"Brands/internal/dto"
	"Brands/internal/repository"
	"Brands/pkg/pool"
	"Brands/pkg/zerohook"
	"context"
	"fmt"
)

// BrandServiceInterface интерфейс для работы с брендами
type BrandServiceInterface interface {
	Create(ctx context.Context, brand *dto.Brand) (int64, error)
	GetByID(ctx context.Context, id int64) (*dto.Brand, error)
	Update(ctx context.Context, brand *dto.Brand) error
	SoftDelete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	GetAll(ctx context.Context, filter map[string]interface{}, sort string) ([]dto.Brand, error)
}

// BrandService представляет слой сервиса для работы с брендами
type BrandService struct {
	repo       *repository.BrandRepository
	workerPool *pool.WorkerPool
}

// NewBrandService создает новый экземпляр BrandService
func NewBrandService(repo *repository.BrandRepository, workerPool *pool.WorkerPool) *BrandService {
	return &BrandService{
		repo:       repo,
		workerPool: workerPool,
	}
}

// Create создает новый бренд
func (s *BrandService) Create(ctx context.Context, brand *dto.Brand) (int64, error) {
	zerohook.Logger.Info().Msg("Starting Create operation")
	if brand.Name == "" {
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
		id, err := s.repo.Create(ctx, brand)
		task <- struct {
			id  int64
			err error
		}{id: id, err: err}
		close(task)
	})

	result := <-task
	if result.err != nil {
		zerohook.Logger.Error().Err(result.err).Msg("Failed to create brand in repository")
		return 0, result.err
	}
	zerohook.Logger.Info().Int64("brand_id", result.id).Msg("Brand created successfully")
	return result.id, nil
}

// GetByID получает бренд по ID
func (s *BrandService) GetByID(ctx context.Context, id int64) (*dto.Brand, error) {
	zerohook.Logger.Info().Int64("brand_id", id).Msg("Starting GetByID operation")
	task := make(chan struct {
		brand *dto.Brand
		err   error
	}, 1)

	s.workerPool.Submit(func(workerID int) {
		zerohook.Logger.Info().Int("worker_id", workerID).Msg("Executing GetByID operation")
		brand, err := s.repo.GetByID(ctx, id)
		task <- struct {
			brand *dto.Brand
			err   error
		}{brand: brand, err: err}
		close(task)
	})

	result := <-task
	if result.err != nil {
		zerohook.Logger.Error().Err(result.err).Int64("brand_id", id).Msg("Failed to get brand in repository")
		return nil, result.err
	}
	zerohook.Logger.Info().Interface("brand", result.brand).Msg("Brand retrieved successfully")
	return result.brand, nil
}

// Update обновляет данные бренда
func (s *BrandService) Update(ctx context.Context, brand *dto.Brand) error {
	zerohook.Logger.Info().Interface("brand", brand).Msg("Starting Update operation")
	if brand.Name == "" {
		err := fmt.Errorf("name is required")
		zerohook.Logger.Error().Err(err).Msg("Validation failed in Update")
		return err
	}

	task := make(chan error, 1)

	s.workerPool.Submit(func(workerID int) {
		zerohook.Logger.Info().Int("worker_id", workerID).Msg("Executing Update operation")
		err := s.repo.Update(ctx, brand)
		task <- err
		close(task)
	})

	err := <-task
	if err != nil {
		zerohook.Logger.Error().Err(err).Msg("Failed to update brand in repository")
		return err
	}
	zerohook.Logger.Info().Int64("brand_id", brand.ID).Msg("Brand updated successfully")
	return nil
}

// SoftDelete мягко удаляет бренд
func (s *BrandService) SoftDelete(ctx context.Context, id int64) error {
	zerohook.Logger.Info().Int64("brand_id", id).Msg("Starting SoftDelete operation")
	task := make(chan error, 1)

	s.workerPool.Submit(func(workerID int) {
		zerohook.Logger.Info().Int("worker_id", workerID).Msg("Executing SoftDelete operation")
		err := s.repo.SoftDelete(ctx, id)
		task <- err
		close(task)
	})

	err := <-task
	if err != nil {
		zerohook.Logger.Error().Err(err).Int64("brand_id", id).Msg("Failed to soft delete brand in repository")
		return err
	}
	zerohook.Logger.Info().Int64("brand_id", id).Msg("Brand soft deleted successfully")
	return nil
}

// Restore восстанавливает мягко удалённый бренд
func (s *BrandService) Restore(ctx context.Context, id int64) error {
	zerohook.Logger.Info().Int64("brand_id", id).Msg("Starting Restore operation")
	task := make(chan error, 1)

	s.workerPool.Submit(func(workerID int) {
		zerohook.Logger.Info().Int("worker_id", workerID).Msg("Executing Restore operation")
		err := s.repo.Restore(ctx, id)
		task <- err
		close(task)
	})

	err := <-task
	if err != nil {
		zerohook.Logger.Error().Err(err).Int64("brand_id", id).Msg("Failed to restore brand in repository")
		return err
	}
	zerohook.Logger.Info().Int64("brand_id", id).Msg("Brand restored successfully")
	return nil
}

// GetAll получает все бренды с фильтрацией и сортировкой
func (s *BrandService) GetAll(ctx context.Context, filter map[string]interface{}, sort string) ([]dto.Brand, error) {
	zerohook.Logger.Info().Interface("filter", filter).Str("sort", sort).Msg("Starting GetAll operation")
	task := make(chan struct {
		brands []dto.Brand
		err    error
	}, 1)

	s.workerPool.Submit(func(workerID int) {
		zerohook.Logger.Info().Int("worker_id", workerID).Msg("Executing GetAll operation")
		brands, err := s.repo.GetAll(ctx, filter, sort)
		task <- struct {
			brands []dto.Brand
			err    error
		}{brands: brands, err: err}
		close(task)
	})

	result := <-task
	if result.err != nil {
		zerohook.Logger.Error().Err(result.err).Msg("Failed to get all brands in repository")
		return nil, result.err
	}
	zerohook.Logger.Info().Int("brands_count", len(result.brands)).Msg("Brands retrieved successfully")
	return result.brands, nil
}
