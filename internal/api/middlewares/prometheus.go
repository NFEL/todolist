package middlewares

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)

	requestDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "http_request_duration_seconds",
			Help:       "HTTP request duration in seconds.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"method", "path", "status"},
	)
)

func init() {
	prometheus.MustRegister(requestCounter, requestDuration)
}

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		method := c.Request.Method
		path := c.Request.URL.Path
		status := c.Writer.Status()

		requestCounter.WithLabelValues(method, path, strconv.Itoa(status)).Inc()

		duration := time.Since(start).Seconds()
		requestDuration.WithLabelValues(method, path, strconv.Itoa(status)).Observe(duration)
	}
}
