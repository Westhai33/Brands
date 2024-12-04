package metrics

import (
	"Brands/pkg/zerohook"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"net/http"
	"regexp"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	TotalRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "total_requests",
		Help: "Общее количество HTTP-запросов",
	}, []string{"method", "path"})

	RequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "app_request_duration_seconds",
		Help:    "Время обработки HTTP-запросов в секундах",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"})

	ActiveGoroutines = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "active_goroutines",
		Help: "Количество активных горутин",
	})

	MemoryUsageBytes = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "memory_usage_bytes",
		Help: "Использование памяти приложением в байтах",
	})

	// TODO: Метрики кэша
	CacheHits = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "cache_hits_total",
		Help: "Общее количество попаданий в кэш",
	})

	CacheMisses = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "cache_misses_total",
		Help: "Общее количество промахов кэша",
	})
)

func RegisterMetrics() {
	_ = prometheus.Unregister(collectors.NewGoCollector())
	goCollector := collectors.NewGoCollector(
		collectors.WithGoCollectorRuntimeMetrics(
			collectors.GoRuntimeMetricsRule{Matcher: regexp.MustCompile("/cpu/classes/total:cpu-seconds")},
			collectors.GoRuntimeMetricsRule{Matcher: regexp.MustCompile("/cpu/classes/user:cpu-seconds")},
			collectors.GoRuntimeMetricsRule{Matcher: regexp.MustCompile("/memory/classes/heap/objects:bytes")},
			collectors.GoRuntimeMetricsRule{Matcher: regexp.MustCompile("/memory/classes/total:bytes")},
		),
	)
	prometheus.MustRegister(goCollector)

	prometheus.MustRegister(TotalRequests)
	prometheus.MustRegister(RequestDuration)

	prometheus.MustRegister(ActiveGoroutines)
	prometheus.MustRegister(MemoryUsageBytes)
	prometheus.MustRegister(CacheHits)
	prometheus.MustRegister(CacheMisses)
	go updateRuntimeMetrics()
}

// StartPrometheusServer запускает сервер для экспорта метрик Prometheus
func StartPrometheusServer(addr string) {
	RegisterMetrics()
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			zerohook.Logger.Error().Err(err).Msg("Ошибка при запуске сервера метрик")
		}
	}()
}

// updateRuntimeMetrics периодически обновляет метрики ActiveGoroutines и MemoryUsageBytes
func updateRuntimeMetrics() {
	for {
		ActiveGoroutines.Set(float64(runtime.NumGoroutine()))
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		MemoryUsageBytes.Set(float64(m.Alloc))
		time.Sleep(10 * time.Second)
	}
}
