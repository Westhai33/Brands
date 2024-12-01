package main

import (
	"Brands/internal/api"
	"Brands/internal/api/handler"
	"Brands/internal/config"
	"Brands/internal/dto"
	"Brands/internal/pg"
	"Brands/internal/repository"
	"Brands/internal/service"
	"Brands/pkg/pool"
	"Brands/pkg/yamlreader"
	"os/signal"
	"syscall"
	"time"

	"Brands/pkg/zerohook"
	"context"
	"flag"
	"os"
	"regexp"

	_ "Brands/docs"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/rs/zerolog"
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

	_ = prometheus.Unregister(collectors.NewGoCollector())
	goCollector := collectors.NewGoCollector(
		collectors.WithGoCollectorRuntimeMetrics(
			collectors.GoRuntimeMetricsRule{Matcher: regexp.MustCompile("/cpu/classes/total:cpu-seconds")},
			collectors.GoRuntimeMetricsRule{Matcher: regexp.MustCompile("/cpu/classes/user:cpu-seconds")},
			collectors.GoRuntimeMetricsRule{Matcher: regexp.MustCompile("/memory/classes/heap/objects:bytes")},
			collectors.GoRuntimeMetricsRule{Matcher: regexp.MustCompile("/memory/classes/total:bytes")},
		),
	)
	prometheus.MustRegister(goCollector)

	// Подключение к базе данных
	pgInstance, err := pg.NewPG(ctx, "postgres://brands:pgpwdbrands@postgres:5432/brands?sslmode=disable", zerohook.Logger)
	if err != nil {
		zerohook.Logger.Fatal().Msg("Ошибка подключения к базе данных")
		return
	}
	defer pgInstance.Close()

	// Создание WorkerPool
	wp := pool.NewWorkerPool(ctx) // Инициализация нового WorkerPool
	defer wp.Stop()

	// Создание репозиториев
	br, err := repository.NewBrandRepository(ctx, pgInstance.Pool(), zerohook.Logger)
	if err != nil {
		zerohook.Logger.Fatal().Err(err)
		return
	}
	mr, err := repository.NewModelRepository(ctx, pgInstance.Pool(), zerohook.Logger)
	if err != nil {
		zerohook.Logger.Fatal().Err(err)
		return
	}

	// Создание сервисов с передачей WorkerPool
	bs := service.NewBrandService(br, wp)
	ms := service.NewModelService(mr, wp)

	// Создание хендлеров
	bh := handler.NewBrandHandler(bs)
	mh := handler.NewModelHandler(ms)

	// Создание API-сервиса
	apiService, err := api.NewService(zerohook.Logger, bh, mh)
	if err != nil {
		zerohook.Logger.Fatal().Err(err)
		return
	}

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

	// Запуск API-сервиса
	go func() {
		if err := apiService.Start(ctx); err != nil {
			zerohook.Logger.Fatal().Err(err)
		}
	}()

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
