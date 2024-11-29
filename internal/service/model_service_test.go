package service_test

import (
	"Brands/internal/dto"
	"Brands/test/mocks"
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestModelService_Create_Success(t *testing.T) {
	mc := minimock.NewController(t)

	mockService := mocks.NewModelServiceInterfaceMock(mc)

	model := &dto.Model{Name: "Test Model"}
	mockService.CreateMock.Expect(context.Background(), model).Return(int64(1), nil)

	createdModelID, err := mockService.Create(context.Background(), model)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), createdModelID)
}

func TestModelService_Create_Error(t *testing.T) {
	mc := minimock.NewController(t)

	mockService := mocks.NewModelServiceInterfaceMock(mc)

	model := &dto.Model{Name: "Test Model"}
	mockService.CreateMock.Expect(context.Background(), model).Return(int64(0), assert.AnError)

	createdModelID, err := mockService.Create(context.Background(), model)

	assert.Error(t, err)
	assert.Equal(t, int64(0), createdModelID)
}

func TestModelService_GetAll_Success(t *testing.T) {
	mc := minimock.NewController(t)

	mockService := mocks.NewModelServiceInterfaceMock(mc)

	filter := "popular"
	sort := "name"
	order := "asc"
	mockService.GetAllMock.Expect(context.Background(), filter, sort, order).Return([]dto.Model{
		{ID: 1, Name: "Model A"},
		{ID: 2, Name: "Model B"},
	}, nil)

	models, err := mockService.GetAll(context.Background(), filter, sort, order)

	assert.NoError(t, err)
	assert.Len(t, models, 2)
	assert.Equal(t, "Model A", models[0].Name)
	assert.Equal(t, "Model B", models[1].Name)
}

func TestModelService_GetAll_Error(t *testing.T) {
	mc := minimock.NewController(t)

	mockService := mocks.NewModelServiceInterfaceMock(mc)

	filter := "popular"
	sort := "name"
	order := "asc"
	mockService.GetAllMock.Expect(context.Background(), filter, sort, order).Return(nil, assert.AnError)

	models, err := mockService.GetAll(context.Background(), filter, sort, order)

	assert.Error(t, err)
	assert.Nil(t, models)
}

func TestModelService_GetByID_Success(t *testing.T) {
	mc := minimock.NewController(t)

	mockService := mocks.NewModelServiceInterfaceMock(mc)

	modelID := int64(1)
	mockService.GetByIDMock.Expect(context.Background(), modelID).Return(&dto.Model{ID: modelID, Name: "Model A"}, nil)

	model, err := mockService.GetByID(context.Background(), modelID)

	assert.NoError(t, err)
	assert.Equal(t, "Model A", model.Name)
}

func TestModelService_GetByID_Error(t *testing.T) {
	mc := minimock.NewController(t)

	mockService := mocks.NewModelServiceInterfaceMock(mc)

	modelID := int64(1)
	mockService.GetByIDMock.Expect(context.Background(), modelID).Return(nil, assert.AnError)

	model, err := mockService.GetByID(context.Background(), modelID)

	assert.Error(t, err)
	assert.Nil(t, model)
}

func TestModelService_Restore_Success(t *testing.T) {
	mc := minimock.NewController(t)

	mockService := mocks.NewModelServiceInterfaceMock(mc)

	modelID := int64(1)
	mockService.RestoreMock.Expect(context.Background(), modelID).Return(nil)

	err := mockService.Restore(context.Background(), modelID)

	assert.NoError(t, err)
}

func TestModelService_Restore_Error(t *testing.T) {
	mc := minimock.NewController(t)

	mockService := mocks.NewModelServiceInterfaceMock(mc)

	modelID := int64(1)
	mockService.RestoreMock.Expect(context.Background(), modelID).Return(assert.AnError)

	err := mockService.Restore(context.Background(), modelID)

	assert.Error(t, err)
}

func TestModelService_SoftDelete_Success(t *testing.T) {
	mc := minimock.NewController(t)

	mockService := mocks.NewModelServiceInterfaceMock(mc)

	modelID := int64(1)
	mockService.SoftDeleteMock.Expect(context.Background(), modelID).Return(nil)

	err := mockService.SoftDelete(context.Background(), modelID)

	assert.NoError(t, err)
}

func TestModelService_SoftDelete_Error(t *testing.T) {
	mc := minimock.NewController(t)

	mockService := mocks.NewModelServiceInterfaceMock(mc)

	modelID := int64(1)
	mockService.SoftDeleteMock.Expect(context.Background(), modelID).Return(assert.AnError)

	err := mockService.SoftDelete(context.Background(), modelID)

	assert.Error(t, err)
}

func TestModelService_Update_Success(t *testing.T) {
	mc := minimock.NewController(t)

	mockService := mocks.NewModelServiceInterfaceMock(mc)

	model := &dto.Model{ID: 1, Name: "Updated Model"}
	mockService.UpdateMock.Expect(context.Background(), model).Return(nil)

	err := mockService.Update(context.Background(), model)

	assert.NoError(t, err)
}

func TestModelService_Update_Error(t *testing.T) {
	mc := minimock.NewController(t)

	mockService := mocks.NewModelServiceInterfaceMock(mc)

	model := &dto.Model{ID: 1, Name: "Updated Model"}
	mockService.UpdateMock.Expect(context.Background(), model).Return(assert.AnError)

	err := mockService.Update(context.Background(), model)

	assert.Error(t, err)
}
