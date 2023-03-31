package controller

import (
	"github.com/gin-gonic/gin"
	"server/common"
	"server/model"
	"server/repository"
	"server/response"
	"time"
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
	total, err := r.RaspDashboardRepository.GetAllCount()
	if err != nil {
		response.Fail(c, nil, "获取总攻击次数失败")
		return
	}
	high, err := r.RaspDashboardRepository.GetHighLevelCount(90)
	if err != nil {
		response.Fail(c, nil, "获取高危漏洞次数失败")
		return
	}
	block, err := r.RaspDashboardRepository.GetBlockCount()
	if err != nil {
		response.Fail(c, nil, "获取拦截漏洞次数失败")
		return
	}
	response.Success(c, gin.H{
		"high": high, "total": total, "block": block,
	}, "获取攻击数据成功")
}
func (r RaspDashboardController) GetRaspAttackTrends(c *gin.Context) {
	now := time.Now()
	date := []string{
		now.AddDate(0, 0, -6).Format("2006-01-02"),
		now.AddDate(0, 0, -5).Format("2006-01-02"),
		now.AddDate(0, 0, -4).Format("2006-01-02"),
		now.AddDate(0, 0, -3).Format("2006-01-02"),
		now.AddDate(0, 0, -2).Format("2006-01-02"),
		now.AddDate(0, 0, -1).Format("2006-01-02"),
		now.Format("2006-01-02"),
	}
	var warn []int64
	var block []int64
	for _, item := range date {
		var warnCount int64
		var blockCount int64
		common.DB.Model(&model.RaspAttack{}).
			Where("attack_time >= ?", item+" 00:00:00").
			Where("attack_time <= ?", item+" 23:59:59").
			Where("is_blocked", false).
			Count(&warnCount)
		warn = append(warn, warnCount)
		common.DB.Model(&model.RaspAttack{}).
			Where("attack_time >= ?", item+" 00:00:00").
			Where("attack_time <= ?", item+" 23:59:59").
			Where("is_blocked", true).
			Count(&blockCount)
		block = append(block, blockCount)
	}
	response.Success(c, gin.H{
		"date": date, "warn": warn, "block": block,
	}, "获取攻击趋势成功")
}
func (r RaspDashboardController) GetRaspAttackTypes(c *gin.Context) {
	list, err := r.RaspDashboardRepository.GetAttackTypes()
	if err != nil {
		response.Fail(c, nil, "获取攻击类型失败")
		return
	}
	response.Success(c, gin.H{
		"list": list,
	}, "获取攻击趋势成功")
}
