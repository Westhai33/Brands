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

type CreateModelResponse struct {
	ID uuid.UUID `json:"id"`
}

// CreateModel godoc
// @Summary Создание новой модели
// @Description Эндпоинт для создания новой модели
// @Tags models
// @Accept json
// @Produce json
// @Param model body dto.Model true "Данные новой модели"
// @Success 200 {object} CreateModelResponse "Model created successfully"
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
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}
	model.ID, err = uuid.NewV7()
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(
			log.String("event", "new_uuid_error"),
			log.String("error", err.Error()),
		)
		zerohook.Logger.Error().Err(err).Send()
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		return
	}

	err = api.ModelService.Create(spanCtx, &model)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to create model: %v", err))
		return
	}

	data, err := json.Marshal(
		CreateModelResponse{
			ID: model.ID,
		},
	)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(
			log.String("event", "marshal_response_error"),
			log.String("error", err.Error()),
		)
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString("Failed to marshal response")
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBodyString(fmt.Sprintf("Model created successfully with ID: %d", model.ID))
	ctx.Response.SetBody(data)
}
