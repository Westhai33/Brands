package model

import (
	"Brands/pkg/zerohook"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
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
// @Success 200 {array} dto.Model "Список моделей"
// @Failure 500 {string} string "Failed to fetch models"
// @Router /models/all [get]
func (api *ModelHandler) GetAllModels(ctx *fasthttp.RequestCtx) {
	var spanCtx context.Context
	spanCtx, ok := ctx.UserValue("traceContext").(context.Context)
	if !ok {
		spanCtx = ctx
	}
	span, spanCtx := opentracing.StartSpanFromContext(spanCtx, "ModelHandler.GetAllModels")
	defer span.Finish()

	models, err := api.ModelService.GetAll(spanCtx)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(
			log.String("event", "failed_to_fetch_models"),
			log.Error(err),
		)
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to fetch models: %v", err))
		return
	}

	// Преобразуем список моделей в JSON
	data, err := json.Marshal(models)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(
			log.String("event", "json_marshal_error"),
			log.Error(err),
			log.Object("models", models),
		)
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to marshal models data: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBody(data)
}

// ModelsFilter godoc
// @Summary Фильтрация моделей
// @Description Возвращает все модели с возможностью фильтрации и сортировки
// @Tags models
// @Accept json
// @Produce json
// @Param name query string false "Фильтр по имени модели"
// @Param brand_id query string false "Фильтр по идентификатору бренда"
// @Param popularity query integer false "Фильтр по популярности (целое число)"
// @Param is_limited query boolean false "Фильтр по признаку премиум-модели"
// @Param sort query string false "Поле сортировки (например, 'name', '-popularity')"
// @Success 200 {array} dto.Model "Список моделей"
// @Failure 500 {string} string "Failed to fetch models"
// @Router /models/filter [get]
func (api *ModelHandler) ModelsFilter(ctx *fasthttp.RequestCtx) {
	var spanCtx context.Context
	spanCtx, ok := ctx.UserValue("traceContext").(context.Context)
	if !ok {
		spanCtx = ctx
	}
	span, spanCtx := opentracing.StartSpanFromContext(spanCtx, "ModelHandler.GetAllModels")
	defer span.Finish()

	filter := make(map[string]any)

	if name := string(ctx.QueryArgs().Peek("name")); name != "" {
		filter["name"] = name
	}
	if brandID := ctx.QueryArgs().Peek("brand_id"); len(brandID) > 0 {
		_, err := uuid.Parse(string(brandID))
		zerohook.Logger.Info().Msg(string(brandID))

		if err != nil {
			span.SetTag("error", true)
			span.LogFields(
				log.String("event", "invalid_brand_id"),
				log.Error(err),
			)
			ctx.SetStatusCode(http.StatusBadRequest)
			ctx.Response.SetBodyString("Invalid brand_id format")
			return
		}
		filter["brand_id"] = string(brandID)
	}
	if isPremium := ctx.QueryArgs().Peek("is_limited"); len(isPremium) > 0 {
		if isPremiumVal, err := strconv.ParseBool(string(isPremium)); err == nil {
			filter["is_limited"] = isPremiumVal
		} else {
			message := "Invalid is_limited value"
			span.SetTag("error", true)
			span.LogFields(log.String("err", message))
			ctx.Response.SetStatusCode(http.StatusBadRequest)
			ctx.Response.SetBodyString("Invalid is_limited value")
			return
		}
	}

	sort := string(ctx.QueryArgs().Peek("sort"))
	if sort != "" {
		validSortFields := map[string]bool{
			"name":       true,
			"created_at": true,
			"updated_at": true,
		}
		sortField := strings.TrimPrefix(sort, "-")
		if !validSortFields[sortField] {
			message := fmt.Sprintf("Invalid sort field: %s", sortField)
			span.SetTag("error", true)
			span.LogFields(log.String("err", message))
			ctx.Response.SetStatusCode(http.StatusBadRequest)
			ctx.Response.SetBodyString(message)
			return
		}
	}

	models, err := api.ModelService.ModelsFilter(spanCtx, filter, sort)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(
			log.String("event", "failed_to_filter_models"),
			log.Error(err),
		)
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to fetch models: %v", err))
		return
	}

	data, err := json.Marshal(models)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(
			log.String("event", "json_marshal_error"),
			log.Error(err),
			log.Object("models", models),
		)
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to marshal models data: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBody(data)
}
