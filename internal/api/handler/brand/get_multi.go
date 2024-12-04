package brand

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
	"strings"
)

// GetAllBrands godoc
// @Summary Получение всех брендов
// @Description Возвращает все бренды с возможностью фильтрации и сортировки
// @Tags brand
// @Accept json
// @Produce json
// @Param name query string false "Фильтр по имени бренда"
// @Param origin_country query string false "Фильтр по стране происхождения"
// @Param popularity query integer false "Фильтр по популярности (целое число)"
// @Param is_premium query boolean false "Фильтр по признаку премиум-бренда"
// @Param is_upcoming query boolean false "Фильтр по признаку предстоящего бренда"
// @Param founded_year query integer false "Фильтр по году основания"
// @Param sort query string false "Поле сортировки (например, 'name', '-popularity')"
// @Success 200 {array} dto.Brand "Список брендов"
// @Failure 500 {string} string "Failed to fetch brands"
// @Router /brands/all [get]
func (api *BrandHandler) GetAllBrands(ctx *fasthttp.RequestCtx) {
	// Извлечение фильтров из query параметров
	filter := make(map[string]interface{})

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
			ctx.Response.SetStatusCode(http.StatusBadRequest)
			ctx.Response.SetBodyString("Invalid popularity value")
			return
		}
	}
	if isPremium := ctx.QueryArgs().Peek("is_premium"); len(isPremium) > 0 {
		if isPremiumVal, err := strconv.ParseBool(string(isPremium)); err == nil {
			filter["is_premium"] = isPremiumVal
		} else {
			ctx.Response.SetStatusCode(http.StatusBadRequest)
			ctx.Response.SetBodyString("Invalid is_premium value")
			return
		}
	}

	sort := string(ctx.QueryArgs().Peek("sort"))
	if sort != "" {
		validSortFields := map[string]bool{
			"name":           true,
			"popularity":     true,
			"founded_year":   true,
			"origin_country": true,
			"created_at":     true,
		}

		sortField := strings.TrimPrefix(sort, "-")
		if !validSortFields[sortField] {
			ctx.Response.SetStatusCode(http.StatusBadRequest)
			ctx.Response.SetBodyString(fmt.Sprintf("Invalid sort field: %s", sortField))
			return
		}
	}

	brands, err := api.BrandService.GetAll(ctx, filter, sort)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to fetch brands: %v", err))
		return
	}

	data, err := json.Marshal(brands)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Failed to marshal brands data: %v", err))
		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBody(data)
}
