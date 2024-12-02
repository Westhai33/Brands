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
func NewPG(ctx context.Context, conn string, log zerolog.Logger) (*postgres, error) {
	var err error
	pgOnce.Do(func() {
		// Парсим конфигурацию подключения
		cfg, parseErr := pgxpool.ParseConfig(conn)
		if parseErr != nil {
			err = errors.Wrap(parseErr, "unable to parse connection config")
			log.Error().Err(err).Send()
			return
		}

		// Добавляем трейсер для логирования запросов
		cfg.ConnConfig.Tracer = &myQueryTracer{
			log: log,
		}

		// Создаем пул соединений с базой данных
		db, poolErr := pgxpool.NewWithConfig(ctx, cfg)
		if poolErr != nil {
			err = errors.Wrap(poolErr, "unable to create connection pool")
			log.Error().Err(err).Send()
			return
		}

		// Проверка соединения с базой данных
		err = checkDBConnection(ctx, db, log)
		if err != nil {
			err = errors.Wrap(err, "database connection check failed")
			log.Error().Err(err).Send()
			return
		}

		// Инициализация глобальной переменной
		pgInstance = &postgres{db}
	})

	if err != nil {
		// Если ошибка произошла в процессе инициализации, возвращаем её
		return nil, err
	}

	// Если инициализация прошла успешно, возвращаем пул соединений
	return pgInstance, nil
}

// Проверка соединения с базой данных
func checkDBConnection(ctx context.Context, db *pgxpool.Pool, log zerolog.Logger) error {
	// Пробный запрос для проверки подключения
	_, err := db.Acquire(ctx)
	if err != nil {
		log.Error().Err(err).Msg("unable to acquire a connection from the pool")
		return err
	}
	return nil
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
