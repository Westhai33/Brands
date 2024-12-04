package model

import (
	"Brands/internal/service/model"
	"github.com/fasthttp/router"
)

type ModelHandler struct {
	ModelService *model.ModelService
}

func New(modelService *model.ModelService) *ModelHandler {
	return &ModelHandler{
		ModelService: modelService,
	}
}

func (api *ModelHandler) SetupRoutes(r *router.Router) {
	group := r.Group("/models")
	group.POST("/create", api.CreateModel)
	group.GET("/{id}", api.GetModelByID)
	group.PUT("/update/{id}", api.UpdateModel)
	group.DELETE("/delete/{id}", api.DeleteModel)
	group.POST("/restore/{id}", api.RestoreModel)
	group.GET("/all", api.GetAllModels)
	group.GET("/filter", api.ModelsFilter)
}
