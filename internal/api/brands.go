package api

import (
	"Brands/internal/dto"
	"Brands/internal/service"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
)

type BrandAPI struct {
	// Экспортируемое поле
	BrandService *service.BrandService
}

// Новый роутер для API
func (api *BrandAPI) SetupRoutes(r *router.Router) {
	group := r.Group("/brands")
	group.POST("/create", api.CreateBrand)
	group.GET("/{id}", api.GetBrandByID)
	group.PUT("/update/{id}", api.UpdateBrand)
	group.DELETE("/delete/{id}", api.DeleteBrand)
	group.POST("/restore/{id}", api.RestoreBrand)
	group.GET("/all", api.GetAllBrands)
	group.GET("/search", api.SearchBrands)
}

// Эндпоинт для создания нового бренда
func (api *BrandAPI) CreateBrand(ctx *fasthttp.RequestCtx) {
	decoder := json.NewDecoder(bytes.NewReader(ctx.PostBody()))
	var brand dto.Brand
	err := decoder.Decode(&brand)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}

	// Вызов метода Create из сервиса
	brandID, err := api.BrandService.Create(ctx, &brand)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to create brand: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBodyString(fmt.Sprintf("Brand created successfully with ID: %d", brandID))
}

// Эндпоинт для получения бренда по ID
func (api *BrandAPI) GetBrandByID(ctx *fasthttp.RequestCtx) {
	idStr := ctx.UserValue("id").(string)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString("Invalid ID format")
		return
	}

	// Вызов метода GetByID из сервиса
	brand, err := api.BrandService.GetByID(ctx, id)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusNotFound)
		ctx.Response.SetBodyString(fmt.Sprintf("Brand not found: %v", err))
		return
	}

	data, err := json.Marshal(brand)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to marshal brand data: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBody(data)
}

// Эндпоинт для обновления бренда
func (api *BrandAPI) UpdateBrand(ctx *fasthttp.RequestCtx) {
	idStr := ctx.UserValue("id").(string)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString("Invalid ID format")
		return
	}

	decoder := json.NewDecoder(bytes.NewReader(ctx.PostBody()))
	var brand dto.Brand
	err = decoder.Decode(&brand)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}

	brand.ID = id // Устанавливаем ID из URL в объект бренда

	// Вызов метода Update из сервиса
	err = api.BrandService.Update(ctx, &brand)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to update brand: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBodyString("Brand updated successfully")
}

// Эндпоинт для мягкого удаления бренда
func (api *BrandAPI) DeleteBrand(ctx *fasthttp.RequestCtx) {
	idStr := ctx.UserValue("id").(string)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString("Invalid ID format")
		return
	}

	// Вызов метода SoftDelete из сервиса
	err = api.BrandService.SoftDelete(ctx, id)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to delete brand: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBodyString("Brand soft-deleted successfully")
}

// Эндпоинт для восстановления мягко удалённого бренда
func (api *BrandAPI) RestoreBrand(ctx *fasthttp.RequestCtx) {
	idStr := ctx.UserValue("id").(string)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString("Invalid ID format")
		return
	}

	// Вызов метода Restore из сервиса
	err = api.BrandService.Restore(ctx, id)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to restore brand: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBodyString("Brand restored successfully")
}

// Эндпоинт для получения всех брендов с фильтрацией и сортировкой
func (api *BrandAPI) GetAllBrands(ctx *fasthttp.RequestCtx) {
	filter := string(ctx.QueryArgs().Peek("filter"))
	sort := string(ctx.QueryArgs().Peek("sort"))

	// Вызов метода GetAll из сервиса
	brands, err := api.BrandService.GetAll(ctx, filter, sort)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to fetch brands: %v", err))
		return
	}

	data, err := json.Marshal(brands)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to marshal brands data: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBody(data)
}

// Эндпоинт для поиска брендов по фильтру
func (api *BrandAPI) SearchBrands(ctx *fasthttp.RequestCtx) {
	filter := string(ctx.QueryArgs().Peek("filter"))
	sort := string(ctx.QueryArgs().Peek("sort"))

	// Вызов метода GetAll с фильтром и сортировкой
	brands, err := api.BrandService.GetAll(ctx, filter, sort)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to search brands: %v", err))
		return
	}

	data, err := json.Marshal(brands)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to marshal search results: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBody(data)
}
