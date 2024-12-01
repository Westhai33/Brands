package service

import (
	"Brands/internal/dto"
	"Brands/internal/repository"
	"Brands/pkg/pool"
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
	if brand.Name == "" {
		return 0, fmt.Errorf("name is required")
	}

	resultChan := make(chan int64, 1)
	errorChan := make(chan error, 1)

	s.workerPool.SubmitTask(func() {
		id, err := s.repo.Create(ctx, brand)
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

// GetByID получает бренд по ID
func (s *BrandService) GetByID(ctx context.Context, id int64) (*dto.Brand, error) {
	resultChan := make(chan *dto.Brand, 1)
	errorChan := make(chan error, 1)

	s.workerPool.SubmitTask(func() {
		brand, err := s.repo.GetByID(ctx, id)
		if err != nil {
			errorChan <- err
			return
		}
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
	if brand.Name == "" {
		return fmt.Errorf("name is required")
	}

	errorChan := make(chan error, 1)

	s.workerPool.SubmitTask(func() {
		err := s.repo.Update(ctx, brand)
		errorChan <- err
	})

	return <-errorChan
}

// SoftDelete мягко удаляет бренд
func (s *BrandService) SoftDelete(ctx context.Context, id int64) error {
	errorChan := make(chan error, 1)

	s.workerPool.SubmitTask(func() {
		err := s.repo.SoftDelete(ctx, id)
		errorChan <- err
	})

	return <-errorChan
}

// Restore восстанавливает мягко удалённый бренд
func (s *BrandService) Restore(ctx context.Context, id int64) error {
	errorChan := make(chan error, 1)

	s.workerPool.SubmitTask(func() {
		err := s.repo.Restore(ctx, id)
		errorChan <- err
	})

	return <-errorChan
}

// GetAll получает все бренды с фильтрацией и сортировкой
func (s *BrandService) GetAll(ctx context.Context, filter map[string]interface{}, sort string) ([]dto.Brand, error) {
	resultChan := make(chan []dto.Brand, 1)
	errorChan := make(chan error, 1)

	s.workerPool.SubmitTask(func() {
		brands, err := s.repo.GetAll(ctx, filter, sort)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- brands
	})

	select {
	case err := <-errorChan:
		return nil, err
	case brands := <-resultChan:
		return brands, nil
	}
}
