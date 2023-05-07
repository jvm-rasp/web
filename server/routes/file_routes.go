package routes

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"server/controller"
	"server/middleware"
)

func InitFileUploadRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	controller := controller.NewFileController()
	router := r.Group("/file")
	router.Use(authMiddleware.MiddlewareFunc())
	router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/list", controller.GetRaspFiles)
		router.POST("/upload", controller.Upload)
		router.POST("/delete/batch", controller.Delete)
		router.GET("/getFileInfo/module", controller.GetModuleInfo)
	}
	return r
}
