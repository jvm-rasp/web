package routes

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"server/controller"
	"server/middleware"
)

func InitRaspAttackRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	raspLogController := controller.NewRaspLogController()
	router := r.Group("/attack")
	// 开启jwt认证中间件
	router.Use(authMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/list", raspLogController.GetAttackLogs)
		router.GET("/detail", raspLogController.GetAttackDetail)
		router.POST("/delete/batch", raspLogController.BatchDeleteLogByIds)
	}
	return r
}
