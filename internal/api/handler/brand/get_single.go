package brand

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	brandrepo "Brands/internal/repository/brand"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/valyala/fasthttp"
)

// GetBrandByID godoc
// @Summary Получение бренда по ID
// @Description Получение информации о бренде по его уникальному ID
// @Tags brand
// @Accept json
// @Produce json
// @Param id path int true "ID бренда"
// @Success 200 {object} dto.Brand "Бренд найден"
// @Failure 400 {string} string "Invalid ID format"
// @Failure 404 {string} string "Brand not found"
// @Router /brands/{id} [get]
func (api *Handler) GetBrandByID(ctx *fasthttp.RequestCtx) {
	spanCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	span, timeoutCtx := opentracing.StartSpanFromContext(spanCtx, "Handler.GetBrandByID")
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

	brand, err := api.BrandService.GetByID(timeoutCtx, id)
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
			log.String("event", "get_brand_error"),
			log.String("error", err.Error()),
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
			log.String("error", err.Error()),
		)
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to marshal brand data: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBody(data)
}
