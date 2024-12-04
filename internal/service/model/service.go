package model

import (
	"Brands/internal/repository/model"
	"github.com/rs/zerolog"
)

// ModelService представляет слой сервиса для работы с моделями
type ModelService struct {
	repo *model.ModelRepository
	log  zerolog.Logger
}

// New создает новый экземпляр ModelService

func New(
	repo *model.ModelRepository,
	logger zerolog.Logger,
) *ModelService {
	return &ModelService{
		repo: repo,
		log:  logger,
	}
}
