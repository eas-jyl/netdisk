package middleware

import (
	"net/http"
	"strings"

	"go-netdisk/internal/config"
	"go-netdisk/internal/service"

	"github.com/gin-gonic/gin"
)

/*
鉴权中间件
	1. 从请求头里面取出token
	2. 校验token是否合法
	3. 解析出当前用户信息
	4. 写入到gin.Context

*/

// 用户信息
const (
	ContextUserIDKey   = "userID"
	ContextUsernameKey = "username"
)

// 接口：从http中获取token， 解析出用户信息存入cin.context
func JWT(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取http中Authorization字段的值
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "missing authorization header",
			})
			return
		}

		// 判断是否属于是Bearer <token>格式
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "invalid authorization header",
			})
			return
		}

		// 从token解析出claims
		claims, err := service.ParseToken(parts[1], cfg.JWT.Secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "invalid or expired token",
			})
			return
		}

		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextUsernameKey, claims.Username)
		c.Next()

	}
}
