package observability

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ==================== Регистрация кастомных метрик ====================
var (
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "method", "status"},
	)

	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method"},
	)
)

func InitMetrics() {
	// Регистрируем метрики
	prometheus.MustRegister(HTTPRequestsTotal, HTTPRequestDuration)
}

// ==================== HTTP handler для /metrics ====================
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}

// ==================== Middleware для Gin ====================
func GinMetricsMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		path := c.FullPath()
		method := c.Request.Method

		timer := prometheus.NewTimer(HTTPRequestDuration.WithLabelValues(path, method))
		defer timer.ObserveDuration()

		c.Next()

		status := fmt.Sprintf("%d", c.Writer.Status())
		HTTPRequestsTotal.WithLabelValues(path, method, status).Inc()
	}
}
