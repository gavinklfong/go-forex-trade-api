package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Define custom metrics
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of response time for handler in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "status"},
	)
)

func init() {
	// Register custom metrics with Prometheus
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

// PrometheusMiddleware is a Gin middleware for collecting Prometheus metrics
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process the request
		c.Next()

		// Collect metrics after the request is processed
		duration := time.Since(startTime).Seconds()
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.FullPath() // Use FullPath to capture the route pattern (e.g., "/users/:id")

		// Update metrics
		httpRequestsTotal.WithLabelValues(method, path, string(statusCode)).Inc()
		httpRequestDuration.WithLabelValues(method, path, string(statusCode)).Observe(duration)
	}
}
