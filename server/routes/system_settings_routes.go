package routes

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"server/controller"
	"server/middleware"
)

func InitSystemSettingRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	settingController := controller.NewSystemSettingController()
	router := r.Group("/settings")
	// 开启jwt认证中间件
	router.Use(authMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	router.Use(middleware.CasbinMiddleware())
	{
		router.POST("/update", settingController.Update)
		router.GET("/list", settingController.List)
		router.POST("/getProjectInfo", settingController.GetProjectInfo)
	}
	return r
}
