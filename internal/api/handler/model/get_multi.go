package model

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
	"strings"
)

// GetAllModels godoc
// @Summary Получение всех моделей
// @Description Возвращает все модели с возможностью фильтрации и сортировки
// @Tags models
// @Accept json
// @Produce json
// @Param name query string false "Фильтр по имени модели"
// @Param brand_id query integer false "Фильтр по идентификатору бренда"
// @Param category query string false "Фильтр по категории модели"
// @Param popularity query integer false "Фильтр по популярности (целое число)"
// @Param is_premium query boolean false "Фильтр по признаку премиум-модели"
// @Param sort query string false "Поле сортировки (например, 'name', '-popularity')"
// @Success 200 {array} dto.Model "Список моделей"
// @Failure 500 {string} string "Failed to fetch models"
// @Router /models/all [get]
func (api *ModelHandler) GetAllModels(ctx *fasthttp.RequestCtx) {
	filter := make(map[string]interface{})

	if name := string(ctx.QueryArgs().Peek("name")); name != "" {
		filter["name"] = name
	}
	if brandID := ctx.QueryArgs().Peek("brand_id"); len(brandID) > 0 {
		if brandIDVal, err := strconv.Atoi(string(brandID)); err == nil {
			filter["brand_id"] = brandIDVal
		} else {
			ctx.Response.SetStatusCode(http.StatusBadRequest)
			ctx.Response.SetBodyString("Invalid brand_id value")
			return
		}
	}
	if category := string(ctx.QueryArgs().Peek("category")); category != "" {
		filter["category"] = category
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

	sort := string(ctx.QueryArgs().Peek("sort"))
	if sort != "" {
		validSortFields := map[string]bool{
			"name":       true,
			"popularity": true,
			"created_at": true,
			"updated_at": true,
		}

		sortField := strings.TrimPrefix(sort, "-")
		if !validSortFields[sortField] {
			ctx.Response.SetStatusCode(http.StatusBadRequest)
			ctx.Response.SetBodyString(fmt.Sprintf("Invalid sort field: %s", sortField))
			return
		}
	}

	models, err := api.ModelService.GetAll(ctx, filter, sort)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to fetch models: %v", err))
		return
	}

	data, err := json.Marshal(models)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to marshal models data: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBody(data)
}
