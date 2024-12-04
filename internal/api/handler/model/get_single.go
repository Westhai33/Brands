package model

import (
	"Brands/internal/api/handler"
	"context"
	"encoding/json"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/valyala/fasthttp"
	"net/http"
)

// GetModelByID godoc
// @Summary Получение модели по ID
// @Description Получение информации о модели по её уникальному ID
// @Tags models
// @Accept json
// @Produce json
// @Param id path string true "ID модели (UUIDv7)"
// @Success 200 {object} dto.Model "Модель найдена"
// @Failure 400 {string} string "Invalid ID format"
// @Failure 404 {string} string "Model not found"
// @Router /models/{id} [get]
func (api *ModelHandler) GetModelByID(ctx *fasthttp.RequestCtx) {
	var spanCtx context.Context
	spanCtx, ok := ctx.UserValue("traceContext").(context.Context)
	if !ok {
		spanCtx = ctx
	}

	span, spanCtx := opentracing.StartSpanFromContext(spanCtx, "ModelHandler.GetModelByID")
	defer span.Finish()

	// Извлечение и парсинг UUID из пути запроса
	id, err := handler.ExtractUUIDFromPath(ctx, "id")
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(
			log.String("event", "invalid_id"),
			log.Error(err),
		)
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString("Invalid ID format")
		return
	}

	model, err := api.ModelService.GetByID(spanCtx, id)
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
