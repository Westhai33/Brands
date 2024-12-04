package brand

import (
	"Brands/internal/repository/brand"
	"github.com/rs/zerolog"
)

// BrandService представляет слой сервиса для работы с брендами
type BrandService struct {
	repo *brand.BrandRepository
	log  zerolog.Logger
}

// New создает новый экземпляр BrandService
func New(
	repo *brand.BrandRepository,
	logger zerolog.Logger,
) *BrandService {
	return &BrandService{
		repo: repo,
		log:  logger,
	}
}
