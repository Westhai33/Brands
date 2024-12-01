package main

import (
	"Brands/internal/api"
	"Brands/internal/api/handler"
	"Brands/internal/config"
	"Brands/internal/dto"
	"Brands/internal/pg"
	"Brands/internal/repository"
	"Brands/internal/service"
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
// @description Микро-сервис для работы с брендами
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
	//TODO: тут трассировку еще надо

	// Подключаемся к базе данных
	pgInstance, err := pg.NewPG(ctx, "postgres://brands:pgpwdbrands@postgres:5432/brands?sslmode=disable", zerohook.Logger)
	if err != nil {
		zerohook.Logger.Fatal().Msg("Ошибка подключения к базе данных")
		return
	}
	defer pgInstance.Close()

	// Создание репозитория и сервиса
	br, err := repository.NewBrandRepository(ctx, pgInstance.Pool(), zerohook.Logger)
	if err != nil {
		zerohook.Logger.Fatal().Err(err)
		return
	}
	bs := service.NewBrandService(br)

	// TODO: add mh
	//mr, err := repository.NewModelRepository(ctx, pgInstance.Pool(), zerohook.Logger)
	//if err != nil {
	//	zerohook.Logger.Fatal().Err(err)
	//	return
	//}
	//ms := service.NewBrandService(br)

	bh := handler.NewBrandHandler(bs)
	apiService, err := api.NewService(zerohook.Logger, bh)
	if err != nil {
		zerohook.Logger.Fatal().Err(err)
		return
	}

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
	zerohook.Logger.Info().Interface("obj", brand).Send()
	ass, err := bs.Create(ctx, brand)
	zerohook.Logger.Info().Err(err).Interface("obj", ass).Send()
	if err != nil {
		zerohook.Logger.Info().Err(err)

		zerohook.Logger.Fatal().Err(err)
	}
	err = apiService.Start(ctx)
	if err != nil {
		zerohook.Logger.Fatal().Err(err)
		return
	}

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
