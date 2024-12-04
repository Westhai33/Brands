package brand

import (
	"Brands/internal/dto"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"net/http"
)

// UpdateBrand godoc
// @Summary Обновление бренда по ID
// @Description Обновление данных бренда по его ID
// @Tags brand
// @Accept json
// @Produce json
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
	span, spanCtx := opentracing.StartSpanFromContext(spanCtx, "BrandHandler.UpdateBrand")
	defer span.Finish()

	decoder := json.NewDecoder(bytes.NewReader(ctx.PostBody()))
	var brand dto.Brand
	err := decoder.Decode(&brand)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(
			log.String("event", "decode_error"),
			log.Error(err),
		)
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}

	if brand.Name == "" {
		err = errors.New("Name is required")
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString(err.Error())
		return
	}

	err = api.BrandService.Update(spanCtx, &brand)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(
			log.String("event", "update_brand_error"),
			log.Error(err),
			log.Object("brand", brand),
		)
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to update brand: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBodyString("Brand updated successfully")
}
