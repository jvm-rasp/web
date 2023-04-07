package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"path"
	"server/common"
	"server/config"
	"server/model"
	"server/repository"
	"server/response"
	"server/socket"
	"server/util"
	"server/vo"
)

type IRaspHostController interface {
	GetRaspHosts(c *gin.Context)
	BatchDeleteHostByIds(c *gin.Context)
	PushConfig(c *gin.Context)
}

var AgentMode = map[uint]string{
	0: "disable",
	1: "dynamic",
	2: "static",
}

type RaspHostController struct {
	RaspHostRepository   repository.IRaspHostRepository
	RaspConfigRepository repository.IRaspConfigRepository
}

func NewRaspHostController() IRaspHostController {
	raspHostRepository := repository.NewRaspHostRepository()
	raspConfigRepository := repository.NewRaspConfigRepository()
	raspHostController := RaspHostController{
		RaspHostRepository:   raspHostRepository,
		RaspConfigRepository: raspConfigRepository,
	}
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

// 发布配置
func (h RaspHostController) PushConfig(c *gin.Context) {
	var req vo.PushConfigRequest
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

	// TODO 新建发布记录表
	configId := req.ConfigId
	raspConfig, err := h.RaspConfigRepository.GetRaspConfigById(configId)
	if err != nil {
		response.Fail(c, nil, "配置不存在")
		return
	}
	// 合并成一个完成的配置文件
	var agentConfigsFields model.AgentConfig
	var moduleConfigsFields []model.RaspModule
	var moduleConfigs []model.ModuleConfig
	err = json.Unmarshal([]byte(raspConfig.AgentConfigs.String()), &agentConfigsFields)
	if err != nil {
		response.Fail(c, nil, "获取配置文本失败:"+err.Error())
		return
	}
	err = json.Unmarshal([]byte(raspConfig.ModuleConfigs.String()), &moduleConfigsFields)
	if err != nil {
		response.Fail(c, nil, "获取配置文本失败:"+err.Error())
		return
	}
	for _, item := range moduleConfigsFields {
		var moduleConfig model.ModuleConfig
		err = json.Unmarshal([]byte(item.Parameters.String()), &moduleConfig)
		if err != nil {
			response.Fail(c, nil, "获取配置文本失败:"+err.Error())
			return
		}
		moduleConfigs = append(moduleConfigs, moduleConfig)
	}
	finalConfig := model.RaspFinalConfig{
		AgentMode:        AgentMode[raspConfig.AgentMode],
		Version:          "1.1.1",
		ConfigId:         raspConfig.ID,
		ModuleAutoUpdate: true,
		LogPath:          raspConfig.LogPath,
		RemoteHosts:      path.Join(fmt.Sprintf("%v:%v", util.GetDefaultIp(), config.Conf.System.Port), config.Conf.System.UrlPathPrefix),
		AgentConfigs:     agentConfigsFields,
		ModuleConfigs:    moduleConfigs,
	}

	content, err := json.Marshal(finalConfig)
	if err != nil {
		response.Fail(c, nil, "获取配置文本失败:"+err.Error())
		return
	}
	// 批量下发考虑 所有主机的状态
	for _, hostName := range req.HostNames {
		// 先判断连接是否存在
		m := socket.WebsocketManager.Group[hostName]
		if m != nil {
			client := m[hostName]
			if client != nil {
				socket.WebsocketManager.Send(hostName, hostName, content)
				response.Success(c, nil, "配置下发成功")
				return
			}
		}
		response.Fail(c, nil, hostName+",配置下发失败: wbsocket连接不存在，请检查rasp在线状态")
		return
	}
	response.Fail(c, nil, "配置下发失败")
}
