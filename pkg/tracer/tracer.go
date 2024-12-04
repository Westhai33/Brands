package tracer

import (
	"fmt"
	"io"

	"Brands/pkg/zerohook"
	"github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	traceconfig "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics/prometheus"
)

var Tracer opentracing.Tracer
var closer io.Closer

// InitTracer инициализирует глобальный трейсер
func InitTracer(serviceName, agentHost string, agentPort int) error {
	cfg := traceconfig.Configuration{
		ServiceName: serviceName,
		Sampler: &traceconfig.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1, // 1 включает полную выборку
		},
		Reporter: &traceconfig.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: fmt.Sprintf("%s:%d", agentHost, agentPort),
		},
	}

	metricsFactory := prometheus.New()

	var err error
	Tracer, closer, err = cfg.NewTracer(
		traceconfig.Logger(jaeger.StdLogger),
		traceconfig.Metrics(metricsFactory),
	)
	if err != nil {
		zerohook.Logger.Error().Err(err).Msg("Не удалось инициализировать трейсер")
		return err
	}
	opentracing.SetGlobalTracer(Tracer)
	zerohook.Logger.Info().Msg("Трейсер инициализирован")
	return nil
}

// CloseTracer закрывает трейсер
func CloseTracer() error {
	if closer != nil {
		return closer.Close()
	}
	return nil
}
