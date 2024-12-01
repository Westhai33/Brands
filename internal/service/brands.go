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

	resultChan := make(chan int64, 1)
	errorChan := make(chan error, 1)

	s.workerPool.SubmitTask(func() {
		id, err := s.repo.Create(ctx, brand)
		if err != nil {
			zerohook.Logger.Error().Err(err).Msg("Failed to create brand in repository")
			errorChan <- err
			return
		}
		zerohook.Logger.Info().Int64("brand_id", id).Msg("Brand created successfully")
		resultChan <- id
	})

	select {
	case err := <-errorChan:
		return 0, err
	case id := <-resultChan:
		return id, nil
	}
}

// GetByID получает бренд по ID
func (s *BrandService) GetByID(ctx context.Context, id int64) (*dto.Brand, error) {
	zerohook.Logger.Info().Int64("brand_id", id).Msg("Starting GetByID operation")
	resultChan := make(chan *dto.Brand, 1)
	errorChan := make(chan error, 1)

	s.workerPool.SubmitTask(func() {
		brand, err := s.repo.GetByID(ctx, id)
		if err != nil {
			zerohook.Logger.Error().Err(err).Int64("brand_id", id).Msg("Failed to get brand in repository")
			errorChan <- err
			return
		}
		zerohook.Logger.Info().Interface("brand", brand).Msg("Brand retrieved successfully")
		resultChan <- brand
	})

	select {
	case err := <-errorChan:
		return nil, err
	case brand := <-resultChan:
		return brand, nil
	}
}

// Update обновляет данные бренда
func (s *BrandService) Update(ctx context.Context, brand *dto.Brand) error {
	zerohook.Logger.Info().Interface("brand", brand).Msg("Starting Update operation")
	if brand.Name == "" {
		err := fmt.Errorf("name is required")
		zerohook.Logger.Error().Err(err).Msg("Validation failed in Update")
		return err
	}

	errorChan := make(chan error, 1)

	s.workerPool.SubmitTask(func() {
		err := s.repo.Update(ctx, brand)
		if err != nil {
			zerohook.Logger.Error().Err(err).Msg("Failed to update brand in repository")
			errorChan <- err
			return
		}
		zerohook.Logger.Info().Int64("brand_id", brand.ID).Msg("Brand updated successfully")
		errorChan <- nil
	})

	return <-errorChan
}

// SoftDelete мягко удаляет бренд
func (s *BrandService) SoftDelete(ctx context.Context, id int64) error {
	zerohook.Logger.Info().Int64("brand_id", id).Msg("Starting SoftDelete operation")
	errorChan := make(chan error, 1)

	s.workerPool.SubmitTask(func() {
		err := s.repo.SoftDelete(ctx, id)
		if err != nil {
			zerohook.Logger.Error().Err(err).Int64("brand_id", id).Msg("Failed to soft delete brand in repository")
			errorChan <- err
			return
		}
		zerohook.Logger.Info().Int64("brand_id", id).Msg("Brand soft deleted successfully")
		errorChan <- nil
	})

	return <-errorChan
}

// Restore восстанавливает мягко удалённый бренд
func (s *BrandService) Restore(ctx context.Context, id int64) error {
	zerohook.Logger.Info().Int64("brand_id", id).Msg("Starting Restore operation")
	errorChan := make(chan error, 1)

	s.workerPool.SubmitTask(func() {
		err := s.repo.Restore(ctx, id)
		if err != nil {
			zerohook.Logger.Error().Err(err).Int64("brand_id", id).Msg("Failed to restore brand in repository")
			errorChan <- err
			return
		}
		zerohook.Logger.Info().Int64("brand_id", id).Msg("Brand restored successfully")
		errorChan <- nil
	})

	return <-errorChan
}

// GetAll получает все бренды с фильтрацией и сортировкой
func (s *BrandService) GetAll(ctx context.Context, filter map[string]interface{}, sort string) ([]dto.Brand, error) {
	zerohook.Logger.Info().Interface("filter", filter).Str("sort", sort).Msg("Starting GetAll operation")
	resultChan := make(chan []dto.Brand, 1)
	errorChan := make(chan error, 1)

	s.workerPool.SubmitTask(func() {
		brands, err := s.repo.GetAll(ctx, filter, sort)
		if err != nil {
			zerohook.Logger.Error().Err(err).Msg("Failed to get all brands in repository")
			errorChan <- err
			return
		}
		zerohook.Logger.Info().Int("brands_count", len(brands)).Msg("Brands retrieved successfully")
		resultChan <- brands
	})

	select {
	case err := <-errorChan:
		return nil, err
	case brands := <-resultChan:
		return brands, nil
	}
}
