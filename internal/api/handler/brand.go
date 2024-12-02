package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"Brands/internal/dto"
	"Brands/internal/service"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	_ "Brands/docs"
)

type BrandHandler struct {
	BrandService *service.BrandService
}

func NewBrandHandler(brandService *service.BrandService) *BrandHandler {
	return &BrandHandler{
		BrandService: brandService,
	}
}

func (api *BrandHandler) SetupRoutes(r *router.Router) {
	group := r.Group("/brands")
	group.POST("/create", api.CreateBrand)
	group.GET("/{id}", api.GetBrandByID)
	group.PUT("/update/{id}", api.UpdateBrand)
	group.DELETE("/delete/{id}", api.DeleteBrand)
	group.POST("/restore/{id}", api.RestoreBrand)
	group.GET("/all", api.GetAllBrands)
}

// CreateBrand godoc
// @Summary Создание нового бренда
// @Description Эндпоинт для создания нового бренда
// @Tags brands
// @Accept json
// @Produce json
// @Param brand body dto.Brand true "Данные нового бренда"
// @Success 200 {string} string "Brand created successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Failed to create brand"
// @Router /brands/create [post]
func (api *BrandHandler) CreateBrand(ctx *fasthttp.RequestCtx) {
	decoder := json.NewDecoder(bytes.NewReader(ctx.PostBody()))
	var brand dto.Brand
	err := decoder.Decode(&brand)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}

	brandID, err := api.BrandService.Create(ctx, &brand)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to create brand: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBodyString(fmt.Sprintf("Brand created successfully with ID: %d", brandID))
}

// GetBrandByID godoc
// @Summary Получение бренда по ID
// @Description Получение информации о бренде по его уникальному ID
// @Tags brands
// @Accept json
// @Produce json
// @Param id path int true "ID бренда"
// @Success 200 {object} dto.Brand "Бренд найден"
// @Failure 400 {string} string "Invalid ID format"
// @Failure 404 {string} string "Brand not found"
// @Router /brands/{id} [get]
func (api *BrandHandler) GetBrandByID(ctx *fasthttp.RequestCtx) {
	idStr := ctx.UserValue("id").(string)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString("Invalid ID format")
		return
	}

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

// UpdateBrand godoc
// @Summary Обновление бренда по ID
// @Description Обновление данных бренда по его ID
// @Tags brands
// @Accept json
// @Produce json
// @Param id path int true "ID бренда"
// @Param brand body dto.Brand true "Обновлённые данные бренда"
// @Success 200 {string} string "Brand updated successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 404 {string} string "Brand not found"
// @Failure 500 {string} string "Failed to update brand"
// @Router /brands/update/{id} [put]
func (api *BrandHandler) UpdateBrand(ctx *fasthttp.RequestCtx) {
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

	brand.ID = id

	err = api.BrandService.Update(ctx, &brand)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to update brand: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBodyString("Brand updated successfully")
}

// DeleteBrand godoc
// @Summary Мягкое удаление бренда
// @Description Выполняет мягкое удаление бренда по его ID
// @Tags brands
// @Accept json
// @Produce json
// @Param id path int true "ID бренда"
// @Success 200 {string} string "Brand soft-deleted successfully"
// @Failure 400 {string} string "Invalid ID format"
// @Failure 500 {string} string "Failed to delete brand"
// @Router /brands/delete/{id} [delete]
func (api *BrandHandler) DeleteBrand(ctx *fasthttp.RequestCtx) {
	idStr := ctx.UserValue("id").(string)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString("Invalid ID format")
		return
	}

	err = api.BrandService.SoftDelete(ctx, id)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to delete brand: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBodyString("Brand soft-deleted successfully")
}

// RestoreBrand godoc
// @Summary Восстановление мягко удалённого бренда
// @Description Восстановление бренда, который был удалён ранее (мягкое удаление)
// @Tags brands
// @Accept json
// @Produce json
// @Param id path int true "ID бренда"
// @Success 200 {string} string "Brand restored successfully"
// @Failure 400 {string} string "Invalid ID format"
// @Failure 500 {string} string "Failed to restore brand"
// @Router /brands/restore/{id} [post]
func (api *BrandHandler) RestoreBrand(ctx *fasthttp.RequestCtx) {
	idStr := ctx.UserValue("id").(string)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString("Invalid ID format")
		return
	}

	err = api.BrandService.Restore(ctx, id)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to restore brand: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBodyString("Brand restored successfully")
}

// GetAllBrands godoc
// @Summary Получение всех брендов
// @Description Возвращает все бренды с возможностью фильтрации и сортировки
// @Tags brands
// @Accept json
// @Produce json
// @Param filter query string false "Фильтрация по полю"
// @Param sort query string false "Сортировка по полю"
// @Success 200 {array} dto.Brand "Список брендов"
// @Failure 500 {string} string "Failed to fetch brands"
// @Router /brands/all [get]
func (api *BrandHandler) GetAllBrands(ctx *fasthttp.RequestCtx) {
	// Извлечение фильтров из query параметров
	filter := make(map[string]interface{})

	// Добавляем фильтры на основе входных параметров
	if name := string(ctx.QueryArgs().Peek("name")); name != "" {
		filter["name"] = name
	}
	if originCountry := string(ctx.QueryArgs().Peek("origin_country")); originCountry != "" {
		filter["origin_country"] = originCountry
	}
	if popularity := ctx.QueryArgs().Peek("popularity"); len(popularity) > 0 {
		if popVal, err := strconv.Atoi(string(popularity)); err == nil {
			filter["popularity"] = popVal
		} else {
			ctx.Response.SetStatusCode(http.StatusBadRequest)
			ctx.Response.SetBodyString("Invalid popularity value")
			return
		}
	}
	if isPremium := ctx.QueryArgs().Peek("is_premium"); len(isPremium) > 0 {
		if isPremiumVal, err := strconv.ParseBool(string(isPremium)); err == nil {
			filter["is_premium"] = isPremiumVal
		} else {
			ctx.Response.SetStatusCode(http.StatusBadRequest)
			ctx.Response.SetBodyString("Invalid is_premium value")
			return
		}
	}

	// Обработка сортировки
	sort := string(ctx.QueryArgs().Peek("sort"))
	if sort != "" {
		validSortFields := map[string]bool{
			"name":           true,
			"popularity":     true,
			"founded_year":   true,
			"origin_country": true,
			"created_at":     true,
		}

		// Проверяем, является ли поле сортировки допустимым
		sortField := strings.TrimPrefix(sort, "-")
		if !validSortFields[sortField] {
			ctx.Response.SetStatusCode(http.StatusBadRequest)
			ctx.Response.SetBodyString(fmt.Sprintf("Invalid sort field: %s", sortField))
			return
		}
	}

	// Вызов метода GetAll из сервиса с фильтрами и сортировкой
	brands, err := api.BrandService.GetAll(ctx, filter, sort)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to fetch brands: %v", err))
		return
	}

	// Преобразуем список брендов в JSON
	data, err := json.Marshal(brands)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to marshal brands data: %v", err))
		return
	}

	// Отправляем список брендов в ответ
	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBody(data)
}

func parseFilter(filterStr string) (map[string]interface{}, error) {
	filter := make(map[string]interface{})
	if filterStr != "" {
		// Пример: filter="name=Apple&popularity=5"
		pairs := strings.Split(filterStr, "&")
		for _, pair := range pairs {
			parts := strings.Split(pair, "=")
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid filter format")
			}
			filter[parts[0]] = parts[1] // Параметры фильтра: ключ=значение
		}
	}
	return filter, nil
}
