package brand

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
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
// @Param id path int true "ID бренда"
// @Success 200 {string} string "Brand soft-deleted successfully"
// @Failure 400 {string} string "Invalid ID format"
// @Failure 500 {string} string "Failed to delete brand"
// @Router /brands/delete/{id} [delete]
func (api *Handler) DeleteBrand(ctx *fasthttp.RequestCtx) {
	var spanCtx context.Context
	spanCtx, ok := ctx.UserValue("traceContext").(context.Context)
	if !ok {
		spanCtx = ctx
	}
	spanCtx, cancel := context.WithTimeout(spanCtx, 5*time.Second)
	defer cancel()
	span, spanCtx := opentracing.StartSpanFromContext(spanCtx, "Handler.DeleteBrand")
	defer span.Finish()

	idStr := ctx.UserValue("id").(string)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(
			log.String("event", "decode_error"),
			log.String("error", err.Error()),
		)
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString("Invalid ID format")
		return
	}

	err = api.BrandService.SoftDelete(spanCtx, id)
	if err != nil {
		span.SetTag("error", true)
		if errors.Is(err, brandrepo.ErrBrandNotFound) {
			span.LogFields(
				log.String("event", "brand_not_found"),
				log.Int64("brand.id", id),
			)
			ctx.Response.SetStatusCode(http.StatusNotFound)
			ctx.Response.SetBodyString(fmt.Sprintf("Brand not found with ID: %d", id))
			return
		}
		span.LogFields(
			log.String("event", "delete_brand_error"),
			log.String("error", err.Error()),
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
// @Param id path int true "ID бренда"
// @Success 200 {string} string "Brand restored successfully"
// @Failure 400 {string} string "Invalid ID format"
// @Failure 500 {string} string "Failed to restore brand"
// @Router /brands/restore/{id} [post]
func (api *Handler) RestoreBrand(ctx *fasthttp.RequestCtx) {
	spanCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	span, spanCtx := opentracing.StartSpanFromContext(spanCtx, "Handler.RestoreBrand")
	defer span.Finish()

	idStr := ctx.UserValue("id").(string)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(
			log.String("event", "decode_error"),
			log.String("error", err.Error()),
		)
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString("Invalid ID format")
		return
	}

	err = api.BrandService.Restore(spanCtx, id)
	if err != nil {
		span.SetTag("error", true)
		if errors.Is(err, brandrepo.ErrBrandNotFound) {
			span.LogFields(
				log.String("event", "brand_not_found"),
				log.Int64("brand.id", id),
			)
			ctx.Response.SetStatusCode(http.StatusNotFound)
			ctx.Response.SetBodyString(fmt.Sprintf("Brand not found with ID: %d", id))
			return
		}
		span.LogFields(
			log.String("event", "restore_brand_error"),
			log.String("error", err.Error()),
		)
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to restore brand: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBodyString("Brand restored successfully")
}
