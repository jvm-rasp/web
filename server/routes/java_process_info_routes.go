package routes

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"server/controller"
	"server/middleware"
)

func InitJavaProcessInfoRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	javaPorcessInfoController := controller.NewJavaProcessInfoController()
	router := r.Group("/process")
	router.Use(authMiddleware.MiddlewareFunc())
	router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/list", javaPorcessInfoController.GetJavaProcessInfos)
	}
	return r
}
