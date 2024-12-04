package brand

import (
	"Brands/internal/dto"
	"Brands/internal/repository/brand"
	"Brands/pkg/pool"
	"context"
	"github.com/rs/zerolog"
)

type BrandService interface {
	Create(ctx context.Context, brand *dto.Brand) (int64, error)
	GetByID(ctx context.Context, id int64) (*dto.Brand, error)
	Update(ctx context.Context, brand *dto.Brand) error
	SoftDelete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	GetAll(ctx context.Context, filter map[string]interface{}, sort string) ([]dto.Brand, error)
}

// Service представляет слой сервиса для работы с брендами
type Service struct {
	repo       brand.BrandRepository
	workerPool *pool.WorkerPool
	log        zerolog.Logger
}

// New создает новый экземпляр Service
func New(
	repo brand.BrandRepository,
	workerPool *pool.WorkerPool,
	logger zerolog.Logger,
) BrandService {
	return &Service{
		repo:       repo,
		workerPool: workerPool,
		log:        logger,
	}
}
