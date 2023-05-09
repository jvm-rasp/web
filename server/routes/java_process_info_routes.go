package routes

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"server/controller"
	"server/middleware"
	"server/repository"
)

func InitJavaProcessInfoRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	repo1 := repository.NewRaspHostRepository()
	javaPorcessInfoController := controller.NewJavaProcessInfoController(repo1)
	router := r.Group("/process")
	router.Use(authMiddleware.MiddlewareFunc())
	router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/list", javaPorcessInfoController.GetJavaProcessInfos)
	}
	return r
}
