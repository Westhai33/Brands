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
	r            *router.Router
	log          zerolog.Logger
	brandHandler *handler.BrandHandler
	modelHandler *handler.ModelHandler
}

func NewService(
	log zerolog.Logger,
	bh *handler.BrandHandler,
	mh *handler.ModelHandler,
) (*service, error) {
	r := router.New()

	// Настройка пути для Swagger UI
	r.GET("/swagger", func(ctx *fasthttp.RequestCtx) {
		fastHttpSwagger.WrapHandler(fastHttpSwagger.InstanceName("swagger"))(ctx)
	})
	r.GET("/{filepath:*}", fasthttpadaptor.NewFastHTTPHandlerFunc(httpSwagger.WrapHandler))

	// Инициализация сервиса
	s := &service{
		log:          log,
		brandHandler: bh,
		modelHandler: mh,
	}

	// Настройка маршрутов
	s.brandHandler.SetupRoutes(r)
	s.modelHandler.SetupRoutes(r)

	s.r = r
	return s, nil
}

// Start запускает HTTP сервер
func (s *service) Start(ctx context.Context) error {
	server := fasthttp.Server{
		Handler: s.r.Handler,
		Name:    "Brands and Models API",
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
