package main

import (
	_ "Brands/docs"
	"Brands/internal/api"
	brandhandler "Brands/internal/api/handler/brand"
	modelhandler "Brands/internal/api/handler/model"
	"Brands/internal/config"
	"Brands/internal/metrics"
	"Brands/internal/pg"
	"Brands/internal/repository/brand"
	"Brands/internal/repository/model"
	brandservice "Brands/internal/service/brand"
	modelservice "Brands/internal/service/model"
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
	bs := brandservice.New(br, zerohook.Logger)
	ms := modelservice.New(mr, zerohook.Logger)

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
