package brand

import (
	"Brands/internal/api/handler"
	"context"
	"fmt"
	"net/http"
	"time"

	brandrepo "Brands/internal/repository/brand"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

// DeleteBrand godoc
// @Summary Мягкое удаление бренда
// @Description Выполняет мягкое удаление бренда по его ID
// @Tags brand
// @Accept json
// @Produce json
// @Param id path string true "ID бренда"
// @Success 200 {string} string "Brand soft-deleted successfully"
// @Failure 400 {string} string "Invalid ID format"
// @Failure 500 {string} string "Failed to delete brand"
// @Router /brands/delete/{id} [delete]
func (api *BrandHandler) DeleteBrand(ctx *fasthttp.RequestCtx) {
	var spanCtx context.Context
	spanCtx, ok := ctx.UserValue("traceContext").(context.Context)
	if !ok {
		spanCtx = ctx
	}

	span, spanCtx := opentracing.StartSpanFromContext(spanCtx, "BrandHandler.DeleteBrand")
	defer span.Finish()

	// Извлечение и парсинг UUID из пути запроса
	id, err := handler.ExtractUUIDFromPath(ctx, "id")
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

	err = api.BrandService.SoftDelete(spanCtx, id)
	if err != nil {

		if errors.Is(err, brandrepo.ErrBrandNotFound) {
			span.LogFields(
				log.String("event", "brand_not_found"),
				log.String("brand.id", id.String()),
			)
			ctx.Response.SetStatusCode(http.StatusNotFound)
			ctx.Response.SetBodyString(fmt.Sprintf("Brand not found with ID: %d", id))
			return
		}
		span.LogFields(
			log.String("event", "delete_brand_error"),
			log.Error(err),
		)
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to delete brand: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBodyString("Brand soft-deleted successfully")
}

// RestoreBrand godoc
// @Summary Восстановление мягко удалённого бренда
// @Description Восстановление бренда, который был удалён ранее (мягкое удаление)
// @Tags brand
// @Accept json
// @Produce json
// @Param id path string true "ID бренда"
// @Success 200 {string} string "Brand restored successfully"
// @Failure 400 {string} string "Invalid ID format"
// @Failure 500 {string} string "Failed to restore brand"
// @Router /brands/restore/{id} [post]
func (api *BrandHandler) RestoreBrand(ctx *fasthttp.RequestCtx) {
	spanCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	span, spanCtx := opentracing.StartSpanFromContext(spanCtx, "BrandHandler.RestoreBrand")
	defer span.Finish()

	// Извлечение и парсинг UUID из пути запроса
	id, err := handler.ExtractUUIDFromPath(ctx, "id")
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

	err = api.BrandService.Restore(spanCtx, id)
	if err != nil {
		if errors.Is(err, brandrepo.ErrBrandNotFound) {
			span.LogFields(
				log.String("event", "brand_not_found"),
				log.String("brand.id", id.String()),
			)
			ctx.Response.SetStatusCode(http.StatusNotFound)
			ctx.Response.SetBodyString(fmt.Sprintf("Brand not found with ID: %d", id))
			return
		}
		span.LogFields(
			log.String("event", "restore_brand_error"),
			log.Error(err),
			log.String("brand.id", id.String()),
		)
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to restore brand: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBodyString("Brand restored successfully")
}
