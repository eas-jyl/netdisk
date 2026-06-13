package handler

import (
	"errors"
	"net/http"
	"strings"

	"go-netdisk/internal/model"
	"go-netdisk/internal/service"

	"github.com/gin-gonic/gin"
)

// 作用：负责处理http请求

type AuthHandler struct {
	authService *service.AuthService
}

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type userResponse struct {
	ID       uint   `json:"id"`
	UUID     string `json:"uuid"`
	Username string `json:"username"`
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// 接口：处理客户端的注册请求
func (h *AuthHandler) Register(c *gin.Context) {
	// 获取请求中的用户名和密码
	var req authRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "invalid request body")
		return
	}
	// 去除空格
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)
	if req.Username == "" || req.Password == "" {
		fail(c, http.StatusBadRequest, "username and password are required")
		return
	}
	// 注册用户
	user, err := h.authService.Register(req.Username, req.Password)
	if err != nil {
		if err.Error() == "username already exists" {
			fail(c, http.StatusBadRequest, err.Error())
			return
		}
		fail(c, http.StatusInternalServerError, "register failed")
		return
	}

	success(c, userToResponse(user))
}

// 接口：处理客户端的登陆请求
func (h *AuthHandler) Login(c *gin.Context) {
	// 获取请求中的用户名和密码
	var req authRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)
	if req.Username == "" || req.Password == "" {
		fail(c, http.StatusBadRequest, "username and password are required")
		return
	}

	// 登陆用户
	user, token, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		// token 验证不正确
		if errors.Is(err, service.ErrInvalidCredentials) {
			fail(c, http.StatusUnauthorized, "invalid username or password")
			return
		}
		fail(c, http.StatusInternalServerError, "login failed")
		return
	}

	success(c, gin.H{
		"token": token,
		"user":  userToResponse(user),
	})
}

func success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, response{
		Code:    0,
		Message: "ok",
		Data:    data,
	})
}

func fail(c *gin.Context, status int, message string) {
	c.JSON(status, response{
		Code:    status,
		Message: message,
	})
}

func userToResponse(user *model.User) userResponse {
	return userResponse{
		ID:       user.ID,
		UUID:     user.UUID,
		Username: user.Username,
	}
}
