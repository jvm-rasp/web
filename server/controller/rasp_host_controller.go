package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"server/common"
	"server/config"
	"server/model"
	"server/repository"
	"server/response"
	"server/socket"
	"server/util"
	"server/vo"
	"strings"
)

type IRaspHostController interface {
	GetRaspHosts(c *gin.Context)
	BatchDeleteHostByIds(c *gin.Context)
	PushConfig(c *gin.Context)
	UpdateConfig(c *gin.Context)
	PushHostsConfig(hostList []string, content []byte) []string
	GeneratePushConfig(configId uint) ([]byte, error)
}

var AgentMode = map[uint]string{
	0: "disable",
	1: "dynamic",
	2: "static",
}

type RaspHostController struct {
	RaspHostRepository   repository.IRaspHostRepository
	RaspConfigRepository repository.IRaspConfigRepository
	RaspFileRepository   repository.IRaspFileRepository
	RaspModuleRepository repository.IRaspModuleRepository
}

func NewRaspHostController() IRaspHostController {
	raspHostRepository := repository.NewRaspHostRepository()
	raspConfigRepository := repository.NewRaspConfigRepository()
	raspFileRepository := repository.NewRaspFileRepository()
	raspModuleRepository := repository.NewRaspModuleRepository()
	raspHostController := RaspHostController{
		RaspHostRepository:   raspHostRepository,
		RaspConfigRepository: raspConfigRepository,
		RaspFileRepository:   raspFileRepository,
		RaspModuleRepository: raspModuleRepository,
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
	// 生成推送的配置
	configId := req.ConfigId
	content, err := h.GeneratePushConfig(configId)
	if err != nil {
		response.Fail(c, nil, "生成配置文件失败")
		return
	}
	h.PushHostsConfig(req.HostNames, content)
	response.Success(c, nil, "配置下发成功")
	return
}

func (h RaspHostController) PushHostsConfig(hostList []string, content []byte) []string {
	var offlineHosts []string
	// 批量下发考虑 所有主机的状态
	for _, hostName := range hostList {
		// 先判断连接是否存在
		m := socket.WebsocketManager.Group[hostName]
		if m != nil {
			client := m[hostName]
			if client != nil {
				socket.WebsocketManager.Send(hostName, hostName, content)
			}
		} else {
			common.Log.Warnf("主机: %v,配置下发失败: websocket连接不存在，请检查rasp在线状态", hostName)
			offlineHosts = append(offlineHosts, hostName)
		}
	}
	return offlineHosts
}

func (h RaspHostController) GeneratePushConfig(configId uint) ([]byte, error) {
	raspConfig, err := h.RaspConfigRepository.GetRaspConfigById(configId)
	if err != nil {
		return nil, errors.New("获取配置文本失败:" + err.Error())
	}
	// 合并成一个完成的配置文件
	var agentConfigsFields model.AgentConfig
	var raspLibInfo model.ZipFileInfo
	var raspBinInfo model.ZipFileInfo

	var moduleConfigsFields []model.RaspModule
	var moduleConfigs []model.ModuleConfig
	var downloadPrefix = fmt.Sprintf("%v://%v:%v/%v",
		util.Ternary(config.Conf.Ssl.Enable, "https", "http"),
		util.GetDefaultIp(),
		config.Conf.System.Port,
		config.Conf.System.UrlPathPrefix)
	err = json.Unmarshal([]byte(raspConfig.AgentConfigs.String()), &agentConfigsFields)
	if err != nil {
		return nil, errors.New("获取配置文本失败:" + err.Error())
	}
	err = json.Unmarshal([]byte(raspConfig.ModuleConfigs.String()), &moduleConfigsFields)
	if err != nil {
		return nil, errors.New("获取配置文本失败:" + err.Error())
	}
	err = json.Unmarshal([]byte(raspConfig.RaspLibInfo.String()), &raspLibInfo)
	if err != nil {
		return nil, errors.New("获取配置文本失败:" + err.Error())
	}
	err = json.Unmarshal([]byte(raspConfig.RaspBinInfo.String()), &raspBinInfo)
	if err != nil {
		return nil, errors.New("获取配置文本失败:" + err.Error())
	}
	// 添加raspLibInfo子文件信息
	if raspLibInfo.Md5 != "" {
		raspLibInfo.DownloadUrl = util.Ternary(raspLibInfo.Md5 == "", "", fmt.Sprintf("%v%v", downloadPrefix, raspLibInfo.DownloadUrl)).(string)
		file, err := h.RaspFileRepository.GetRaspFileByHash(raspLibInfo.Md5)
		if err != nil {
			return nil, errors.New("获取配置文本失败:" + err.Error())
		}
		zipItemInfo, err := util.GetZipItemInfo(file.DiskPath)
		if err != nil {
			return nil, errors.New("获取配置文本失败:" + err.Error())
		}
		raspLibInfo.ItemsInfo = zipItemInfo
	}
	// 添加raspBinInfo子文件信息
	if raspBinInfo.Md5 != "" {
		raspBinInfo.DownloadUrl = util.Ternary(raspBinInfo.Md5 == "", "", fmt.Sprintf("%v%v", downloadPrefix, raspBinInfo.DownloadUrl)).(string)
		file, err := h.RaspFileRepository.GetRaspFileByHash(raspBinInfo.Md5)
		if err != nil {
			return nil, errors.New("获取配置文本失败:" + err.Error())
		}
		zipItemInfo, err := util.GetZipItemInfo(file.DiskPath)
		if err != nil {
			return nil, errors.New("获取配置文本失败:" + err.Error())
		}
		raspBinInfo.ItemsInfo = zipItemInfo
	}
	// 添加模块参数信息
	for _, item := range moduleConfigsFields {
		moduleId := item.ID
		moduleInfo, err := h.RaspModuleRepository.GetRaspModuleById(moduleId)
		if err != nil {
			return nil, errors.New("获取模块信息失败:" + err.Error())
		}
		var moduleConfig model.ModuleConfig
		err = json.Unmarshal([]byte(moduleInfo.Parameters.String()), &moduleConfig)
		// 如果是外部下载地址则直接赋值
		if strings.HasPrefix(moduleInfo.DownLoadURL, "http://") || strings.HasPrefix(moduleInfo.DownLoadURL, "https://") {
			moduleConfig.DownLoadUrl = moduleInfo.DownLoadURL
		} else {
			moduleConfig.DownLoadUrl = fmt.Sprintf("%v%v", downloadPrefix, moduleInfo.DownLoadURL)
		}
		moduleConfig.Md5 = moduleInfo.Md5
		if err != nil {
			return nil, errors.New("获取配置文本失败:" + err.Error())
		}
		moduleConfigs = append(moduleConfigs, moduleConfig)
	}
	finalConfig := model.RaspFinalConfig{
		AgentMode:        AgentMode[raspConfig.AgentMode],
		ConfigId:         raspConfig.ID,
		ModuleAutoUpdate: true,
		LogPath:          raspConfig.LogPath,
		RemoteHosts: fmt.Sprintf("%v://%v:%v/%v",
			util.Ternary(config.Conf.Ssl.Enable, "wss", "ws"),
			util.GetDefaultIp(),
			config.Conf.System.Port,
			config.Conf.System.UrlPathPrefix),
		EnableMdns:     config.Conf.Mdns.Enable,
		AgentConfigs:   agentConfigsFields,
		RaspLibConfigs: raspLibInfo,
		RaspBinConfigs: raspBinInfo,
		ModuleConfigs:  moduleConfigs,
	}

	content, err := json.Marshal(finalConfig)
	if err != nil {
		return nil, errors.New("获取配置文本失败:" + err.Error())
	}
	return content, nil
}

// 更新实例配置文件
func (h RaspHostController) UpdateConfig(c *gin.Context) {
	var req vo.UpdateRaspHostRequest
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
	host, err := h.RaspHostRepository.GetRaspHostById(req.Id)
	if err != nil {
		response.Fail(c, nil, "获取实例失败")
		return
	}
	host.HostName = req.HostName
	host.ConfigId = req.ConfigId
	err = h.RaspHostRepository.UpdateRaspHost(host)
	if err != nil {
		response.Fail(c, nil, "更新实例失败")
		return
	}
	response.Success(c, nil, "更新实例成功")
}
