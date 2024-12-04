package model

import (
	"Brands/internal/dto"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"net/http"
)

// UpdateModel godoc
// @Summary Обновление модели по ID
// @Description Обновление данных модели по её ID
// @Tags models
// @Accept json
// @Produce json
// @Param model body dto.Model true "Обновлённые данные модели"
// @Success 200 {string} string "Model updated successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 404 {string} string "Model not found"
// @Failure 500 {string} string "Failed to update model"
// @Router /models/update/{id} [put]
func (api *ModelHandler) UpdateModel(ctx *fasthttp.RequestCtx) {
	var spanCtx context.Context
	spanCtx, ok := ctx.UserValue("traceContext").(context.Context)
	if !ok {
		spanCtx = ctx
	}
	span, spanCtx := opentracing.StartSpanFromContext(spanCtx, "ModelHandler.UpdateModel")
	defer span.Finish()

	decoder := json.NewDecoder(bytes.NewReader(ctx.PostBody()))
	var model dto.Model
	err := decoder.Decode(&model)
	if err != nil {
		span.LogFields(
			log.String("event", "decode_error"),
			log.Error(err),
		)
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}
	if model.Name == "" {
		err = errors.New("Name is required")
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString(err.Error())
		return
	}

	err = api.ModelService.Update(spanCtx, &model)
	if err != nil {
		span.LogFields(
			log.String("event", "update_model_error"),
			log.Error(err),
			log.Object("model", model),
		)
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to update model: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBodyString("Model updated successfully")
}
