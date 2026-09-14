package httpapi

import (
	"http-server-projeto-korp/internal/httpapi/handler"
	"http-server-projeto-korp/internal/httpapi/middleware"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter() (*gin.Engine, error) {
	registry := prometheus.NewRegistry()
	router := gin.New()
	router.Use(gin.Logger(), middleware.Metrics(registry), gin.Recovery())

	if err := router.SetTrustedProxies(nil); err != nil {
		return nil, err
	}

	router.GET("/projeto-korp", handler.GetProjetoKorp)
	router.GET("/metrics", gin.WrapH(promhttp.HandlerFor(registry, promhttp.HandlerOpts{})))

	return router, nil
}
