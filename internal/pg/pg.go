package pg

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"sync"
)

// Объект для работы с PostgreSQL
type postgres struct {
	db *pgxpool.Pool
}

var (
	pgInstance *postgres
	pgOnce     sync.Once
)

// NewPG создает и инициализирует подключение к базе данных PostgreSQL через пул соединений.
// Теперь принимает только два параметра: контекст и строку подключения
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

		// Создание пула соединений с базой данных
		db, err := pgxpool.NewWithConfig(ctx, cfg)
		if err != nil {
			log.Fatal().Err(errors.Wrap(err, "unable to create connection pool")).Send()
			return
		}
		pgInstance = &postgres{db}
	})

	return pgInstance
}

// Close закрывает пул соединений с базой данных
func (pg *postgres) Close() {
	pg.db.Close()
}

// Pool возвращает пул соединений для использования в других местах программы
func (pg *postgres) Pool() *pgxpool.Pool {
	return pg.db
}

// Мок для логирования запросов
type myQueryTracer struct {
	log zerolog.Logger
}

// TraceQueryStart логирует начало SQL-запроса
func (tracer *myQueryTracer) TraceQueryStart(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TraceQueryStartData) context.Context {
	tracer.log.Info().Str("sql", data.SQL).Interface("args", data.Args).Send()
	return ctx
}

// TraceQueryEnd логирует завершение SQL-запроса
func (tracer *myQueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	if data.Err != nil {
		tracer.log.Error().Err(data.Err).Interface("args", data.CommandTag).Send()
	}
}
