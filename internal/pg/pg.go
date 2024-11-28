package pg

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"sync"
)

type postgres struct {
	db *pgxpool.Pool
}

var (
	pgInstance *postgres
	pgOnce     sync.Once
)

func NewPG(ctx context.Context, conn string, log zerolog.Logger) *postgres {
	pgOnce.Do(func() {
		cfg, err := pgxpool.ParseConfig(conn)
		if err != nil {
			log.Fatal().Err(errors.Wrap(err, "unable to create connection pool")).Send()
			return
		}
		cfg.ConnConfig.Tracer = &myQueryTracer{
			log: log,
		}
		db, err := pgxpool.NewWithConfig(ctx, cfg)
		if err != nil {
			log.Fatal().Err(errors.Wrap(err, "unable to create connection pool")).Send()
			return
		}
		pgInstance = &postgres{db}
	})

	return pgInstance
}

func (pg *postgres) Close() {
	pg.db.Close()
}

func (pg *postgres) Pool() *pgxpool.Pool {
	return pg.db
}

type myQueryTracer struct {
	log zerolog.Logger
}

func (tracer *myQueryTracer) TraceQueryStart(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TraceQueryStartData) context.Context {
	tracer.log.Info().Str("sql", data.SQL).Interface("args", data.Args).Send()

	return ctx
}

func (tracer *myQueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	if data.Err != nil {
		tracer.log.Error().Err(data.Err).Interface("args", data.CommandTag).Send()
	}
}
