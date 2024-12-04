package model

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
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
	idStr := ctx.UserValue("id").(string)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString("Invalid ID format")
		return
	}

	model, err := api.ModelService.GetByID(ctx, id)
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
