package main

import (
	"Brands/internal/controller"
	"Brands/internal/pg"         // Пакет для работы с PostgreSQL
	"Brands/internal/repository" // Пакет с репозиториями
	"Brands/internal/service"    // Пакет с сервисами
	"context"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
	"os"
)

func main() {
	// Инициализация логирования
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Задаем параметры для подключения к базе данных
	dbName := "postgres"     // Имя базы данных
	dbUser := "postgres"     // Имя пользователя
	dbPassword := "postgres" // Пароль пользователя // Имя контейнера с базой данных в Docker Compose (или хост)
	dbPort := 5432           // Порт базы данных

	// Формируем строку подключения
	connString := fmt.Sprintf("postgres://%s:%s@localhost:%d/%s?sslmode=disable", dbUser, dbPassword, dbPort, dbName)

	// Подключаемся к базе данных
	ctx := context.Background()
	pgInstance := pg.NewPG(ctx, connString, log) // Обновленный вызов
	if pgInstance == nil {
		log.Fatal().Msg("Ошибка подключения к базе данных")
	}
	defer pgInstance.Close() // Закрытие соединения с базой данных при завершении работы

	// Создание репозитория и сервиса
	repo := repository.NewBrandRepository(pgInstance.Pool()) // Передаем пул соединений
	brandService := service.NewBrandService(repo)            // Создание сервиса

	// Создание контроллеров для работы с запросами
	brandController := controller.NewBrandController(brandService) // Используем конструктор контроллера

	// Настройка роутера
	r := router.New()

	// Пример маршрутов с использованием контроллера
	r.POST("/brands", brandController.CreateBrand)
	r.PUT("/brands/{id}", brandController.UpdateBrand)
	r.GET("/brands/{id}", brandController.GetBrandByID)
	r.DELETE("/brands/{id}", brandController.SoftDeleteBrand)
	r.PATCH("/brands/{id}/restore", brandController.RestoreBrand)

	// Запуск сервера
	log.Info().Msg("Сервер запущен на порту :8080")
	err := fasthttp.ListenAndServe(":8080", r.Handler)
	if err != nil {
		log.Fatal().Err(err).Msg("Ошибка запуска сервера")
	}
}
