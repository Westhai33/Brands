package model

import (
	"Brands/internal/dto"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type ModelRepository interface {
	Create(ctx context.Context, model *dto.Model) (int64, error)
	GetByID(ctx context.Context, id int64) (*dto.Model, error)
	Update(ctx context.Context, model *dto.Model) error
	SoftDelete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
	GetAll(ctx context.Context, filter map[string]interface{}, sortBy string) ([]dto.Model, error)
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
) (ModelRepository, error) {
	return &Repository{ctx: ctx, pool: pool, log: logger}, nil
}
