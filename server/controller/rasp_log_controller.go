package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"server/common"
	"server/repository"
	"server/response"
	"server/vo"
)

type IRaspLogController interface {
	GetRaspErrorLog(c *gin.Context)
	BatchDeleteLogByIds(c *gin.Context)
}

type RaspLogController struct {
	RaspLogRepository repository.IRaspErrorLogsRepository
}

func NewRaspLogController() IRaspLogController {
	errorLogsRepository := repository.NewRaspErrorLogsRepository()
	raspLogController := RaspLogController{RaspLogRepository: errorLogsRepository}
	return raspLogController
}

func (r RaspLogController) GetRaspErrorLog(c *gin.Context) {
	var req vo.RaspLogsListRequest
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
	raspHosts, total, err := r.RaspLogRepository.GetRaspLogs(&req)
	if err != nil {
		response.Fail(c, nil, "获取实例列表失败")
		return
	}
	response.Success(c, gin.H{
		"list": raspHosts, "total": total,
	}, "获取实例列表成功")
}

func (r RaspLogController) BatchDeleteLogByIds(c *gin.Context) {
	var req vo.RaspLogsDeleteRequest
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

	err := r.RaspLogRepository.DeleteRaspLogs(req.Ids)
	if err != nil {
		response.Fail(c, nil, "删除日志失败")
		return
	}
	response.Success(c, nil, "删除日志成功")
}
