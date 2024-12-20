package model

import (
	"Brands/internal/dto"
	"Brands/pkg/zerohook"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/valyala/fasthttp"
	"net/http"
)

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
	var spanCtx context.Context
	spanCtx, ok := ctx.UserValue("traceContext").(context.Context)
	if !ok {
		spanCtx = ctx
	}

	span, spanCtx := opentracing.StartSpanFromContext(spanCtx, "ModelHandler.CreateModel")
	defer span.Finish()

	decoder := json.NewDecoder(bytes.NewReader(ctx.PostBody()))
	var model dto.Model
	err := decoder.Decode(&model)
	if err != nil {
		span.SetTag("error", true)
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}
	model.ID, err = uuid.NewV7()
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(
			log.String("event", "new_uuid_error"),
			log.Error(err),
			log.Object("model", model),
		)
		zerohook.Logger.Error().Err(err).Send()
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		return
	}

	err = api.ModelService.Create(spanCtx, &model)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(
			log.String("event", "create_model_error"),
			log.Object("model", model),
			log.Error(err),
		)
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to create model: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBodyString(fmt.Sprintf("Model created successfully with ID: %s", model.ID))
}
