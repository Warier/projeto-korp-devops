package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func Metrics(registry prometheus.Registerer) gin.HandlerFunc {
	requests := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "korp_http_requests_total",
		Help: "Total number of http responses by method, route, and status.",
	}, []string{"method", "route", "status"})

	duration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "korp_http_request_duration_seconds",
		Help:    "http request duration in seconds.",
		Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
	}, []string{"method", "route"})

	registry.MustRegister(requests, duration)

	return func(c *gin.Context) {
		if c.FullPath() == "/metrics" {
			c.Next()
			return
		}

		start := time.Now()
		c.Next()

		route := c.FullPath()
		if route == "" {
			route = "unmatched"
		}

		method := c.Request.Method
		switch method {
		case http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut,
			http.MethodPatch, http.MethodDelete, http.MethodOptions, http.MethodConnect, http.MethodTrace:
		default:
			method = "OTHER"
		}

		requests.WithLabelValues(method, route, strconv.Itoa(c.Writer.Status())).Inc()
		duration.WithLabelValues(method, route).Observe(time.Since(start).Seconds())
	}
}
