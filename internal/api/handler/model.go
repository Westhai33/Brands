package handler

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
	"strings"
)

type ModelHandler struct {
	ModelService *service.ModelService
}

func NewModelHandler(modelService *service.ModelService) *ModelHandler {
	return &ModelHandler{
		ModelService: modelService,
	}
}

func (api *ModelHandler) SetupRoutes(r *router.Router) {
	group := r.Group("/models")
	group.POST("/create", api.CreateModel)
	group.GET("/{id}", api.GetModelByID)
	group.PUT("/update/{id}", api.UpdateModel)
	group.DELETE("/delete/{id}", api.DeleteModel)
	group.POST("/restore/{id}", api.RestoreModel)
	group.GET("/all", api.GetAllModels)
	group.GET("/search", api.SearchModels)
}

// CreateModel godoc
// @Summary Создание новой модели
// @Description Эндпоинт для создания новой модели
// @Tags models
// @Accept json
// @Produce json
// @Param model body dto.Model true "Данные новой модели"
// @Success 200 {string} string "Model created successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Failed to create model"
// @Router /models/create [post]
func (api *ModelHandler) CreateModel(ctx *fasthttp.RequestCtx) {
	decoder := json.NewDecoder(bytes.NewReader(ctx.PostBody()))
	var model dto.Model
	err := decoder.Decode(&model)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}

	modelID, err := api.ModelService.Create(ctx, &model)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to create model: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBodyString(fmt.Sprintf("Model created successfully with ID: %d", modelID))
}

// GetModelByID godoc
// @Summary Получение модели по ID
// @Description Получение информации о модели по её уникальному ID
// @Tags models
// @Accept json
// @Produce json
// @Param id path int true "ID модели"
// @Success 200 {object} dto.Model "Модель найдена"
// @Failure 400 {string} string "Invalid ID format"
// @Failure 404 {string} string "Model not found"
// @Router /models/{id} [get]
func (api *ModelHandler) GetModelByID(ctx *fasthttp.RequestCtx) {
	idStr := ctx.UserValue("id").(string)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString("Invalid ID format")
		return
	}

	model, err := api.ModelService.GetByID(ctx, id)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusNotFound)
		ctx.Response.SetBodyString(fmt.Sprintf("Model not found: %v", err))
		return
	}

	data, err := json.Marshal(model)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to marshal model data: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBody(data)
}

// UpdateModel godoc
// @Summary Обновление модели по ID
// @Description Обновление данных модели по её ID
// @Tags models
// @Accept json
// @Produce json
// @Param id path int true "ID модели"
// @Param model body dto.Model true "Обновлённые данные модели"
// @Success 200 {string} string "Model updated successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 404 {string} string "Model not found"
// @Failure 500 {string} string "Failed to update model"
// @Router /models/update/{id} [put]
func (api *ModelHandler) UpdateModel(ctx *fasthttp.RequestCtx) {
	idStr := ctx.UserValue("id").(string)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString("Invalid ID format")
		return
	}

	decoder := json.NewDecoder(bytes.NewReader(ctx.PostBody()))
	var model dto.Model
	err = decoder.Decode(&model)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}

	model.ID = id

	err = api.ModelService.Update(ctx, &model)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to update model: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBodyString("Model updated successfully")
}

// DeleteModel godoc
// @Summary Мягкое удаление модели
// @Description Выполняет мягкое удаление модели по её ID
// @Tags models
// @Accept json
// @Produce json
// @Param id path int true "ID модели"
// @Success 200 {string} string "Model soft-deleted successfully"
// @Failure 400 {string} string "Invalid ID format"
// @Failure 500 {string} string "Failed to delete model"
// @Router /models/delete/{id} [delete]
func (api *ModelHandler) DeleteModel(ctx *fasthttp.RequestCtx) {
	idStr := ctx.UserValue("id").(string)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString("Invalid ID format")
		return
	}

	err = api.ModelService.SoftDelete(ctx, id)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to delete model: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBodyString("Model soft-deleted successfully")
}

// RestoreModel godoc
// @Summary Восстановление мягко удалённой модели
// @Description Восстановление модели, которая была удалена ранее
// @Tags models
// @Accept json
// @Produce json
// @Param id path int true "ID модели"
// @Success 200 {string} string "Model restored successfully"
// @Failure 400 {string} string "Invalid ID format"
// @Failure 500 {string} string "Failed to restore model"
// @Router /models/restore/{id} [post]
func (api *ModelHandler) RestoreModel(ctx *fasthttp.RequestCtx) {
	idStr := ctx.UserValue("id").(string)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString("Invalid ID format")
		return
	}

	err = api.ModelService.Restore(ctx, id)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to restore model: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBodyString("Model restored successfully")
}

// GetAllModels godoc
// @Summary Получение всех моделей
// @Description Возвращает все модели с возможностью фильтрации и сортировки
// @Tags models
// @Accept json
// @Produce json
// @Param filter query string false "Фильтрация по полю"
// @Param sort query string false "Сортировка по полю"
// @Param order query string false "Порядок сортировки (asc, desc)"
// @Success 200 {array} dto.Model "Список моделей"
// @Failure 500 {string} string "Failed to fetch models"
// @Router /models/all [get]
func (api *ModelHandler) GetAllModels(ctx *fasthttp.RequestCtx) {
	// Извлечение строки фильтров из query параметра "filter"
	filterStr := string(ctx.QueryArgs().Peek("filter"))

	// Парсим строку фильтров с помощью parseFilter
	filter, err := parseFilter(filterStr)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString(fmt.Sprintf("Invalid filter format: %v", err))
		return
	}

	// Добавляем дополнительные фильтры из query параметров
	if name := string(ctx.QueryArgs().Peek("name")); name != "" {
		filter["name"] = name
	}
	if brandID := string(ctx.QueryArgs().Peek("brand_id")); brandID != "" {
		if brandIDVal, err := strconv.ParseInt(brandID, 10, 64); err == nil {
			filter["brand_id"] = brandIDVal
		} else {
			ctx.Response.SetStatusCode(http.StatusBadRequest)
			ctx.Response.SetBodyString("Invalid brand_id value")
			return
		}
	}
	if isUpcoming := ctx.QueryArgs().Peek("is_upcoming"); len(isUpcoming) > 0 {
		if isUpcomingVal, err := strconv.ParseBool(string(isUpcoming)); err == nil {
			filter["is_upcoming"] = isUpcomingVal
		} else {
			ctx.Response.SetStatusCode(http.StatusBadRequest)
			ctx.Response.SetBodyString("Invalid is_upcoming value")
			return
		}
	}
	if isLimited := ctx.QueryArgs().Peek("is_limited"); len(isLimited) > 0 {
		if isLimitedVal, err := strconv.ParseBool(string(isLimited)); err == nil {
			filter["is_limited"] = isLimitedVal
		} else {
			ctx.Response.SetStatusCode(http.StatusBadRequest)
			ctx.Response.SetBodyString("Invalid is_limited value")
			return
		}
	}

	// Обработка сортировки
	sort := string(ctx.QueryArgs().Peek("sort"))
	if sort != "" {
		validSortFields := map[string]bool{
			"name":         true,
			"release_date": true,
			"created_at":   true,
			"updated_at":   true,
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
	models, err := api.ModelService.GetAll(ctx, filter, sort)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to fetch models: %v", err))
		return
	}

	// Преобразуем список моделей в JSON
	data, err := json.Marshal(models)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to marshal models data: %v", err))
		return
	}

	// Отправляем список моделей в ответ
	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBody(data)
}
