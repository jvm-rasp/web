package controller

import (
	"github.com/gin-gonic/gin"
	"server/repository"
)

type IRaspLogController interface {
	HandleLog(c *gin.Context)
}

type RaspLogController struct {
	// 主机表、攻击表、进程表
	RaspConfigRepository repository.IRaspModuleRepository
}
