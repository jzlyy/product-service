package middlewares

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "product_service_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "product_service_http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: []float64{0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path", "status"},
	)

	productOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "product_service_product_operations_total",
			Help: "Total number of product operations",
		},
		[]string{"operation", "status"},
	)

	categoryOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "product_service_category_operations_total",
			Help: "Total number of category operations",
		},
		[]string{"operation", "status"},
	)
)

// PrometheusMiddleware 收集 Prometheus 指标
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		httpRequestsTotal.WithLabelValues(
			c.Request.Method,
			path,
			status,
		).Inc()

		httpRequestDuration.WithLabelValues(
			c.Request.Method,
			path,
			status,
		).Observe(duration)
	}
}

// RecordProductOperation 记录商品操作指标
func RecordProductOperation(operation string, success bool) {
	status := "success"
	if !success {
		status = "error"
	}
	productOperations.WithLabelValues(operation, status).Inc()
}

// RecordCategoryOperation 记录分类操作指标
func RecordCategoryOperation(operation string, success bool) {
	status := "success"
	if !success {
		status = "error"
	}
	categoryOperations.WithLabelValues(operation, status).Inc()
}
