package routes

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"server/controller"
	"server/middleware"
)

func InitRaspDashboardRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	raspDashboardController := controller.NewRaspDashboardController()
	router := r.Group("/dashboard")
	router.Use(authMiddleware.MiddlewareFunc())
	router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/attackData", raspDashboardController.GetRaspAttackData)
		router.GET("/attackTrends", raspDashboardController.GetRaspAttackTrends)
		router.GET("/attackTypes", raspDashboardController.GetRaspAttackTypes)
	}
	return r
}
