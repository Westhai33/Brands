package model

import (
	"Brands/internal/dto"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
	"time"
)

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
func (api *Handler) UpdateModel(ctx *fasthttp.RequestCtx) {
	var spanCtx context.Context
	spanCtx, ok := ctx.UserValue("traceContext").(context.Context)
	if !ok {
		spanCtx = ctx
	}
	spanCtx, cancel := context.WithTimeout(spanCtx, 5*time.Second)
	defer cancel()
	span, spanCtx := opentracing.StartSpanFromContext(spanCtx, "Handler.UpdateModel")
	defer span.Finish()

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

	err = api.ModelService.Update(spanCtx, &model)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to update model: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBodyString("Model updated successfully")
}
