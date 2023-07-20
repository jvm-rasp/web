package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net"
	"server/Interface"
	"server/common"
	"server/config"
	"server/global"
	"server/model"
	"server/repository"
	"server/response"
	"server/socket"
	"server/util"
	"server/vo"
	"strings"
	"time"
)

var AgentMode = map[uint]string{
	0: "disable",
	1: "dynamic",
	2: "static",
}

type RaspHostController struct {
	RaspHostRepository          repository.IRaspHostRepository
	JavaProcessInfoRepository   repository.IJavaProcessInfoRepository
	RaspConfigRepository        repository.IRaspConfigRepository
	RaspConfigHistoryRepository repository.IRaspConfigHistoryRepository
	RaspFileRepository          repository.IRaspFileRepository
	RaspModuleRepository        repository.IRaspModuleRepository
	RaspComponentRepository     repository.IRaspComponentRepository
}

func NewRaspHostController() Interface.IRaspHostController {
	if global.IRaspHostController == nil {
		repo1 := repository.NewRaspHostRepository()
		repo2 := repository.NewRaspConfigRepository()
		repo3 := repository.NewRaspFileRepository()
		repo4 := repository.NewRaspModuleRepository()
		repo5 := repository.NewRaspComponentRepository()
		repo6 := repository.NewRaspConfigHistoryRepository()
		repo7 := repository.NewJavaProcessInfoRepository(repo1)
		raspHostController := RaspHostController{
			RaspHostRepository:          repo1,
			RaspConfigRepository:        repo2,
			RaspFileRepository:          repo3,
			RaspModuleRepository:        repo4,
			RaspComponentRepository:     repo5,
			RaspConfigHistoryRepository: repo6,
			JavaProcessInfoRepository:   repo7,
		}
		raspHostController.InitPushConfigService()
		global.IRaspHostController = raspHostController
	}
	return global.IRaspHostController
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
	}, "获取实例列表成功")
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
	for _, item := range req.Ids {
		hostInfo, err := h.RaspHostRepository.GetRaspHostById(item)
		if err != nil {
			response.Fail(c, nil, "删除实例失败: "+err.Error())
			return
		}
		err = h.JavaProcessInfoRepository.DeleteProcessByHostName(hostInfo.HostName)
		if err != nil {
			response.Fail(c, nil, "删除实例失败: "+err.Error())
			return
		}
		err = h.RaspHostRepository.DeleteRaspHostById(hostInfo.ID)
		if err != nil {
			response.Fail(c, nil, "删除实例失败: "+err.Error())
			return
		}
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
	if raspConfig == nil {
		return nil, errors.New(fmt.Sprintf("获取配置文本失败: 未找到configId=%v的配置", configId))
	}
	raspConfigHistory, err := h.RaspConfigHistoryRepository.GetRaspConfigHistoryDataByGuid(raspConfig.RowGuid, raspConfig.Version)
	if err != nil {
		return nil, errors.New("获取配置版本失败:" + err.Error())
	}
	var historyModuleConfig []map[string]interface{}
	err = json.Unmarshal(raspConfigHistory.ModuleConfigs, &historyModuleConfig)
	if err != nil {
		return nil, errors.New("反序列化用户配置信息失败:" + err.Error())
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
	err = json.Unmarshal(raspConfigHistory.AgentConfigs, &agentConfigsFields)
	if err != nil {
		return nil, errors.New("获取配置文本失败:" + err.Error())
	}
	err = json.Unmarshal(raspConfigHistory.ModuleConfigs, &moduleConfigsFields)
	if err != nil {
		return nil, errors.New("获取配置文本失败:" + err.Error())
	}
	err = json.Unmarshal(raspConfigHistory.RaspLibInfo, &raspLibInfo)
	if err != nil {
		return nil, errors.New("获取配置文本失败:" + err.Error())
	}
	err = json.Unmarshal(raspConfigHistory.RaspBinInfo, &raspBinInfo)
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
		moduleInfo, err := h.RaspModuleRepository.GetRaspModuleByName(item.ModuleName)
		if err != nil {
			return nil, errors.New("获取模块信息失败:" + err.Error())
		}
		componentsInfo := item.Components
		// 如果是外部下载地址则直接赋值
		for _, component := range componentsInfo {
			var moduleConfig model.ModuleConfig
			moduleConfig.ModuleName = component.ComponentName
			moduleConfig.ModuleVersion = component.ComponentVersion
			moduleConfig.Md5 = component.Md5
			if strings.HasPrefix(component.DownLoadURL, "http://") || strings.HasPrefix(component.DownLoadURL, "https://") {
				moduleConfig.DownLoadUrl = component.DownLoadURL
			} else {
				moduleConfig.DownLoadUrl = fmt.Sprintf("%v%v", downloadPrefix, component.DownLoadURL)
			}
			var componentParameters map[string]interface{}
			err = json.Unmarshal(component.Parameters, &componentParameters)
			if err != nil {
				return nil, errors.New("反序列化组件配置失败:" + err.Error())
			}
			moduleConfig.Parameters = componentParameters
			// 使用策略中用户定义的拦截放行规则替换默认规则
			if component.ComponentType == 2 {
				actions := h.GetUserBlockParameterList(historyModuleConfig, moduleInfo.ModuleName)
				for k, v := range actions {
					moduleConfig.Parameters[k] = v
				}
			}
			moduleConfigs = append(moduleConfigs, moduleConfig)
		}
	}
	finalConfig := model.RaspFinalConfig{
		AgentMode:        AgentMode[raspConfigHistory.AgentMode],
		ConfigId:         raspConfig.ID,
		ModuleAutoUpdate: true,
		LogPath:          raspConfigHistory.LogPath,
		RemoteHosts: fmt.Sprintf("%v://%v:%v/%v",
			util.Ternary(config.Conf.Ssl.Enable, "wss", "ws"),
			util.GetDefaultIp(),
			config.Conf.System.Port,
			config.Conf.System.UrlPathPrefix),
		EnableMdns:     true,
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

func (h RaspHostController) GetUserBlockParameterList(userBlockParameter []map[string]interface{}, moduleName string) map[string]interface{} {
	var result = map[string]interface{}{}
	for _, item := range userBlockParameter {
		if item["moduleName"] == moduleName {
			parameters := item["parameters"].(map[string]interface{})
			result = parameters["action"].(map[string]interface{})
			return result
		}
	}
	return result
}

func (h RaspHostController) AddHost(c *gin.Context) {
	var req vo.AddHostRequest
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
	udpAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%v:%v", req.Ip, req.Port))
	if err != nil {
		response.Fail(c, nil, "添加主机失败, err: "+err.Error())
		return
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	message := fmt.Sprintf("%v://%v:%v/%v",
		util.Ternary(config.Conf.Ssl.Enable, "wss", "ws"),
		util.GetDefaultIp(),
		config.Conf.System.Port,
		config.Conf.System.UrlPathPrefix)
	pack := &socket.Package{
		Magic:     socket.MagicBytes,
		Version:   socket.PROTOCOL_VERSION,
		Type:      socket.UPDATE_SERVER,
		BodySize:  int32(len(message)),
		TimeStamp: time.Now().Unix(),
		Signature: socket.EmptySignature,
		Body:      []byte((message)),
	}
	buf := bytes.NewBuffer(nil)
	err = pack.Pack(buf)
	if err != nil {
		response.Fail(c, nil, "添加主机失败, err: "+err.Error())
		return
	}
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		response.Fail(c, nil, "添加主机失败, err: "+err.Error())
		return
	}
	response.Success(c, nil, "已发送请求, 等待主机上线")
}

func (h RaspHostController) InitPushConfigService() {
	global.PushConfigQueue = make(chan *vo.PushConfigRequest, 128)
	go h.handlePushConfig()
}

func (h RaspHostController) handlePushConfig() {
	for {
		select {
		case pushConfig, ok := <-global.PushConfigQueue:
			if !ok {
				common.Log.Warn("the chan of push config is closed")
				return
			}
			content, err := h.GeneratePushConfig(pushConfig.ConfigId)
			if err != nil {
				common.Log.Errorf("生成配置失败, error: %v", err)
				continue
			}
			h.PushHostsConfig(pushConfig.HostNames, content)
		}
	}
}
