package main

import (
	_ "Brands/docs"
	"Brands/internal/api"
	brandhandler "Brands/internal/api/handler/brand"
	modelhandler "Brands/internal/api/handler/model"
	"Brands/internal/config"
	"Brands/internal/dto"
	"Brands/internal/metrics"
	"Brands/internal/pg"
	"Brands/internal/repository/brand"
	"Brands/internal/repository/model"
	brandservice "Brands/internal/service/brand"
	modelservice "Brands/internal/service/model"
	"Brands/pkg/pool"
	"Brands/pkg/tracer"
	"Brands/pkg/yamlreader"
	"Brands/pkg/zerohook"
	"context"
	"flag"
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Brands API
// @version 1.0
// @description Микро-сервис для работы с брендами и моделями
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host 127.0.0.1:8080
// @BasePath /

func main() {
	ctx := context.Background()
	cfg := MustNewConfig(parseFlags(), zerohook.Logger)
	zerohook.InitLogger(cfg.Log)

	err := tracer.InitTracer(
		cfg.Jaeger.ServiceName,
		cfg.Jaeger.AgentHost,
		cfg.Jaeger.AgentPort,
	)
	if err != nil {
		zerohook.Logger.Fatal().Msgf("Ошибка инициализации трейсера: %v", err)
		return
	}

	// Подключение к базе данных
	pgInstance, err := pg.NewPG(ctx, cfg.Postgres.Conn, zerohook.Logger)
	if err != nil {
		zerohook.Logger.Fatal().Msg("Ошибка подключения к базе данных")
		return
	}
	defer pgInstance.Close()

	// Создание WorkerPool
	wp := pool.NewWorkerPool(ctx)
	defer wp.Stop()

	// Создание репозиториев
	br, err := brand.New(ctx, pgInstance.Pool(), zerohook.Logger)
	if err != nil {
		zerohook.Logger.Fatal().Err(err)
		return
	}
	mr, err := model.New(ctx, pgInstance.Pool(), zerohook.Logger)
	if err != nil {
		zerohook.Logger.Fatal().Err(err)
		return
	}

	// Создание сервисов с передачей WorkerPool
	bs := brandservice.New(br, wp, zerohook.Logger)
	ms := modelservice.New(mr, wp, zerohook.Logger)

	// Создание хендлеров
	bh := brandhandler.New(bs)
	mh := modelhandler.New(ms)

	// Создание API-сервиса
	apiService, err := api.NewService(zerohook.Logger, bh, mh)
	if err != nil {
		zerohook.Logger.Fatal().Err(err)
		return
	}
	// Запуск API-сервиса
	go func() {
		if err := apiService.Start(ctx); err != nil {
			zerohook.Logger.Fatal().Err(err)
		}
	}()

	go metrics.StartPrometheusServer(fmt.Sprintf(":%d", cfg.Prometheus.Port))

	// Пример создания бренда с использованием WorkerPool
	wp.Submit(func(workerID int) {
		brand := &dto.Brand{
			Name:          "SuperBrand",
			Link:          "https://superbrand.com",
			Description:   "SuperBrand is known for its high-quality products and innovative designs.",
			LogoURL:       "https://superbrand.com/logo.png",
			CoverImageURL: "https://superbrand.com/cover.jpg",
			FoundedYear:   1998,
			OriginCountry: "USA",
			Popularity:    85,
			IsPremium:     true,
			IsUpcoming:    false,
			IsDeleted:     false,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		zerohook.Logger.Info().Interface("obj", brand).Int("worker_id", workerID).Send()
		_, err := bs.Create(ctx, brand)
		if err != nil {
			zerohook.Logger.Error().Err(err).Msg("Failed to create brand asynchronously")
		}
	})

	// Пример создания модели с использованием WorkerPool
	wp.Submit(func(workerID int) {
		model := &dto.Model{
			BrandID:     1, // Замените на реальный ID бренда
			Name:        "Model X",
			ReleaseDate: time.Now(),
			IsUpcoming:  false,
			IsLimited:   true,
		}
		zerohook.Logger.Info().Interface("obj", model).Int("worker_id", workerID).Send()
		_, err := ms.Create(ctx, model)
		if err != nil {
			zerohook.Logger.Error().Err(err).Msg("Failed to create model asynchronously")
		}
	})

	// Завершение программы
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func MustNewConfig(path string, lgr zerolog.Logger) *config.Config {
	cfg, err := yamlreader.NewConfig[config.Config](path)
	if err != nil {
		lgr.Fatal().Str("path", path).Err(err).Msg("ошибка чтения конфигурации приложения")
		return nil
	}
	return cfg
}

func parseFlags() string {
	var configPath string
	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.Parse()
	return configPath
}
