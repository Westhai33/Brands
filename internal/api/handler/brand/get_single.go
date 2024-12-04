package brand

import (
	"Brands/internal/api/handler"
	brandrepo "Brands/internal/repository/brand"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/valyala/fasthttp"
	"net/http"
)

// GetBrandByID godoc
// @Summary Получение бренда по ID
// @Description Получение информации о бренде по его уникальному ID
// @Tags brand
// @Accept json
// @Produce json
// @Param id path string true "ID бренда"
// @Success 200 {object} dto.Brand "Бренд найден"
// @Failure 400 {string} string "Invalid ID format"
// @Failure 404 {string} string "Brand not found"
// @Router /brands/{id} [get]
func (api *BrandHandler) GetBrandByID(ctx *fasthttp.RequestCtx) {
	var spanCtx context.Context
	spanCtx, ok := ctx.UserValue("traceContext").(context.Context)
	if !ok {
		spanCtx = ctx
	}
	span, spanCtx := opentracing.StartSpanFromContext(spanCtx, "BrandHandler.GetBrandByID")
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

	brand, err := api.BrandService.GetByID(spanCtx, id)
	if err != nil {
		span.SetTag("error", true)
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
			log.String("event", "get_brand_error"),
			log.Error(err),
			log.String("brand.id", id.String()),
		)
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to get brand: %v", err))
		return
	}

	data, err := json.Marshal(brand)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(
			log.String("event", "json_marshal_error"),
			log.Error(err),
			log.Object("brand", brand),
		)
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to marshal brand data: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBody(data)
}
