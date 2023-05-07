package routes

import (
	"server/controller"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 注册基础路由
func InitBaseRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	logController := controller.NewLogController()
	configController := controller.NewRaspConfigController()
	fileController := controller.NewFileController()
	router := r.Group("/base")
	{
		// 登录登出刷新token无需鉴权
		router.POST("/login", authMiddleware.LoginHandler)
		router.POST("/logout", authMiddleware.LogoutHandler)
		router.POST("/refreshToken", authMiddleware.RefreshHandler)
		router.POST("/report", logController.ReportLog)
		router.POST("/remote/config", configController.GetViperRaspConfig)
		router.GET("/file/download", fileController.Download)
	}
	return r
}
