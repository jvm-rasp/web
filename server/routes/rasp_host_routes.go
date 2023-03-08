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
		router.DELETE("/delete/batch", raspHostController.BatchDeleteHostByIds)
	}
	return r
}