package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"server/common"
	"server/model"
	"server/repository"
	"server/response"
	"server/util"
	"server/vo"
)

type IRaspConfigController interface {
	CreateRaspConfig(c *gin.Context)
	GetRaspConfigs(c *gin.Context)
	BatchDeleteConfigByIds(c *gin.Context)
}

type RaspConfigController struct {
	RaspConfigRepository repository.IRaspConfigRepository
}

func NewRaspConfigController() IRaspConfigController {
	raspConfigRepository := repository.NewRaspConfigRepository()
	raspConfigController := RaspConfigController{RaspConfigRepository: raspConfigRepository}
	return raspConfigController
}

func (r RaspConfigController) GetRaspConfigs(c *gin.Context) {
	var req vo.RaspConfigListRequest
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
	raspConfigs, total, err := r.RaspConfigRepository.GetRaspConfigs(&req)
	if err != nil {
		response.Fail(c, nil, "获取配置列表失败")
		return
	}
	response.Success(c, gin.H{
		"list": raspConfigs, "total": total,
	}, "获取配置列表成功")
}

func (r RaspConfigController) CreateRaspConfig(c *gin.Context) {
	var req vo.CreateRaspConfigRequest
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

	// 获取当前用户
	ur := repository.NewUserRepository()
	ctxUser, err := ur.GetCurrentUser(c)
	if err != nil {
		response.Fail(c, nil, "获取当前用户信息失败")
		return
	}

	raspConfig := model.RaspConfig{
		Name:          req.Name,
		Desc:          req.Desc,
		Status:        req.Status,
		Creator:       ctxUser.Username,
		Operator:      ctxUser.Username,
		AgentMode:     req.AgentMode,
		ModuleConfigs: util.Struct2Json(req.ModuleConfigs),
		LogPath:       req.LogPath,
		AgentConfigs:  util.Struct2Json(req.AgentConfigs),
		BinFileUrl:    req.BinFileUrl,
		BinFileHash:   req.BinFileHash,
	}

	// 获取
	err = r.RaspConfigRepository.CreateRaspConfig(&raspConfig)
	if err != nil {
		response.Fail(c, nil, "获取接口列表失败"+err.Error())
		return
	}
	response.Success(c, nil, "创建配置成功")
	return
}

// 批量删除接口
func (r RaspConfigController) BatchDeleteConfigByIds(c *gin.Context) {
	var req vo.DeleteRaspConfigRequest
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
	err := r.RaspConfigRepository.DeleteRaspConfig(req.Ids)
	if err != nil {
		response.Fail(c, nil, "删除配置失败: "+err.Error())
		return
	}
	response.Success(c, nil, "删除配置成功")
}
