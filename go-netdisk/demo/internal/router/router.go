package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 接口：创建一个路由引擎
func New() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	return r
}
