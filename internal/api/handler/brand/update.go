package brand

import (
	"Brands/internal/dto"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
	"time"
)

// UpdateBrand godoc
// @Summary Обновление бренда по ID
// @Description Обновление данных бренда по его ID
// @Tags brand
// @Accept json
// @Produce json
// @Param id path int true "ID бренда"
// @Param brand body dto.Brand true "Обновлённые данные бренда"
// @Success 200 {string} string "Brand updated successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 404 {string} string "Brand not found"
// @Failure 500 {string} string "Failed to update brand"
// @Router /brands/update/{id} [put]
func (api *BrandHandler) UpdateBrand(ctx *fasthttp.RequestCtx) {
	var spanCtx context.Context
	spanCtx, ok := ctx.UserValue("traceContext").(context.Context)
	if !ok {
		spanCtx = ctx
	}
	spanCtx, cancel := context.WithTimeout(spanCtx, 5*time.Second)
	defer cancel()
	span, spanCtx := opentracing.StartSpanFromContext(spanCtx, "BrandHandler.UpdateBrand")
	defer span.Finish()

	idStr := ctx.UserValue("id").(string)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(
			log.String("event", "invalid_id_format"),
			log.String("error", err.Error()),
		)
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString("Invalid ID format")
		return
	}

	decoder := json.NewDecoder(bytes.NewReader(ctx.PostBody()))
	var brand dto.Brand
	err = decoder.Decode(&brand)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(
			log.String("event", "decode_error"),
			log.String("error", err.Error()),
		)
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}

	brand.ID = id

	err = api.BrandService.Update(spanCtx, &brand)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(
			log.String("event", "update_brand_error"),
			log.String("error", err.Error()),
		)
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to update brand: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBodyString("Brand updated successfully")
}
