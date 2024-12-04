package model

import (
	"Brands/internal/dto"
	"Brands/internal/repository/model"
	"Brands/pkg/pool"
	"context"
	"github.com/rs/zerolog"
)

// ModelService интерфейс для работы с моделями
type ModelService interface {
	Create(ctx context.Context, model *dto.Model) (int64, error)
	GetByID(ctx context.Context, id int64) (*dto.Model, error)
	Update(ctx context.Context, model *dto.Model) error
	SoftDelete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	GetAll(ctx context.Context, filter map[string]interface{}, sort string) ([]dto.Model, error)
}

// Service представляет слой сервиса для работы с моделями
type Service struct {
	repo       model.ModelRepository
	workerPool *pool.WorkerPool
	log        zerolog.Logger
}

// New создает новый экземпляр Service

func New(
	repo model.ModelRepository,
	workerPool *pool.WorkerPool,
	logger zerolog.Logger,
) ModelService {
	return &Service{
		repo:       repo,
		workerPool: workerPool,
		log:        logger,
	}
}
