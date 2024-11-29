package controller

import (
	"Brands/internal/dto"
	"Brands/internal/service"
	"fmt"
	"github.com/valyala/fasthttp"
	"strconv"
)

// BrandController контролирует обработку запросов для брендов
type BrandController struct {
	BrandService *service.BrandService
}

// NewBrandController создаёт новый контроллер
func NewBrandController(brandService *service.BrandService) *BrandController {
	return &BrandController{BrandService: brandService}
}

// CreateBrand создает новый бренд из данных в запросе
func (c *BrandController) CreateBrand(ctx *fasthttp.RequestCtx) {
	// Извлекаем данные из тела запроса
	name := string(ctx.FormValue("name"))
	link := string(ctx.FormValue("link"))
	description := string(ctx.FormValue("description"))
	logoURL := string(ctx.FormValue("logo_url"))
	coverImageURL := string(ctx.FormValue("cover_image_url"))
	foundedYearStr := string(ctx.FormValue("founded_year"))
	foundedYear, err := strconv.Atoi(foundedYearStr)
	if err != nil {
		ctx.Error(fmt.Sprintf("не удалось преобразовать founded_year в число: %v", err), fasthttp.StatusBadRequest)
		return
	}

	originCountry := string(ctx.FormValue("origin_country"))
	popularityStr := string(ctx.FormValue("popularity"))
	popularity, err := strconv.Atoi(popularityStr)
	if err != nil {
		ctx.Error(fmt.Sprintf("не удалось преобразовать popularity в число: %v", err), fasthttp.StatusBadRequest)
		return
	}

	isPremium := string(ctx.FormValue("is_premium")) == "true"
	isUpcoming := string(ctx.FormValue("is_upcoming")) == "true"

	if name == "" {
		ctx.Error("название обязательно", fasthttp.StatusBadRequest)
		return
	}

	brand := &dto.Brand{
		Name:          name,
		Link:          link,
		Description:   description,
		LogoURL:       logoURL,
		CoverImageURL: coverImageURL,
		FoundedYear:   foundedYear,
		OriginCountry: originCountry,
		Popularity:    popularity,
		IsPremium:     isPremium,
		IsUpcoming:    isUpcoming,
	}

	brandID, err := c.BrandService.Create(ctx, brand)
	if err != nil {
		ctx.Error(fmt.Sprintf("ошибка при создании бренда: %v", err), fasthttp.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ
	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.SetBody([]byte(fmt.Sprintf("Бренд с ID %d успешно создан", brandID)))
}

// GetBrandByID получает бренд по ID
func (c *BrandController) GetBrandByID(ctx *fasthttp.RequestCtx) {
	// Преобразуем id из строки в int64
	brandIDStr := ctx.UserValue("id").(string) // предполагаем, что id приходит как строка
	brandID, err := strconv.ParseInt(brandIDStr, 10, 64)
	if err != nil {
		ctx.Error(fmt.Sprintf("не удалось преобразовать id в число: %v", err), fasthttp.StatusBadRequest)
		return
	}

	brand, err := c.BrandService.GetByID(ctx, brandID)
	if err != nil {
		ctx.Error(fmt.Sprintf("бренд с ID %d не найден: %v", brandID, err), fasthttp.StatusNotFound)
		return
	}

	// Преобразуем бренд в JSON
	brandJSON, err := brand.ToJSON()
	if err != nil {
		ctx.Error(fmt.Sprintf("ошибка преобразования бренда в JSON: %v", err), fasthttp.StatusInternalServerError)
		return
	}

	// Преобразуем строку в []byte и отправляем в теле ответа
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(brandJSON)) // Конвертируем строку в []byte
}

// UpdateBrand обновляет данные бренда с проверкой существования
func (c *BrandController) UpdateBrand(ctx *fasthttp.RequestCtx) {
	// Преобразуем id из строки в int64
	brandIDStr := ctx.UserValue("id").(string) // предполагаем, что id приходит как строка
	brandID, err := strconv.ParseInt(brandIDStr, 10, 64)
	if err != nil {
		ctx.Error(fmt.Sprintf("не удалось преобразовать id в число: %v", err), fasthttp.StatusBadRequest)
		return
	}

	// Извлекаем данные из тела запроса
	name := string(ctx.FormValue("name"))
	link := string(ctx.FormValue("link"))
	description := string(ctx.FormValue("description"))
	logoURL := string(ctx.FormValue("logo_url"))
	coverImageURL := string(ctx.FormValue("cover_image_url"))
	foundedYearStr := string(ctx.FormValue("founded_year"))
	foundedYear, err := strconv.Atoi(foundedYearStr)
	if err != nil {
		ctx.Error(fmt.Sprintf("не удалось преобразовать founded_year в число: %v", err), fasthttp.StatusBadRequest)
		return
	}

	originCountry := string(ctx.FormValue("origin_country"))
	popularityStr := string(ctx.FormValue("popularity"))
	popularity, err := strconv.Atoi(popularityStr)
	if err != nil {
		ctx.Error(fmt.Sprintf("не удалось преобразовать popularity в число: %v", err), fasthttp.StatusBadRequest)
		return
	}

	isPremium := string(ctx.FormValue("is_premium")) == "true"
	isUpcoming := string(ctx.FormValue("is_upcoming")) == "true"

	brand, err := c.BrandService.GetByID(ctx, brandID)
	if err != nil {
		ctx.Error(fmt.Sprintf("бренд с ID %d не найден: %v", brandID, err), fasthttp.StatusNotFound)
		return
	}

	brand.Name = name
	brand.Link = link
	brand.Description = description
	brand.LogoURL = logoURL
	brand.CoverImageURL = coverImageURL
	brand.FoundedYear = foundedYear
	brand.OriginCountry = originCountry
	brand.Popularity = popularity
	brand.IsPremium = isPremium
	brand.IsUpcoming = isUpcoming

	err = c.BrandService.Update(ctx, brand)
	if err != nil {
		ctx.Error(fmt.Sprintf("ошибка обновления бренда с ID %d: %v", brandID, err), fasthttp.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(fmt.Sprintf("Бренд с ID %d успешно обновлен", brandID)))
}

// SoftDeleteBrand мягко удаляет бренд с проверкой существования
func (c *BrandController) SoftDeleteBrand(ctx *fasthttp.RequestCtx) {
	brandIDStr := ctx.UserValue("id").(string)
	brandID, err := strconv.ParseInt(brandIDStr, 10, 64)
	if err != nil {
		ctx.Error(fmt.Sprintf("не удалось преобразовать id в число: %v", err), fasthttp.StatusBadRequest)
		return
	}

	err = c.BrandService.SoftDelete(ctx, brandID)
	if err != nil {
		ctx.Error(fmt.Sprintf("ошибка мягкого удаления бренда с ID %d: %v", brandID, err), fasthttp.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(fmt.Sprintf("Бренд с ID %d успешно удалён", brandID)))
}

// RestoreBrand восстанавливает мягко удалённый бренд
func (c *BrandController) RestoreBrand(ctx *fasthttp.RequestCtx) {
	brandIDStr := ctx.UserValue("id").(string)
	brandID, err := strconv.ParseInt(brandIDStr, 10, 64)
	if err != nil {
		ctx.Error(fmt.Sprintf("не удалось преобразовать id в число: %v", err), fasthttp.StatusBadRequest)
		return
	}

	err = c.BrandService.Restore(ctx, brandID)
	if err != nil {
		ctx.Error(fmt.Sprintf("ошибка восстановления бренда с ID %d: %v", brandID, err), fasthttp.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(fmt.Sprintf("Бренд с ID %d успешно восстановлен", brandID)))
}
