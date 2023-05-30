package routes

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"server/controller"
	"server/middleware"
)

func InitRaspHostRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	raspHostController := controller.NewRaspHostController()
	router := r.Group("/host")
	router.Use(authMiddleware.MiddlewareFunc())
	router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/list", raspHostController.GetRaspHosts)
		router.POST("/push/config", raspHostController.PushConfig)
		router.POST("/delete/batch", raspHostController.BatchDeleteHostByIds)
		router.POST("/update", raspHostController.UpdateConfig)
		router.POST("/add", raspHostController.AddHost)
	}
	return r
}
