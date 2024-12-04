package model

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
)

// DeleteModel godoc
// @Summary Мягкое удаление модели
// @Description Выполняет мягкое удаление модели по её ID
// @Tags models
// @Accept json
// @Produce json
// @Param id path int true "ID модели"
// @Success 200 {string} string "Model soft-deleted successfully"
// @Failure 400 {string} string "Invalid ID format"
// @Failure 500 {string} string "Failed to delete model"
// @Router /models/delete/{id} [delete]
func (api *Handler) DeleteModel(ctx *fasthttp.RequestCtx) {
	idStr := ctx.UserValue("id").(string)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString("Invalid ID format")
		return
	}

	err = api.ModelService.SoftDelete(ctx, id)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to delete model: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBodyString("Model soft-deleted successfully")
}

// RestoreModel godoc
// @Summary Восстановление мягко удалённой модели
// @Description Восстановление модели, которая была удалена ранее
// @Tags models
// @Accept json
// @Produce json
// @Param id path int true "ID модели"
// @Success 200 {string} string "Model restored successfully"
// @Failure 400 {string} string "Invalid ID format"
// @Failure 500 {string} string "Failed to restore model"
// @Router /models/restore/{id} [post]
func (api *Handler) RestoreModel(ctx *fasthttp.RequestCtx) {
	idStr := ctx.UserValue("id").(string)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString("Invalid ID format")
		return
	}

	err = api.ModelService.Restore(ctx, id)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to restore model: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBodyString("Model restored successfully")
}
