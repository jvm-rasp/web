package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"server/common"
	"server/model"
	"server/repository"
	"server/response"
	"server/vo"
	"strings"
)

type IRaspConfigController interface {
	CreateRaspConfig(c *gin.Context)
	UpdateRaspConfig(c *gin.Context)
	GetRaspConfigs(c *gin.Context)
	BatchDeleteConfigByIds(c *gin.Context)
	GetViperRaspConfig(c *gin.Context)
	UpdateRaspConfigStatusById(c *gin.Context)
	UpdateRaspConfigDefaultById(c *gin.Context)
	PushRaspConfig(c *gin.Context)
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
		ModuleConfigs: req.ModuleConfigs,
		LogPath:       req.LogPath,
		AgentConfigs:  req.AgentConfigs,
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

func (r RaspConfigController) UpdateRaspConfig(c *gin.Context) {
	var req vo.UpdateRaspConfigRequest
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

	id := req.ID
	config, err := r.RaspConfigRepository.GetRaspConfigById(id)
	if err != nil {
		response.Fail(c, nil, "获取当前模块失败")
		return
	}

	config.Name = req.Name
	config.Desc = req.Desc
	config.Status = req.Status
	config.Operator = ctxUser.Username
	config.AgentMode = req.AgentMode
	config.ModuleConfigs = req.ModuleConfigs
	config.LogPath = req.LogPath
	config.AgentConfigs = req.AgentConfigs
	config.BinFileUrl = req.BinFileUrl
	config.BinFileHash = req.BinFileHash

	err = r.RaspConfigRepository.UpdateRaspConfig(config)
	if err != nil {
		response.Fail(c, nil, "更新当前配置失败")
		return
	}
	response.Success(c, nil, "更新配置成功")
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

// GetViperRaspConfig viper remote get
func (l RaspConfigController) GetViperRaspConfig(c *gin.Context) {
	var req vo.RaspConfigRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.Trans)
		response.Fail(c, nil, errStr)
		return
	}
	name := strings.TrimSpace(req.Key)
	config, err := l.RaspConfigRepository.GetRaspConfig(name)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, gin.H{"key": name, "value": config.AgentConfigs}, "获取配置成功")
}

func (l RaspConfigController) UpdateRaspConfigStatusById(c *gin.Context) {
	var req vo.UpdateRaspConfigStatusRequest
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

	id := req.ID
	config, err := l.RaspConfigRepository.GetRaspConfigById(id)
	if err != nil {
		response.Fail(c, nil, "更新当前策略失败")
		return
	}
	config.Status = !config.Status
	config.Operator = ctxUser.Username

	err = l.RaspConfigRepository.UpdateRaspConfig(config)
	if err != nil {
		response.Fail(c, nil, "更新当前策略失败")
		return
	}
	response.Success(c, nil, "更新策略成功")
}

func (l RaspConfigController) UpdateRaspConfigDefaultById(c *gin.Context) {
	var req vo.UpdateRaspConfigDefaultRequest
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
	// 全部设置为非默认
	if err := common.DB.Model(&model.RaspConfig{}).Where("1=1").Updates(map[string]interface{}{"is_default": false}).Error; err != nil {
		response.Fail(c, nil, "更新默认策略失败")
		return
	}

	// 获取当前用户
	ur := repository.NewUserRepository()
	ctxUser, err := ur.GetCurrentUser(c)
	if err != nil {
		response.Fail(c, nil, "获取当前用户信息失败")
		return
	}
	id := req.ID
	config, err := l.RaspConfigRepository.GetRaspConfigById(id)
	if err != nil {
		response.Fail(c, nil, "更新当前策略失败")
		return
	}
	config.IsDefault = req.IsDefault
	config.Operator = ctxUser.Username

	err = l.RaspConfigRepository.UpdateRaspConfig(config)
	if err != nil {
		response.Fail(c, nil, "更新当前策略失败")
		return
	}
	response.Success(c, nil, "更新策略成功")
}

func (l RaspConfigController) PushRaspConfig(c *gin.Context) {
	var req vo.PushRaspConfigRequest
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
	// 生成推送的配置信息
	hostController := NewRaspHostController()
	hostRepository := repository.NewRaspHostRepository()
	content, err := hostController.GeneratePushConfig(req.ID)
	if err != nil {
		response.Fail(c, nil, "推送配置失败")
		return
	}
	// 先更新数据库中的配置Id
	hostList, _, err := hostRepository.GetRaspHosts(new(vo.RaspHostListRequest))
	var hostNameList []string
	for _, item := range hostList {
		item.ConfigId = req.ID
		err = hostRepository.UpdateRaspHost(item)
		if err != nil {
			common.Log.Errorf("更新配置Id失败, %v", err)
			continue
		}
		hostNameList = append(hostNameList, item.HostName)
	}
	// 开始推送最新的配置
	hostController.PushHostsConfig(hostNameList, content)

	response.Success(c, nil, "推送策略成功")
}
