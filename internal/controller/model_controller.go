package controller

import (
	"Brands/internal/dto"
	"Brands/internal/service"
	"context"
	"fmt"
	"strconv"
	"time"
)

// CreateModel создает новую модель.
func CreateModel(ctx context.Context, modelService *service.ModelService, brandService *service.BrandService, name, brandID string, isUpcoming, isLimited bool) (int64, error) {
	if name == "" {
		return 0, fmt.Errorf("name is required")
	}

	brandIDInt, err := strconv.ParseInt(brandID, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid brandID: %w", err)
	}

	_, err = brandService.GetByID(ctx, brandIDInt)
	if err != nil {
		return 0, fmt.Errorf("brand with ID %d not found: %w", brandIDInt, err)
	}

	// Создание модели через сервис
	model := &dto.Model{
		Name:        name,
		BrandID:     brandIDInt,
		IsUpcoming:  isUpcoming,
		IsLimited:   isLimited,
		ReleaseDate: time.Now(),
	}

	// Вызов сервиса для создания модели
	modelID, err := modelService.Create(ctx, model)
	if err != nil {
		return 0, fmt.Errorf("error creating model: %w", err)
	}

	return modelID, nil
}

// GetModelByID получает модель по ID.
func GetModelByID(ctx context.Context, modelService *service.ModelService, modelID int64) (*dto.Model, error) {
	model, err := modelService.GetByID(ctx, modelID)
	if err != nil {
		return nil, fmt.Errorf("error getting model with ID %d: %w", modelID, err)
	}

	return model, nil
}

// GetAllModels получает все модели с фильтрацией, сортировкой и порядком.
func GetAllModels(ctx context.Context, modelService *service.ModelService, filter, sort, order string) ([]dto.Model, error) {
	models, err := modelService.GetAll(ctx, filter, sort, order)
	if err != nil {
		return nil, fmt.Errorf("error getting all models: %w", err)
	}

	return models, nil
}

// UpdateModel обновляет данные модели с проверкой существования.
func UpdateModel(ctx context.Context, modelService *service.ModelService, brandService *service.BrandService, modelID int64, name, brandID string, isUpcoming, isLimited bool) error {
	model, err := modelService.GetByID(ctx, modelID)
	if err != nil {
		return fmt.Errorf("model with ID %d not found: %w", modelID, err)
	}

	brandIDInt, err := strconv.ParseInt(brandID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid brandID: %w", err)
	}

	_, err = brandService.GetByID(ctx, brandIDInt)
	if err != nil {
		return fmt.Errorf("brand with ID %d not found: %w", brandIDInt, err)
	}

	model.Name = name
	model.BrandID = brandIDInt
	model.IsUpcoming = isUpcoming
	model.IsLimited = isLimited
	model.UpdatedAt = time.Now()

	err = modelService.Update(ctx, model)
	if err != nil {
		return fmt.Errorf("error updating model with ID %d: %w", modelID, err)
	}

	return nil
}

// SoftDeleteModel мягко удаляет модель с проверкой существования.
func SoftDeleteModel(ctx context.Context, modelService *service.ModelService, modelID int64) error {
	_, err := modelService.GetByID(ctx, modelID)
	if err != nil {
		return fmt.Errorf("model with ID %d not found: %w", modelID, err)
	}

	err = modelService.SoftDelete(ctx, modelID)
	if err != nil {
		return fmt.Errorf("error soft-deleting model with ID %d: %w", modelID, err)
	}

	return nil
}

// RestoreModel восстанавливает мягко удалённую модель.
func RestoreModel(ctx context.Context, modelService *service.ModelService, modelID int64) error {
	_, err := modelService.GetByID(ctx, modelID)
	if err != nil {
		return fmt.Errorf("model with ID %d not found: %w", modelID, err)
	}

	err = modelService.Restore(ctx, modelID)
	if err != nil {
		return fmt.Errorf("error restoring model with ID %d: %w", modelID, err)
	}

	return nil
}
