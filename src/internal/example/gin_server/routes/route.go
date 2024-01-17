package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaoz194/GoFlux/src/internal/example/gin_server/controller"
	"github.com/xiaoz194/GoFlux/src/internal/example/gin_server/middleware"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	// v1 接口，提供给内部后端服务调用
	v1 := r.Group("/api/")
	v1.Use(middleware.RecoverMiddleware())
	{
		v1.GET("/v1/test_get/uid/:uid/", controller.TestGet)
		v1.POST("/v1/test_post/uid/:uid/", controller.TestPost)
	}

	// v3 接口 提供给前端访问
	v3 := r.Group("/api/")
	v3.Use(middleware.RecoverMiddleware())

	return r
}
