package model

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
	"time"
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
func (api *ModelHandler) DeleteModel(ctx *fasthttp.RequestCtx) {
	var spanCtx context.Context
	spanCtx, ok := ctx.UserValue("traceContext").(context.Context)
	if !ok {
		spanCtx = ctx
	}
	spanCtx, cancel := context.WithTimeout(spanCtx, 5*time.Second)
	defer cancel()
	span, spanCtx := opentracing.StartSpanFromContext(spanCtx, "ModelHandler.DeleteModel")
	defer span.Finish()

	idStr := ctx.UserValue("id").(string)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString("Invalid ID format")
		return
	}

	err = api.ModelService.SoftDelete(spanCtx, id)
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
func (api *ModelHandler) RestoreModel(ctx *fasthttp.RequestCtx) {
	var spanCtx context.Context
	spanCtx, ok := ctx.UserValue("traceContext").(context.Context)
	if !ok {
		spanCtx = ctx
	}
	spanCtx, cancel := context.WithTimeout(spanCtx, 5*time.Second)
	defer cancel()
	span, spanCtx := opentracing.StartSpanFromContext(spanCtx, "ModelHandler.RestoreModel")
	defer span.Finish()

	idStr := ctx.UserValue("id").(string)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString("Invalid ID format")
		return
	}

	err = api.ModelService.Restore(spanCtx, id)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to restore model: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBodyString("Model restored successfully")
}
