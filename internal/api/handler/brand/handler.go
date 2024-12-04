package brand

import (
	_ "Brands/docs"
	"Brands/internal/service/brand"
	"github.com/fasthttp/router"
)

type BrandHandler struct {
	BrandService *brand.BrandService
}

func New(brandService *brand.BrandService) *BrandHandler {
	return &BrandHandler{
		BrandService: brandService,
	}
}

func (api *BrandHandler) SetupRoutes(r *router.Router) {
	group := r.Group("/brands")
	group.POST("/create", api.CreateBrand)
	group.GET("/{id}", api.GetBrandByID)
	group.PUT("/update/{id}", api.UpdateBrand)
	group.DELETE("/delete/{id}", api.DeleteBrand)
	group.POST("/restore/{id}", api.RestoreBrand)
	group.GET("/all", api.GetAllBrands)
}
