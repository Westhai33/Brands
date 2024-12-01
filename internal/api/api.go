package api

import (
	"Brands/internal/api/handler"
	"context"
	"github.com/fasthttp/router"
	"github.com/rs/zerolog"
	fastHttpSwagger "github.com/swaggo/fasthttp-swagger"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

type service struct {
	r   *router.Router
	log zerolog.Logger

	brandHandler *handler.BrandHandler
	// TODO: modelHandler *ModelHandler
}

func NewService(
	log zerolog.Logger,
	bh *handler.BrandHandler,
	// TODO: mh  *ModelHandler,
) (*service, error) {
	r := router.New()
	// Настройка пути для Swagger UI
	r.GET("/swagger", func(ctx *fasthttp.RequestCtx) {
		fastHttpSwagger.WrapHandler(fastHttpSwagger.InstanceName("swagger"))(ctx)
	})
	r.GET("/{filepath:*}", fasthttpadaptor.NewFastHTTPHandlerFunc(httpSwagger.WrapHandler))

	s := &service{
		log:          log,
		brandHandler: bh,
		// TODO: modelsHandler: mh,
	}
	s.brandHandler.SetupRoutes(r)
	// TODO: s.modelHandler.SetupRoutes(r)

	s.r = r
	return s, nil
}

// CORS middleware
//func corsMiddleware(ctx *fasthttp.RequestCtx) {
//	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*") // Разрешаем все источники
//	ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
//	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin")
//
//	if string(ctx.Method()) == "OPTIONS" {
//		ctx.Response.SetStatusCode(fasthttp.StatusNoContent) // Ожидаемый ответ на OPTIONS запрос
//		return
//	}
//
//	ctx.Next() // Передаем запрос дальше в цепочку
//}

func (s *service) Start(ctx context.Context) error {
	server := fasthttp.Server{
		Handler: s.r.Handler,
		Name:    "Brands API",
	}
	emergencyShutdown := make(chan error)
	go func() {
		err := server.ListenAndServe(":8080")
		emergencyShutdown <- err
	}()
	select {
	case <-ctx.Done():
		return server.Shutdown()
	case e := <-emergencyShutdown:
		return e
	}
}
