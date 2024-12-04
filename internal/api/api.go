package api

import (
	"Brands/internal/api/handler/brand"
	"Brands/internal/api/handler/model"
	"Brands/pkg/zerohook"
	"context"
	"github.com/fasthttp/router"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
	"strings"
	"time"
)

type service struct {
	r            *router.Router
	log          zerolog.Logger
	brandHandler *brand.BrandHandler
	modelHandler *model.ModelHandler
}

func NewService(
	log zerolog.Logger,
	bh *brand.BrandHandler,
	mh *model.ModelHandler,
) (*service, error) {
	r := router.New()

	// Инициализация сервиса
	s := &service{
		log:          log,
		brandHandler: bh,
		modelHandler: mh,
	}
	// Health check маршрут
	r.GET("/health", func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SetBodyString("OK")
	})

	// Настройка маршрутов
	s.brandHandler.SetupRoutes(r)
	s.modelHandler.SetupRoutes(r)

	s.r = r
	return s, nil
}

func (s *service) Start(ctx context.Context) error {
	server := fasthttp.Server{
		Handler: RecoveryMiddleware(
			CORS(
				TraceMiddleware(
					LoggingMiddleware(s.r.Handler),
				),
			),
		),
		Name: "Brands and Models API",

		// The timeout values is also the maximum time it can take
		// for a complete page of Server-Sent Events (SSE).
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 3600 * time.Second,

		MaxRequestBodySize: 100_000,

		ErrorHandler: func(ctx *fasthttp.RequestCtx, err error) {
			if strings.Contains(err.Error(), "error when reading request headers: ") {
				// Suppress this particular error message
				return
			}
			zerohook.Logger.Error().
				Err(err).
				Msgf("Error when serving connection %q<->%q", ctx.RemoteAddr(), ctx.LocalAddr())
		},
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
