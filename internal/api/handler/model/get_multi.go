package model

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
	"strings"
	"time"
)

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
func (api *Handler) GetAllModels(ctx *fasthttp.RequestCtx) {
	var spanCtx context.Context
	spanCtx, ok := ctx.UserValue("traceContext").(context.Context)
	if !ok {
		spanCtx = ctx
	}
	spanCtx, cancel := context.WithTimeout(spanCtx, 5*time.Second)
	defer cancel()
	span, spanCtx := opentracing.StartSpanFromContext(spanCtx, "Handler.GetAllModels")
	defer span.Finish()

	filterStr := string(ctx.QueryArgs().Peek("filter"))

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

		sortField := strings.TrimPrefix(sort, "-")
		if !validSortFields[sortField] {
			ctx.Response.SetStatusCode(http.StatusBadRequest)
			ctx.Response.SetBodyString(fmt.Sprintf("Invalid sort field: %s", sortField))
			return
		}
	}

	models, err := api.ModelService.GetAll(spanCtx, filter, sort)
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

func parseFilter(filterStr string) (map[string]interface{}, error) {
	filter := make(map[string]interface{})
	if filterStr != "" {
		pairs := strings.Split(filterStr, "&")
		for _, pair := range pairs {
			parts := strings.Split(pair, "=")
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid filter format")
			}
			filter[parts[0]] = parts[1]
		}
	}
	return filter, nil
}
