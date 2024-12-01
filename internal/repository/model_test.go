package repository_test

//
//import (
//	"Brands/internal/dto"
//	"Brands/internal/repository"
//	"context"
//	"testing"
//	"time"
//
//	"github.com/stretchr/testify/assert"
//)
//
//// Тест для создания модели
//func TestCreateModel(t *testing.T) {
//	pool := setupTestDB(t)
//	defer teardownTestDB(pool)
//
//	repo := repository.NewModelRepository(pool)
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
//	brandID, err := repository.NewBrandRepository(pool).Create(context.Background(), brand)
//	assert.NoError(t, err)
//
//	model := &dto.Model{
//		BrandID:     brandID,
//		Name:        "Test Model",
//		ReleaseDate: time.Now(),
//		IsUpcoming:  true,
//		IsLimited:   false,
//	}
//
//	modelID, err := repo.Create(context.Background(), model)
//	assert.NoError(t, err)
//	assert.NotZero(t, modelID)
//
//	createdModel, err := repo.GetByID(context.Background(), modelID)
//	assert.NoError(t, err)
//	assert.Equal(t, model.Name, createdModel.Name)
//}
//
//// Тест для обновления модели
//func TestUpdateModel(t *testing.T) {
//	pool := setupTestDB(t)
//	defer teardownTestDB(pool)
//
//	repo := repository.NewModelRepository(pool)
//
//	brand := &dto.Brand{
//		Name:          "Brand for Update",
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
//	brandID, err := repository.NewBrandRepository(pool).Create(context.Background(), brand)
//	assert.NoError(t, err)
//
//	model := &dto.Model{
//		BrandID:     brandID,
//		Name:        "Model to Update",
//		ReleaseDate: time.Now(),
//		IsUpcoming:  true,
//		IsLimited:   false,
//	}
//
//	modelID, err := repo.Create(context.Background(), model)
//	assert.NoError(t, err)
//
//	model.Name = "Updated Model"
//	err = repo.Update(context.Background(), model)
//	assert.NoError(t, err)
//
//	updatedModel, err := repo.GetByID(context.Background(), modelID)
//	assert.NoError(t, err)
//	assert.Equal(t, "Updated Model", updatedModel.Name)
//}
//
//// Тест для мягкого удаления модели
//func TestSoftDeleteModel(t *testing.T) {
//	pool := setupTestDB(t)
//	defer teardownTestDB(pool)
//
//	repo := repository.NewModelRepository(pool)
//
//	brand := &dto.Brand{
//		Name:          "Brand for Deletion",
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
//	brandID, err := repository.NewBrandRepository(pool).Create(context.Background(), brand)
//	assert.NoError(t, err)
//
//	model := &dto.Model{
//		BrandID:     brandID,
//		Name:        "Model to Soft Delete",
//		ReleaseDate: time.Now(),
//		IsUpcoming:  true,
//		IsLimited:   false,
//	}
//
//	modelID, err := repo.Create(context.Background(), model)
//	assert.NoError(t, err)
//
//	err = repo.SoftDelete(context.Background(), modelID)
//	assert.NoError(t, err)
//
//	softDeletedModel, err := repo.GetByID(context.Background(), modelID)
//	assert.NoError(t, err)
//	assert.Nil(t, softDeletedModel) // Теперь модель должна быть nil, так как она была помечена как удалённая
//}
//
//// Тест для восстановления мягко удалённой модели
//func TestRestoreModel(t *testing.T) {
//	pool := setupTestDB(t)
//	defer teardownTestDB(pool)
//
//	repo := repository.NewModelRepository(pool)
//
//	brand := &dto.Brand{
//		Name:          "Brand to Restore Model",
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
//	brandID, err := repository.NewBrandRepository(pool).Create(context.Background(), brand)
//	assert.NoError(t, err)
//
//	model := &dto.Model{
//		BrandID:     brandID,
//		Name:        "Model to Restore",
//		ReleaseDate: time.Now(),
//		IsUpcoming:  true,
//		IsLimited:   false,
//	}
//
//	modelID, err := repo.Create(context.Background(), model)
//	assert.NoError(t, err)
//
//	err = repo.SoftDelete(context.Background(), modelID)
//	assert.NoError(t, err)
//
//	err = repo.Restore(context.Background(), modelID)
//	assert.NoError(t, err)
//
//	restoredModel, err := repo.GetByID(context.Background(), modelID)
//	assert.NoError(t, err)
//	assert.False(t, restoredModel.IsDeleted)
//}
