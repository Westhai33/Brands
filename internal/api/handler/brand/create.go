package brand

import (
	"Brands/internal/dto"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/valyala/fasthttp"
	"net/http"
	"time"
)

// CreateBrand godoc
// @Summary Создание нового бренда
// @Description Эндпоинт для создания нового бренда
// @Tags brand
// @Accept json
// @Produce json
// @Param brand body dto.Brand true "Данные нового бренда"
// @Success 200 {string} string "Brand created successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Failed to create brand"
// @Router /brands/create [post]
func (api *Handler) CreateBrand(ctx *fasthttp.RequestCtx) {
	var spanCtx context.Context
	spanCtx, ok := ctx.UserValue("traceContext").(context.Context)
	if !ok {
		spanCtx = ctx
	}
	spanCtx, cancel := context.WithTimeout(spanCtx, 5*time.Second)
	defer cancel()
	span, spanCtx := opentracing.StartSpanFromContext(spanCtx, "Handler.CreateBrand")
	defer span.Finish()

	decoder := json.NewDecoder(bytes.NewReader(ctx.PostBody()))
	var brand dto.Brand
	err := decoder.Decode(&brand)
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

	brandID, err := api.BrandService.Create(spanCtx, &brand)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(
			log.String("event", "create_brand_error"),
			log.String("error", err.Error()),
		)
		if errors.Is(spanCtx.Err(), context.DeadlineExceeded) {
			ctx.Response.SetStatusCode(http.StatusRequestTimeout)
			ctx.Response.SetBodyString("Request timeout exceeded")
		} else {
			ctx.Response.SetStatusCode(http.StatusInternalServerError)
			ctx.Response.SetBodyString(fmt.Sprintf("Failed to create brand: %v", err))
		}
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBodyString(fmt.Sprintf("Brand created successfully with ID: %d", brandID))
}
