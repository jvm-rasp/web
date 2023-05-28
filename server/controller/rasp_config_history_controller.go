package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"server/common"
	"server/repository"
	"server/response"
	"server/vo"
)

type IRaspConfigHistoryController interface {
	GetRaspConfigHistory(c *gin.Context)
	GetRaspConfigHistoryData(c *gin.Context)
}

type RaspConfigHistoryController struct {
	RaspConfigHistoryRepository repository.IRaspConfigHistoryRepository
}

func NewRaspConfigHistoryController() IRaspConfigHistoryController {
	repo1 := repository.NewRaspConfigHistoryRepository()
	raspConfigHistoryController := RaspConfigHistoryController{
		RaspConfigHistoryRepository: repo1,
	}
	return raspConfigHistoryController
}

func (this RaspConfigHistoryController) GetRaspConfigHistory(c *gin.Context) {
	var req vo.RaspConfigHistoryListRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.Trans)
		response.Fail(c, nil, errStr)
		return
	}
	// 获取
	raspConfigHistory, total, err := this.RaspConfigHistoryRepository.GetRaspConfigHistoryByGuid(req.ConfigGuid)
	if err != nil {
		response.Fail(c, nil, "获取配置列表失败")
		return
	}
	response.Success(c, gin.H{
		"list": raspConfigHistory, "total": total,
	}, "获取配置历史版本列表成功")
}

func (this RaspConfigHistoryController) GetRaspConfigHistoryData(c *gin.Context) {
	var req vo.RaspConfigHistoryDataRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.Trans)
		response.Fail(c, nil, errStr)
		return
	}
	// 获取
	raspConfigHistoryData, err := this.RaspConfigHistoryRepository.GetRaspConfigHistoryDataByGuid(req.ConfigGuid, req.Version)
	if err != nil {
		response.Fail(c, nil, "获取配置历史信息失败")
		return
	}
	response.Success(c,
		gin.H{
			"data": raspConfigHistoryData,
		}, "获取配置历史信息成功")
}
