package model

import (
	"Brands/internal/api/handler/utils"
	modelrepo "Brands/internal/repository/model"
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"net/http"
)

// DeleteModel godoc
// @Summary Мягкое удаление модели
// @Description Выполняет мягкое удаление модели по её ID
// @Tags models
// @Accept json
// @Produce json
// @Param id path string true "ID модели (UUIDv7)"
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

	span, spanCtx := opentracing.StartSpanFromContext(spanCtx, "ModelHandler.DeleteModel")
	defer span.Finish()

	// Извлечение и парсинг UUID из пути запроса
	id, err := utils.ExtractUUIDFromPath(ctx, "id")
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(
			log.String("event", "invalid_id"),
			log.Error(err),
		)
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString("Invalid ID format")
		return
	}
	err = api.ModelService.SoftDelete(spanCtx, id)
	if err != nil {
		span.SetTag("error", true)
		if errors.Is(err, modelrepo.ErrModelNotFound) {
			span.LogFields(
				log.String("event", "model_not_found"),
				log.String("model.id", id.String()),
			)
			ctx.Response.SetStatusCode(http.StatusNotFound)
			ctx.Response.SetBodyString(fmt.Sprintf("Model not found with ID: %s", id))
			return
		}
		span.LogFields(
			log.String("event", "delete_model_error"),
			log.Error(err),
			log.String("model.id", id.String()),
		)
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
// @Param id path string true "ID модели (UUIDv7)"
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

	span, spanCtx := opentracing.StartSpanFromContext(spanCtx, "ModelHandler.RestoreModel")
	defer span.Finish()

	// Извлечение и парсинг UUID из пути запроса
	id, err := utils.ExtractUUIDFromPath(ctx, "id")
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(
			log.String("event", "invalid_id"),
			log.Error(err),
		)
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString("Invalid ID format")
		return
	}

	err = api.ModelService.Restore(spanCtx, id)
	if err != nil {
		span.SetTag("error", true)
		if errors.Is(err, modelrepo.ErrModelNotFound) {
			span.LogFields(
				log.String("event", "model_not_found"),
				log.String("model.id", id.String()),
			)
			ctx.Response.SetStatusCode(http.StatusNotFound)
			ctx.Response.SetBodyString(fmt.Sprintf("Model not found with ID: %s", id))
			return
		}
		span.LogFields(
			log.String("event", "restore_model_error"),
			log.Error(err),
			log.String("model.id", id.String()),
		)
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to restore model: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBodyString("Model restored successfully")
}
