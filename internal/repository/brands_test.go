package repository_test

//
//import (
//	"Brands/internal/dto"
//	"Brands/internal/repository"
//	"context"
//	"log"
//	"testing"
//
//	"github.com/jackc/pgx/v5/pgxpool"
//	"github.com/stretchr/testify/assert"
//)
//
//// Инициализация подключения к базе данных
//func setupTestDB(t *testing.T) *pgxpool.Pool {
//	connStr := "postgres://brands:pgpwdbrands@localhost:5432/brands?sslmode=disable"
//	var err error
//	pool, err := pgxpool.New(context.Background(), connStr)
//	assert.NoError(t, err, "unable to connect to database")
//
//	// Проверим, что таблица brands_test существует, если нет, создадим её
//	var tableExists bool
//	err = pool.QueryRow(context.Background(), "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'brands_test')").Scan(&tableExists)
//	assert.NoError(t, err)
//	if !tableExists {
//		// Создание таблицы brands_test, если она отсутствует
//		_, err = pool.Exec(context.Background(), `
//			CREATE TABLE brands_test (
//				id BIGSERIAL PRIMARY KEY,
//				name VARCHAR(255) NOT NULL,
//				link VARCHAR(255),
//				description TEXT,
//				logo_url VARCHAR(255),
//				cover_image_url VARCHAR(255),
//				founded_year INTEGER,
//				origin_country VARCHAR(100),
//				popularity INTEGER DEFAULT 0,
//				is_premium BOOLEAN DEFAULT FALSE,
//				is_upcoming BOOLEAN DEFAULT FALSE,
//				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
//			);
//		`)
//		assert.NoError(t, err, "unable to create table brands_test")
//	}
//
//	// Проверим, что таблица models_test существует, если нет, создадим её
//	err = pool.QueryRow(context.Background(), "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'models_test')").Scan(&tableExists)
//	assert.NoError(t, err)
//	if !tableExists {
//		// Создание таблицы models_test, если она отсутствует
//		_, err = pool.Exec(context.Background(), `
//			CREATE TABLE models_test (
//				id BIGSERIAL PRIMARY KEY,
//				brand_id BIGINT NOT NULL,
//				name VARCHAR(255) NOT NULL,
//				release_date DATE,
//				is_upcoming BOOLEAN DEFAULT FALSE,
//				is_limited BOOLEAN DEFAULT FALSE,
//				is_deleted BOOLEAN DEFAULT FALSE,
//				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//				CONSTRAINT fk_brand FOREIGN KEY(brand_id) REFERENCES brands_test(id) ON DELETE CASCADE
//			);
//		`)
//		assert.NoError(t, err, "unable to create table models_test")
//	}
//
//	return pool
//}
//
//// Очистка базы данных после тестов
//func teardownTestDB(pool *pgxpool.Pool) {
//	_, err := pool.Exec(context.Background(), "DROP TABLE IF EXISTS models_test CASCADE")
//	if err != nil {
//		log.Printf("Error dropping models_test table: %v", err)
//	}
//
//	_, err = pool.Exec(context.Background(), "DROP TABLE IF EXISTS brands_test CASCADE")
//	if err != nil {
//		log.Printf("Error dropping brands_test table: %v", err)
//	}
//
//	pool.Close()
//}
//
//// Тест для создания бренда
//func TestCreateBrand(t *testing.T) {
//	pool := setupTestDB(t)
//	defer teardownTestDB(pool)
//
//	repo := repository.NewBrandRepository(pool)
//
//	brand := &dto.Brand{
//		Name:          "Test Brand",
//		Link:          "https://example.com",
//		Description:   "Test Description",
//		LogoURL:       "https://example.com/logo.png",
//		CoverImageURL: "https://example.com/cover.jpg",
//		FoundedYear:   2023,
//		OriginCountry: "USA",
//		Popularity:    10,
//		IsPremium:     false,
//		IsUpcoming:    true,
//	}
//
//	id, err := repo.Create(context.Background(), brand)
//	assert.NoError(t, err)
//	assert.NotZero(t, id)
//
//	createdBrand, err := repo.GetByID(context.Background(), id)
//	assert.NoError(t, err)
//	assert.Equal(t, brand.Name, createdBrand.Name)
//}
//
//// Тест для обновления бренда
//func TestUpdateBrand(t *testing.T) {
//	pool := setupTestDB(t)
//	defer teardownTestDB(pool)
//
//	repo := repository.NewBrandRepository(pool)
//
//	brand := &dto.Brand{
//		Name:          "Brand to Update",
//		Link:          "https://example.com",
//		Description:   "Description",
//		LogoURL:       "https://example.com/logo.png",
//		CoverImageURL: "https://example.com/cover.jpg",
//		FoundedYear:   2023,
//		OriginCountry: "USA",
//		Popularity:    5,
//		IsPremium:     false,
//		IsUpcoming:    true,
//	}
//
//	id, err := repo.Create(context.Background(), brand)
//	assert.NoError(t, err)
//
//	brand.Name = "Updated Brand"
//	err = repo.Update(context.Background(), brand)
//	assert.NoError(t, err)
//
//	updatedBrand, err := repo.GetByID(context.Background(), id)
//	assert.NoError(t, err)
//	assert.Equal(t, "Updated Brand", updatedBrand.Name)
//}
//
//// Тест для мягкого удаления бренда
//func TestSoftDeleteBrand(t *testing.T) {
//	pool := setupTestDB(t)
//	defer teardownTestDB(pool)
//
//	repo := repository.NewBrandRepository(pool)
//
//	// Создаем новый бренд
//	brand := &dto.Brand{
//		Name:          "Brand to Delete",
//		Link:          "https://example.com",
//		Description:   "Description",
//		LogoURL:       "https://example.com/logo.png",
//		CoverImageURL: "https://example.com/cover.jpg",
//		FoundedYear:   2023,
//		OriginCountry: "USA",
//		Popularity:    5,
//		IsPremium:     false,
//		IsUpcoming:    true,
//	}
//
//	id, err := repo.Create(context.Background(), brand)
//	assert.NoError(t, err)
//
//	err = repo.SoftDelete(context.Background(), id)
//	assert.NoError(t, err)
//
//	softDeletedBrand, err := repo.GetByID(context.Background(), id)
//	assert.Error(t, err) // Ошибка должна быть, так как бренд помечен как удалённый
//	assert.Nil(t, softDeletedBrand)
//}
//
//// Тест для восстановления мягко удалённого бренда
//func TestRestoreBrand(t *testing.T) {
//	pool := setupTestDB(t)
//	defer teardownTestDB(pool)
//
//	repo := repository.NewBrandRepository(pool)
//
//	brand := &dto.Brand{
//		Name:          "Brand to Restore",
//		Link:          "https://example.com",
//		Description:   "Description",
//		LogoURL:       "https://example.com/logo.png",
//		CoverImageURL: "https://example.com/cover.jpg",
//		FoundedYear:   2023,
//		OriginCountry: "USA",
//		Popularity:    5,
//		IsPremium:     false,
//		IsUpcoming:    true,
//	}
//
//	id, err := repo.Create(context.Background(), brand)
//	assert.NoError(t, err)
//
//	err = repo.SoftDelete(context.Background(), id)
//	assert.NoError(t, err)
//
//	err = repo.Restore(context.Background(), id)
//	assert.NoError(t, err)
//
//	restoredBrand, err := repo.GetByID(context.Background(), id)
//	assert.NoError(t, err)
//	assert.False(t, restoredBrand.IsDeleted)
//}
//
//// Тест для получения всех брендов
//func TestGetAllBrands(t *testing.T) {
//	pool := setupTestDB(t)
//	defer teardownTestDB(pool)
//
//	repo := repository.NewBrandRepository(pool)
//
//	brands, err := repo.GetAll(context.Background(), "", "popularity")
//	assert.NoError(t, err)
//	assert.NotEmpty(t, brands)
//}
