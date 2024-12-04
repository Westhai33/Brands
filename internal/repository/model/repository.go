package model

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type ModelRepository struct {
	ctx  context.Context
	pool *pgxpool.Pool
	log  zerolog.Logger
}

func New(
	ctx context.Context,
	pool *pgxpool.Pool,
	logger zerolog.Logger,
) (*ModelRepository, error) {
	return &ModelRepository{ctx: ctx, pool: pool, log: logger}, nil
}
