package model

import (
	"Brands/internal/dto"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
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

	err = api.ModelService.Update(ctx, &model)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to update model: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBodyString("Model updated successfully")
}
