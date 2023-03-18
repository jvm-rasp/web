package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"server/common"
	"server/repository"
	"server/response"
	"server/vo"
)

type IRaspHostController interface {
	GetRaspHosts(c *gin.Context)
	BatchDeleteHostByIds(c *gin.Context)
}

type RaspHostController struct {
	RaspHostRepository        repository.IRaspHostRepository
}

func NewRaspHostController() IRaspHostController {
	raspHostRepository := repository.NewRaspHostRepository()
	raspHostController := RaspHostController{RaspHostRepository: raspHostRepository}
	return raspHostController
}

func (h RaspHostController) GetRaspHosts(c *gin.Context) {
	var req vo.RaspHostListRequest
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
	raspHosts, total, err := h.RaspHostRepository.GetRaspHosts(&req)
	if err != nil {
		response.Fail(c, nil, "获取实例列表失败")
		return
	}
	response.Success(c, gin.H{
		"data": raspHosts, "total": total,
	}, "获取实例列表失败")
}

// 批量删除接口
func (h RaspHostController) BatchDeleteHostByIds(c *gin.Context) {
	var req vo.DeleteRaspHostRequest
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
	// 删除接口
	err := h.RaspHostRepository.DeleteRaspHost(req.Ids)
	if err != nil {
		response.Fail(c, nil, "删除实例失败: "+err.Error())
		return
	}
	response.Success(c, nil, "删除实例成功")
}
