package routes

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"server/controller"
	"server/middleware"
)

func InitRaspLogRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	raspLogController := controller.NewRaspLogController()
	router := r.Group("/rasp-log")
	// 开启jwt认证中间件
	router.Use(authMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/list", raspLogController.GetRaspErrorLog)
		router.POST("/delete/batch", raspLogController.BatchDeleteLogByIds)
	}
	return r
}
