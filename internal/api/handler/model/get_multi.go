package model

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/valyala/fasthttp"
	"net/http"
	"time"
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
	spanCtx, cancel := context.WithTimeout(spanCtx, 5*time.Second)
	defer cancel()
	span, spanCtx := opentracing.StartSpanFromContext(spanCtx, "ModelHandler.GetAllModels")
	defer span.Finish()

	models, err := api.ModelService.GetAll(spanCtx)
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
