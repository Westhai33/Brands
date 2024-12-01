package service_test

import (
	"Brands/internal/dto"
	"Brands/test/mocks"
	"context"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestBrandService_Create_Success(t *testing.T) {
	mc := minimock.NewController(t)

	mockService := mocks.NewBrandServiceInterfaceMock(mc)

	brand := &dto.Brand{Name: "Test Brand"}
	mockService.CreateMock.Expect(context.Background(), brand).Return(int64(1), nil)

	createdBrandID, err := mockService.Create(context.Background(), brand)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), createdBrandID)
}

func TestBrandService_Create_Error(t *testing.T) {
	mc := minimock.NewController(t)

	mockService := mocks.NewBrandServiceInterfaceMock(mc)

	brand := &dto.Brand{Name: "Test Brand"}
	mockService.CreateMock.Expect(context.Background(), brand).Return(int64(0), assert.AnError)

	createdBrandID, err := mockService.Create(context.Background(), brand)

	assert.Error(t, err)
	assert.Equal(t, int64(0), createdBrandID)
}

func TestBrandService_GetAll_Success(t *testing.T) {
	mc := minimock.NewController(t)

	mockService := mocks.NewBrandServiceInterfaceMock(mc)

	filter := "premium"
	sort := "popularity"
	mockService.GetAllMock.Expect(context.Background(), filter, sort).Return([]dto.Brand{
		{ID: 1, Name: "Brand A"},
		{ID: 2, Name: "Brand B"},
	}, nil)

	brands, err := mockService.GetAll(context.Background(), filter, sort)

	assert.NoError(t, err)
	assert.Len(t, brands, 2)
	assert.Equal(t, "Brand A", brands[0].Name)
	assert.Equal(t, "Brand B", brands[1].Name)
}

func TestBrandService_GetAll_Error(t *testing.T) {
	mc := minimock.NewController(t)

	mockService := mocks.NewBrandServiceInterfaceMock(mc)

	filter := "premium"
	sort := "popularity"
	mockService.GetAllMock.Expect(context.Background(), filter, sort).Return(nil, assert.AnError)

	brands, err := mockService.GetAll(context.Background(), filter, sort)

	assert.Error(t, err)
	assert.Nil(t, brands)
}

func TestBrandService_Update_Success(t *testing.T) {
	mc := minimock.NewController(t)

	mockService := mocks.NewBrandServiceInterfaceMock(mc)

	brand := &dto.Brand{ID: 1, Name: "Updated Brand"}
	mockService.UpdateMock.Expect(context.Background(), brand).Return(nil)

	err := mockService.Update(context.Background(), brand)

	assert.NoError(t, err)
}

func TestBrandService_Update_Error(t *testing.T) {
	mc := minimock.NewController(t)

	mockService := mocks.NewBrandServiceInterfaceMock(mc)

	brand := &dto.Brand{ID: 1, Name: "Updated Brand"}
	mockService.UpdateMock.Expect(context.Background(), brand).Return(assert.AnError)

	err := mockService.Update(context.Background(), brand)

	assert.Error(t, err)
}

func TestBrandService_SoftDelete_Success(t *testing.T) {
	mc := minimock.NewController(t)

	mockService := mocks.NewBrandServiceInterfaceMock(mc)

	mockService.SoftDeleteMock.Expect(context.Background(), int64(1)).Return(nil)

	err := mockService.SoftDelete(context.Background(), int64(1))

	assert.NoError(t, err)
}

func TestBrandService_SoftDelete_Error(t *testing.T) {
	mc := minimock.NewController(t)

	mockService := mocks.NewBrandServiceInterfaceMock(mc)

	mockService.SoftDeleteMock.Expect(context.Background(), int64(1)).Return(assert.AnError)

	err := mockService.SoftDelete(context.Background(), int64(1))

	assert.Error(t, err)
}

func TestBrandService_Restore_Success(t *testing.T) {
	mc := minimock.NewController(t)

	mockService := mocks.NewBrandServiceInterfaceMock(mc)

	mockService.RestoreMock.Expect(context.Background(), int64(1)).Return(nil)

	err := mockService.Restore(context.Background(), int64(1))

	assert.NoError(t, err)
}

func TestBrandService_Restore_Error(t *testing.T) {
	mc := minimock.NewController(t)

	mockService := mocks.NewBrandServiceInterfaceMock(mc)

	mockService.RestoreMock.Expect(context.Background(), int64(1)).Return(assert.AnError)

	err := mockService.Restore(context.Background(), int64(1))

	assert.Error(t, err)
}
