package brand

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type BrandRepository struct {
	ctx  context.Context
	pool *pgxpool.Pool
	log  zerolog.Logger
}

func New(
	ctx context.Context,
	pool *pgxpool.Pool,
	logger zerolog.Logger,
) (*BrandRepository, error) {
	return &BrandRepository{ctx: ctx, pool: pool, log: logger}, nil
}
