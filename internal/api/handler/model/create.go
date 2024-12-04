package model

import (
	"Brands/internal/dto"
	"bytes"
	"encoding/json"
	"fmt"
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
func (api *Handler) CreateModel(ctx *fasthttp.RequestCtx) {
	decoder := json.NewDecoder(bytes.NewReader(ctx.PostBody()))
	var model dto.Model
	err := decoder.Decode(&model)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}

	modelID, err := api.ModelService.Create(ctx, &model)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to create model: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBodyString(fmt.Sprintf("Model created successfully with ID: %d", modelID))
}
