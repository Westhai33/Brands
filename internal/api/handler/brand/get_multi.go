package brand

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/valyala/fasthttp"
	"net/http"
	"time"
)

// GetAllBrands godoc
// @Summary Получение всех брендов
// @Description Возвращает все бренды с возможностью фильтрации и сортировки
// @Tags brand
// @Accept json
// @Produce json
// @Success 200 {array} dto.Brand "Список брендов"
// @Failure 500 {string} string "Failed to fetch brand"
// @Router /brands/all [get]
func (api *BrandHandler) GetAllBrands(ctx *fasthttp.RequestCtx) {
	var spanCtx context.Context
	spanCtx, ok := ctx.UserValue("traceContext").(context.Context)
	if !ok {
		spanCtx = ctx
	}
	spanCtx, cancel := context.WithTimeout(spanCtx, 5*time.Second)
	defer cancel()
	span, spanCtx := opentracing.StartSpanFromContext(spanCtx, "BrandHandler.GetAllBrands")

	// Вызов метода GetAll из сервиса с фильтрами и сортировкой
	brands, err := api.BrandService.GetAll(spanCtx)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(
			log.String("event", "failed_to_fetch_brands"),
			log.String("error", err.Error()),
		)
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to fetch brand: %v", err))
		return
	}

	// Преобразуем список брендов в JSON
	data, err := json.Marshal(brands)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(
			log.String("event", "failed_to_marshal_brands"),
			log.String("error", err.Error()),
		)
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to marshal brand data: %v", err))
		return
	}

	// Отправляем список брендов в ответ
	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBody(data)
}
