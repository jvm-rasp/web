package routes

import (
	"server/socket"

	"github.com/gin-gonic/gin"
)

// 注册基础路由
func InitSocketRoutes(r *gin.RouterGroup) gin.IRoutes {
	router := r.Group("/ws")
	{
		router.GET("/:hostName", socket.WebsocketManager.WsClient)
	}
	return r
}
