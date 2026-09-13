package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetProjetoKorp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"nome":    "Projeto Korp",
		"horario": time.Now().UTC().Format(time.RFC3339Nano),
	})
}
