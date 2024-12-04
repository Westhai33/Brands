package brand

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// GetAllBrands godoc
// @Summary Получение всех брендов
// @Description Возвращает все бренды с возможностью фильтрации и сортировки
// @Tags brand
// @Accept json
// @Produce json
// @Param filter query string false "Фильтрация по полю"
// @Param sort query string false "Сортировка по полю"
// @Success 200 {array} dto.Brand "Список брендов"
// @Failure 500 {string} string "Failed to fetch brand"
// @Router /brands/all [get]
func (api *Handler) GetAllBrands(ctx *fasthttp.RequestCtx) {
	var spanCtx context.Context
	spanCtx, ok := ctx.UserValue("traceContext").(context.Context)
	if !ok {
		spanCtx = ctx
	}
	spanCtx, cancel := context.WithTimeout(spanCtx, 5*time.Second)
	defer cancel()
	span, spanCtx := opentracing.StartSpanFromContext(spanCtx, "Handler.GetAllBrands")
	defer span.Finish()

	// Извлечение фильтров из query параметров
	filter := make(map[string]interface{})

	// Добавляем фильтры на основе входных параметров
	if name := string(ctx.QueryArgs().Peek("name")); name != "" {
		filter["name"] = name
	}
	if originCountry := string(ctx.QueryArgs().Peek("origin_country")); originCountry != "" {
		filter["origin_country"] = originCountry
	}
	if popularity := ctx.QueryArgs().Peek("popularity"); len(popularity) > 0 {
		if popVal, err := strconv.Atoi(string(popularity)); err == nil {
			filter["popularity"] = popVal
		} else {
			span.SetTag("error", true)
			span.LogFields(
				log.String("event", "invalid_popularity_value"),
				log.String("error", err.Error()),
			)
			ctx.Response.SetStatusCode(http.StatusBadRequest)
			ctx.Response.SetBodyString("Invalid popularity value")
			return
		}
	}
	if isPremium := ctx.QueryArgs().Peek("is_premium"); len(isPremium) > 0 {
		if isPremiumVal, err := strconv.ParseBool(string(isPremium)); err == nil {
			filter["is_premium"] = isPremiumVal
		} else {
			span.SetTag("error", true)
			span.LogFields(
				log.String("event", "invalid_is_premium_value"),
				log.String("error", err.Error()),
			)
			ctx.Response.SetStatusCode(http.StatusBadRequest)
			ctx.Response.SetBodyString("Invalid is_premium value")
			return
		}
	}

	// Обработка сортировки
	sort := string(ctx.QueryArgs().Peek("sort"))
	if sort != "" {
		validSortFields := map[string]bool{
			"name":           true,
			"popularity":     true,
			"founded_year":   true,
			"origin_country": true,
			"created_at":     true,
		}

		// Проверяем, является ли поле сортировки допустимым
		sortField := strings.TrimPrefix(sort, "-")
		if !validSortFields[sortField] {
			span.SetTag("error", true)
			errMessage := fmt.Sprintf("Invalid sort field: %s", sortField)
			span.LogFields(
				log.String("event", "invalid_sort_field"),
				log.String("error", errMessage),
			)
			ctx.Response.SetStatusCode(http.StatusBadRequest)
			ctx.Response.SetBodyString(errMessage)
			return
		}
	}

	// Вызов метода GetAll из сервиса с фильтрами и сортировкой
	brands, err := api.BrandService.GetAll(spanCtx, filter, sort)
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
