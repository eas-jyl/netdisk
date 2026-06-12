package main

// 用户管理系统

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// 用户数据字段
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"gte=0,lte=150"`
}

// 创建用户请求
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"gte=0,lte=150"`
}

// 更新用户请求
type UpdateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"gte=0,lte=150"`
}

// 响应
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 批量声明变量
var (
	users = map[int]User{
		1: {ID: 1, Name: "Tom", Email: "tom@example.com", Age: 18},
		2: {ID: 2, Name: "Jerry", Email: "jerry@example.com", Age: 20},
	}
	nextID = 3
	mx     sync.Mutex
)

func main() {
	r := gin.Default()

	r.Use(RequestCostMiddleware())

	api := r.Group("/api/v1")
	{
		api.GET("/health", healthCheck)
		api.GET("/users", listUsers)
		api.POST("/users", createUser)
		api.GET("/users/:id", getUser)
		api.PUT("/users/:id", updateUser)
		api.DELETE("/users/:id", deleteUser)

		// 需要鉴权的接口
		auth := api.Group("")
		auth.Use(AuthMiddleware())
		{
			auth.GET("/profile", getProfile)
		}
	}

	r.Run(":8080")
}

// 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: 0, Message: "success", Data: data})
}

// 失败响应
func Fail(c *gin.Context, status int, code int, message string) {
	c.JSON(status, Response{Code: code, Message: message})
}

// 中间件：一次http请求处理花了多久
func RequestCostMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 执行后面的中间件和接口函数
		c.Next()

		// 计算耗时
		cost := time.Since(start)

		// 设置响应头
		c.Header("X-Request-Cost", cost.String())
	}
}

// 中间件：鉴权
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "Bearer demo-token" {
			Fail(c, http.StatusUnauthorized, 401, "unauthorized")

			c.Abort()
			return
		}

		// 在当前上下文中存入一个键值对
		c.Set("user_id", 1001)
		c.Next()
	}

}

func healthCheck(c *gin.Context) {
	Success(c, gin.H{"status": "ok"})
}

// 接口：将当前内存里面所有用户取出来， 返回给客户端
func listUsers(c *gin.Context) {
	mx.Lock()
	defer mx.Unlock()

	list := make([]User, 0, len(users))
	for _, user := range users {
		list = append(list, user)
	}

	Success(c, list)
}

// 接口：创建用户
func createUser(c *gin.Context) {
	var req CreateUserRequest

	// 将http中的请求body映射到req结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	mx.Lock()
	defer mx.Unlock()

	// 创建新用户
	user := User{
		ID:    nextID,
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}
	users[nextID] = user
	nextID++

	// 返回http响应
	c.JSON(http.StatusCreated, Response{Code: 0, Message: "created", Data: user})
}

// 接口：获取用户
func getUser(c *gin.Context) {

	// c.Param : 获取路由参数的id ， strconv : 转换为int类型
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		Fail(c, http.StatusBadRequest, 400, "invalid user id")
		return
	}

	mu.Lock()
	defer mu.Unlock()

	user, ok := users[id]
	if !ok {
		Fail(c, http.StatusNotFound, 404, "user not found")
		return
	}

	Success(c, user)
}

func updateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		Fail(c, http.StatusBadRequest, 400, "invalid user id")
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if _, ok := users[id]; !ok {
		Fail(c, http.StatusNotFound, 404, "user not found")
		return
	}

	user := User{ID: id, Name: req.Name, Email: req.Email, Age: req.Age}
	users[id] = user

	Success(c, user)
}

func deleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		Fail(c, http.StatusBadRequest, 400, "invalid user id")
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if _, ok := users[id]; !ok {
		Fail(c, http.StatusNotFound, 404, "user not found")
		return
	}

	delete(users, id)
	Success(c, gin.H{"deleted": true})
}

func getProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	Success(c, gin.H{
		"user_id": userID,
		"name":    "demo user",
	})
}
