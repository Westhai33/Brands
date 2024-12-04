package brand

import (
	"Brands/internal/dto"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type BrandRepository interface {
	GetAll(
		ctx context.Context,
		filter map[string]interface{},
		sortBy string,
	) ([]dto.Brand, error)
	Create(ctx context.Context, brand *dto.Brand) (int64, error)
	GetByID(ctx context.Context, id int64) (*dto.Brand, error)
	Update(ctx context.Context, brand *dto.Brand) error
	SoftDelete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
}

type Repository struct {
	ctx  context.Context
	pool *pgxpool.Pool
	log  zerolog.Logger
}

func New(
	ctx context.Context,
	pool *pgxpool.Pool,
	logger zerolog.Logger,
) (BrandRepository, error) {
	return &Repository{ctx: ctx, pool: pool, log: logger}, nil
}
