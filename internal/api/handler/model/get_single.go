package model

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
	"time"
)

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
func (api *Handler) GetModelByID(ctx *fasthttp.RequestCtx) {
	var spanCtx context.Context
	spanCtx, ok := ctx.UserValue("traceContext").(context.Context)
	if !ok {
		spanCtx = ctx
	}
	spanCtx, cancel := context.WithTimeout(spanCtx, 5*time.Second)
	defer cancel()
	span, spanCtx := opentracing.StartSpanFromContext(spanCtx, "Handler.GetModelByID")
	defer span.Finish()

	idStr := ctx.UserValue("id").(string)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusBadRequest)
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
