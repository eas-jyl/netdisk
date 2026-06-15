package router

import (
	"github.com/gin-gonic/gin"
	"go-netdisk/internal/config"
	"go-netdisk/internal/handler"
	"go-netdisk/internal/middleware"
	"go-netdisk/internal/service"
	"gorm.io/gorm"
	"net/http"
)

// 接口：创建一个路由引擎
func New(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()
	authHandler := handler.NewAuthHandler(service.NewAuthService(db, cfg))
	userHandler := handler.NewUserHandler(service.NewUserService(db))

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	api := r.Group("/api/account")
	api.POST("/register", authHandler.Register)
	api.POST("/login", authHandler.Login)

	userGroup := r.Group("/api/user")
	userGroup.Use(middleware.JWT(cfg))
	userGroup.GET("/me", userHandler.Me)

	return r
}
