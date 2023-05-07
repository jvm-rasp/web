package routes

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"server/controller"
	"server/middleware"
)

// 注册用户路由
func InitUserRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	userController := controller.NewUserController()
	router := r.Group("/user")
	// 开启jwt认证中间件
	router.Use(authMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	router.Use(middleware.CasbinMiddleware())
	{
		router.POST("/info", userController.GetUserInfo)
		router.GET("/list", userController.GetUsers)
		router.POST("/changePwd", userController.ChangePwd)
		router.POST("/create", userController.CreateUser)
		router.POST("/update/:userId", userController.UpdateUserById)
		router.POST("/delete/batch", userController.BatchDeleteUserByIds)
	}
	return r
}
