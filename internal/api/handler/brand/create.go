package brand

import (
	"Brands/internal/dto"
	"Brands/pkg/zerohook"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"net/http"
)

type CreateBrandResponse struct {
	ID uuid.UUID `json:"id"`
}

// CreateBrand godoc
// @Summary Создание нового бренда
// @Description Эндпоинт для создания нового бренда
// @Tags brand
// @Accept json
// @Produce json
// @Param brand body dto.Brand true "Данные нового бренда"
// @Success 200 {object} CreateBrandResponse "Brand created successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Failed to create brand"
// @Router /brands/create [post]
func (api *BrandHandler) CreateBrand(ctx *fasthttp.RequestCtx) {
	var spanCtx context.Context
	spanCtx, ok := ctx.UserValue("traceContext").(context.Context)
	if !ok {
		spanCtx = ctx
	}
	span, spanCtx := opentracing.StartSpanFromContext(spanCtx, "BrandHandler.CreateBrand")
	defer span.Finish()

	decoder := json.NewDecoder(bytes.NewReader(ctx.PostBody()))
	var brand dto.Brand
	err := decoder.Decode(&brand)
	if err != nil {
		span.LogFields(
			log.String("event", "decode_error"),
			log.String("error", err.Error()),
		)
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}
	brand.ID, err = uuid.NewV7()
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(
			log.String("event", "new_uuid_error"),
			log.String("error", err.Error()),
		)
		zerohook.Logger.Error().Err(err).Send()
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		return
	}
	if brand.Name == "" {
		err = errors.New("Name is required")
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBodyString(err.Error())
		return
	}

	err = api.BrandService.Create(spanCtx, &brand)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(
			log.String("event", "create_brand_error"),
			log.String("error", err.Error()),
		)
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to create brand: %v", err))
		return
	}
	data, err := json.Marshal(
		CreateBrandResponse{
			ID: brand.ID,
		},
	)
	if err != nil {

		span.SetTag("error", true)
		span.LogFields(
			log.String("event", "marshal_response_error"),
			log.String("error", err.Error()),
		)
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString("Failed to marshal response")
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBodyString(fmt.Sprintf("Brand created successfully with ID: %d", brand.ID))
	ctx.Response.SetBody(data)
}
