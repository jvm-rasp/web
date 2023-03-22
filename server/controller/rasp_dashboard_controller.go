package controller

import (
	"github.com/gin-gonic/gin"
	"server/repository"
	"server/response"
)

type IRaspDashboardController interface {
	GetRaspAttackData(c *gin.Context)
	GetRaspAttackTrends(c *gin.Context)
	GetRaspAttackTypes(c *gin.Context)
}

type RaspDashboardController struct {
	RaspDashboardRepository repository.IRaspDashboardRepository
}

func NewRaspDashboardController() IRaspDashboardController {
	raspDashboardRepository := repository.NewRaspDashboardRepository()
	raspDashboardController := RaspDashboardController{RaspDashboardRepository: raspDashboardRepository}
	return raspDashboardController
}

func (r RaspDashboardController) GetRaspAttackData(c *gin.Context) {
	response.Success(c, gin.H{
		"high": 93481, "total": 58213, "block": 16395,
	}, "获取攻击数据成功")
}
func (r RaspDashboardController) GetRaspAttackTrends(c *gin.Context) {
	date := []string{"2023-03-01", "2023-03-02", "2023-03-03", "2023-03-04", "2023-03-05", "2023-03-06", "2023-03-07"}
	warn := []int{751, 738, 460, 577, 675, 667, 658}
	block := []int{488, 620, 886, 649, 405, 411, 572}
	response.Success(c, gin.H{
		"date": date, "warn": warn, "block": block,
	}, "获取攻击趋势成功")
}
func (r RaspDashboardController) GetRaspAttackTypes(c *gin.Context) {
	sqlInject := map[string]interface{}{"name": "SQL注入", "value": 586}
	command := map[string]interface{}{"name": "命令执行", "value": 542}
	list := []map[string]interface{}{sqlInject, command}
	response.Success(c, gin.H{
		"list": list,
	}, "获取攻击趋势成功")
}
