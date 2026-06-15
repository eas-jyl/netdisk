package handler

import (
	"net/http"

	"go-netdisk/internal/middleware"
	"go-netdisk/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// 接口：根据用户id， 回传用户信息
func (h *UserHandler) Me(c *gin.Context) {
	// 获取用户id
	userIDValue, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		fail(c, http.StatusUnauthorized, "user not found in context")
		return
	}

	// 类型断言：转换为int
	userID, ok := userIDValue.(uint)
	if !ok {
		fail(c, http.StatusUnauthorized, "invalid user context")
		return
	}

	// 返回用户信息
	user, err := h.userService.GetByID(userID)
	if err != nil {
		fail(c, http.StatusInternalServerError, "query current user failed")
		return
	}

	success(c, userToResponse(user))
}
