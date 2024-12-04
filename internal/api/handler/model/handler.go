package model

import (
	"Brands/internal/service/model"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type ModelHandler interface {
	SetupRoutes(r *router.Router)
	CreateModel(ctx *fasthttp.RequestCtx)
	DeleteModel(ctx *fasthttp.RequestCtx)
	RestoreModel(ctx *fasthttp.RequestCtx)
	GetAllModels(ctx *fasthttp.RequestCtx)
	GetModelByID(ctx *fasthttp.RequestCtx)
	UpdateModel(ctx *fasthttp.RequestCtx)
}

type Handler struct {
	ModelService model.ModelService
}

func New(modelService model.ModelService) ModelHandler {
	return &Handler{
		ModelService: modelService,
	}
}

func (api *Handler) SetupRoutes(r *router.Router) {
	group := r.Group("/models")
	group.POST("/create", api.CreateModel)
	group.GET("/{id}", api.GetModelByID)
	group.PUT("/update/{id}", api.UpdateModel)
	group.DELETE("/delete/{id}", api.DeleteModel)
	group.POST("/restore/{id}", api.RestoreModel)
	group.GET("/all", api.GetAllModels)
}
