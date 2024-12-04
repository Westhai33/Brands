package brand

import (
	_ "Brands/docs"
	"Brands/internal/service/brand"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type BrandHandler interface {
	SetupRoutes(r *router.Router)
	CreateBrand(ctx *fasthttp.RequestCtx)
	DeleteBrand(ctx *fasthttp.RequestCtx)
	RestoreBrand(ctx *fasthttp.RequestCtx)
	GetAllBrands(ctx *fasthttp.RequestCtx)
	GetBrandByID(ctx *fasthttp.RequestCtx)
	UpdateBrand(ctx *fasthttp.RequestCtx)
}

type Handler struct {
	BrandService brand.BrandService
}

func New(brandService brand.BrandService) BrandHandler {
	return &Handler{
		BrandService: brandService,
	}
}

func (api *Handler) SetupRoutes(r *router.Router) {
	group := r.Group("/brands")
	group.POST("/create", api.CreateBrand)
	group.GET("/{id}", api.GetBrandByID)
	group.PUT("/update/{id}", api.UpdateBrand)
	group.DELETE("/delete/{id}", api.DeleteBrand)
	group.POST("/restore/{id}", api.RestoreBrand)
	group.GET("/all", api.GetAllBrands)
}
