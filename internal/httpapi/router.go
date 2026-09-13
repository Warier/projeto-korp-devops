package httpapi

import (
	"http-server-projeto-korp/internal/httpapi/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter() (*gin.Engine, error) {
	router := gin.Default()

	if err := router.SetTrustedProxies(nil); err != nil {
		return nil, err
	}

	router.GET("/projeto-korp", handler.GetProjetoKorp)

	return router, nil
}
