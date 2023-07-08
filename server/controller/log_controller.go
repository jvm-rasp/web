package controller

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/vjeantet/grok"
	"gorm.io/datatypes"
	"server/common"
	"server/global"
	"server/model"
	"server/repository"
	"server/response"
	"server/util"
	"server/vo"
	"strconv"
	"strings"
	"time"
)

// grok
const pattern = "%{TIMESTAMP_ISO8601:time}\\s*%{LOGLEVEL:level}\\s*%{DATA:host}\\s*\\[%{DATA:thread}\\]\\s*\\[%{DATA:api}\\]\\s*%{GREEDYDATA:message}"

type ILogController interface {
	ReportLog(c *gin.Context)
	GetAttackLogs(c *gin.Context)
	GetAttackDetail(c *gin.Context)
	BatchDeleteLogByIds(c *gin.Context)
	UpdateStatusById(c *gin.Context)
}

type LogController struct {
	RaspHostRepository     repository.IRaspHostRepository
	JavaProcessRepository  repository.IJavaProcessInfoRepository
	RaspAttackRepository   repository.IRaspAttackRepository
	RaspErrorRepository    repository.IRaspErrorLogsRepository
	RaspConfigRepository   repository.IRaspConfigRepository
	HostResourceRepository repository.IHostResourceRepository
}

func NewLogController() ILogController {
	repository1 := repository.NewRaspHostRepository()
	repository2 := repository.NewJavaProcessInfoRepository(repository1)
	repository3 := repository.NewRaspAttackRepository()
	repository4 := repository.NewRaspErrorLogsRepository()
	repository5 := repository.NewRaspConfigRepository()
	repository6 := repository.NewHostResourceRepository()
	controller := LogController{
		RaspHostRepository:     repository1,
		JavaProcessRepository:  repository2,
		RaspAttackRepository:   repository3,
		RaspErrorRepository:    repository4,
		RaspConfigRepository:   repository5,
		HostResourceRepository: repository6,
	}
	return controller
}

func (l LogController) ReportLog(c *gin.Context) {
	var req vo.RaspLogRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	switch req.Fields.KafkaTopic {
	case vo.JRASP_DAEMON:
		l.handleDaemonLog(req)
	case vo.JRASP_AGENT:
	case vo.JRASP_MODULE:
	case vo.JRASP_ATTACK:
		l.handleAttackLog(req)
	default:
		panic(errors.New("unknown topic: " + req.Fields.KafkaTopic))
	}
	l.handleErrorLog(req.Fields.KafkaTopic, req)
}

func (l LogController) GetAttackLogs(c *gin.Context) {
	var req vo.RaspAttackListRequest
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
	raspHosts, total, err := l.RaspAttackRepository.GetRaspAttacks(&req)
	if err != nil {
		response.Fail(c, nil, "获取实例列表失败")
		return
	}
	response.Success(c, gin.H{
		"list": raspHosts, "total": total,
	}, "获取攻击日志列表失败")
}

func (l LogController) GetAttackDetail(c *gin.Context) {
	var req vo.RaspAttackDetailRequest
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
	record, err := l.RaspAttackRepository.GetRaspAttackDetail(req.Id)
	if err != nil {
		response.Fail(c, nil, "获取实例列表失败")
		return
	}
	response.Success(c, gin.H{
		"detail": record,
	}, "获取攻击日志列表失败")
}

// 批量删除操作日志
func (l LogController) BatchDeleteLogByIds(c *gin.Context) {
	var req vo.DeleteRaspAttackRequest
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
	err := l.RaspAttackRepository.DeleteRaspAttack(req.Guids)
	if err != nil {
		response.Fail(c, nil, "删除日志失败: "+err.Error())
		return
	}
	err = l.RaspAttackRepository.DeleteRaspDetail(req.Guids)
	if err != nil {
		response.Fail(c, nil, "删除日志失败: "+err.Error())
		return
	}

	response.Success(c, nil, "删除日志成功")
}

func (l LogController) UpdateStatusById(c *gin.Context) {
	var req vo.UpdateRaspStatusRequest
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
	// 更新状态
	id := req.Id
	attack, err := l.RaspAttackRepository.GetRaspAttackById(id)
	if err != nil {
		response.Fail(c, nil, "更新状态失败")
		return
	}
	attack.HandleResult = req.Result
	err = l.RaspAttackRepository.UpdateRaspAttack(attack)
	if err != nil {
		response.Fail(c, nil, "更新状态失败")
		return
	}
	response.Success(c, nil, "更新状态成功")
}

// 详情处理过程

// daemon 1000~1999
const (
	DAEMON_STARTUP_LOGID  = 1000
	HOST_ENV_LOGID        = 1002
	HEART_BEAT_LOGID      = 1011
	DEPENDENCY_JAR_LOGID  = 1016
	AGENT_SUCCESS_UNLOAD  = 1017
	JAVA_PROCESS_STARTUP  = 1018
	JAVA_PROCESS_SHUTDOWN = 1019
	AGENT_SUCCESS_INIT    = 1020
	NACOS_INIT_INFO       = 1024
	Agent_CONFIG_UPDATE   = 1025
	CONFIG_ID             = 1030
	RESOURCE_NAME_UPDATE  = 1033
)

func (l LogController) handleDaemonLog(req vo.RaspLogRequest) {
	// 不同logid 处理
	switch req.LogId {
	case DAEMON_STARTUP_LOGID:
		l.handleStartupLog(req)
	case HOST_ENV_LOGID:
		l.handleHostEnvLog(req)
	case HEART_BEAT_LOGID:
		l.handleHeartbeatLog(req)
	case AGENT_SUCCESS_INIT:
	case AGENT_SUCCESS_UNLOAD:
		l.handleAgentInitAndUnloadLog(req)
	case JAVA_PROCESS_STARTUP:
		l.handleFindJavaProcessLog(req)
	case JAVA_PROCESS_SHUTDOWN:
		l.handleRemoveJavaProcessLog(req)
	case Agent_CONFIG_UPDATE:
		l.handleAgentConfigUpdateLog(req)
	case CONFIG_ID:
		l.handleUpdateConfigId(req)
	case RESOURCE_NAME_UPDATE:
		l.handleUpdateResourceName(req)
	default:
	}
}

func (l LogController) handleStartupLog(req vo.RaspLogRequest) {
	host := &model.RaspHost{
		HostName:      req.HostName,
		Ip:            req.Ip,
		HeartbeatTime: req.Ts,
	}
	dbData, err := l.RaspHostRepository.QueryRaspHost(host.HostName)
	if err != nil {
		panic(err)
	}

	// 获取 agentMode
	detailMap := make(map[string]string)
	err = json.Unmarshal([]byte(req.Detail), &detailMap)
	if err != nil {
		panic(err)
	}
	host.AgentMode = detailMap["agentMode"]

	if len(dbData) <= 0 {
		return
	}

	err = l.RaspHostRepository.UpdateRaspHostByHostName(host)
	if err != nil {
		panic(err)
	}
}

func (l LogController) handleHostEnvLog(req vo.RaspLogRequest) {
	detailMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(req.Detail), &detailMap)
	if err != nil {
		panic(err)
	}

	installDir := detailMap["installDir"].(string)
	version := detailMap["version"].(string)
	binFileHash := detailMap["binFileHash"].(string)
	osType := detailMap["osType"].(string)

	totalMem := detailMap["totalMem"].(float64)
	cpuCounts := detailMap["cpuCounts"].(float64)
	freeDisk := detailMap["freeDisk"].(float64)

	buildDateTime := detailMap["buildDateTime"].(string)
	buildGitBranch := detailMap["buildGitBranch"].(string)
	buildGitCommit := detailMap["buildGitCommit"].(string)

	dbData, err := l.RaspHostRepository.QueryRaspHost(req.HostName)
	if err != nil {
		panic(err)
		return
	}

	var host *model.RaspHost
	if len(dbData) == 0 {
		host = &model.RaspHost{
			HostName: req.HostName,
			Ip:       req.Ip,
		}
	} else {
		host = dbData[0]
	}

	host.InstallDir = installDir
	host.Version = version
	host.ExeFileHash = binFileHash
	host.OsType = osType
	host.TotalMem = totalMem
	host.CpuCounts = cpuCounts
	host.FreeDisk = freeDisk
	host.BuildDateTime = buildDateTime
	host.BuildGitBranch = buildGitBranch
	host.BuildGitBranch = buildGitCommit

	if len(dbData) == 0 {
		return
	} else {
		err := l.RaspHostRepository.UpdateRaspHostByHostName(host)
		if err != nil {
			panic(err)
		}
	}
}

func (l LogController) handleHeartbeatLog(req vo.RaspLogRequest) {
	hostInfo, err := l.RaspHostRepository.GetRaspHostByHostName(req.HostName)
	if err != nil {
		panic(err)
	}

	if hostInfo == nil {
		hostInfo = &model.RaspHost{
			HostName: req.HostName,
			Ip:       req.Ip,
		}
		configId, err := l.RaspHostRepository.CreateRaspHost(hostInfo)
		if err != nil {
			panic(err)
		}
		// 推送默认配置
		if configId != 0 {
			global.PushConfigQueue <- &vo.PushConfigRequest{
				ConfigId:  configId,
				HostNames: []string{hostInfo.HostName},
			}
		}
	} else {
		err = l.RaspHostRepository.UpdateRaspHostByHostName(hostInfo)
		if err != nil {
			panic(err)
		}
	}

	// db 数据
	dbList, err := l.JavaProcessRepository.GetAllJavaProcessInfos(req.HostName)
	if err != nil {
		panic(err)
	}

	// 上报数据
	detailMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(req.Detail), &detailMap)
	if err != nil {
		panic(err)
	}

	for _, v := range dbList {
		// db有，上报数据中没有，删除db数据
		if len(detailMap) == 0 || detailMap[strconv.Itoa(v.Pid)] == nil {
			err := l.JavaProcessRepository.DeleteProcess(v.ID)
			if err != nil {
				panic(err)
			}
			continue
		}
		// db有,上报数据中也有，无需处理了
	}

	// db没有，上报数据中有，新增db数据
	for k, v := range detailMap {
		var existed = false
		for _, v2 := range dbList {
			if strconv.Itoa(v2.Pid) == k {
				existed = true
			}
		}
		if existed {
			continue
		}
		processDetail := v.(map[string]interface{})

		// todo json 转map int ---> float
		pidFloat := processDetail["pid"].(float64)
		pidint64 := strconv.FormatInt(int64(int(pidFloat)), 10)
		pid, _ := strconv.ParseInt(pidint64, 10, 32)

		startTime := processDetail["startTime"].(string)
		message := processDetail["status"].(string)
		// 缺少命令后信息
		process := &model.JavaProcessInfo{HostName: req.HostName, Pid: int(pid), StartTime: startTime, Status: l.convertMessageToStatus(message), Message: message}
		err = l.JavaProcessRepository.SaveProcessInfo(process)
		if err != nil {
			panic(err)
		}
	}

}

func (l LogController) handleAgentInitAndUnloadLog(req vo.RaspLogRequest) {
	detailMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(req.Detail), &detailMap)
	if err != nil {
		panic(err)
	}

	pid := detailMap["pid"].(float64)
	startTime := detailMap["startTime"].(string)
	messages := detailMap["status"].(string)
	status := l.convertMessageToStatus(messages)
	processInfo, err := l.JavaProcessRepository.GetProcessByPid(req.HostName, uint(pid))
	if processInfo != nil {
		processInfo.Status = status
		processInfo.Message = messages
		processInfo.StartTime = startTime
		err = l.JavaProcessRepository.UpdateProcessInfo(processInfo)
		if err != nil {
			panic(err)
		}
	} else {
		err = l.JavaProcessRepository.DeleteProcessByPid(req.HostName, uint(pid))
		if err != nil {
			panic(err)
		}
	}
}

func (l LogController) handleAgentConfigUpdateLog(req vo.RaspLogRequest) {
	host := &model.RaspHost{
		HostName:              req.HostName,
		AgentConfigUpdateTime: req.Ts,
	}
	err := l.RaspHostRepository.UpdateRaspHostByHostName(host)
	if err != nil {
		panic(err)
	}
}

func (l LogController) handleUpdateConfigId(req vo.RaspLogRequest) {
	// 获取configId
	detailMap := make(map[string]uint)
	err := json.Unmarshal([]byte(req.Detail), &detailMap)
	if err != nil {
		panic(err)
	}
	configId := detailMap["configId"]
	if err != nil {
		panic(err)
	}
	dbData, err := l.RaspHostRepository.QueryRaspHost(req.HostName)
	if err != nil {
		panic(err)
	}
	if len(dbData) > 0 {
		dbConfigId := dbData[0].ConfigId
		if dbConfigId != configId && dbConfigId > 0 {
			global.PushConfigQueue <- &vo.PushConfigRequest{
				ConfigId:  configId,
				HostNames: []string{req.HostName},
			}
		}
	} else {
		common.Log.Warnf("主机: %v 不存在, 无法推送消息", req.HostName)
	}
}

func (l LogController) handleFindJavaProcessLog(req vo.RaspLogRequest) {
	detailMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(req.Detail), &detailMap)
	if err != nil {
		panic(err)
		return
	}

	pid := detailMap["javaPid"].(float64)
	startTime := detailMap["startTime"].(string)
	message := detailMap["injectedStatus"].(string)
	status := l.convertMessageToStatus(message)
	cmdLines := detailMap["cmdLines"].([]interface{})
	appNames := detailMap["appNames"]
	var paramSlice []string
	for _, param := range cmdLines {
		switch v := param.(type) {
		case string:
			paramSlice = append(paramSlice, v)
		case int:
			strV := strconv.FormatInt(int64(v), 10)
			paramSlice = append(paramSlice, strV)
		default:
			panic("params type not supported")
		}
	}
	var appNamesStr = ""
	if appNames != nil {
		for _, item := range appNames.([]interface{}) {
			appNamesStr += item.(string) + ";"
		}
	}
	// 先判断表中是否有对应的pid
	processInfo, err := l.JavaProcessRepository.GetProcessByPid(req.HostName, uint(pid))
	if err != nil {
		panic(err)
	}
	// 如果库中没有则新增, 如果有则更新
	if processInfo == nil {
		process := &model.JavaProcessInfo{
			Status:       status,
			Message:      message,
			Pid:          int(pid),
			StartTime:    startTime,
			CmdlineInfo:  strings.Join(paramSlice, ","),
			AppNamesInfo: appNamesStr,
			HostName:     req.HostName,
		}
		err = l.JavaProcessRepository.SaveProcessInfo(process)
		if err != nil {
			panic(err)
		}
	} else {
		processInfo.Status = status
		processInfo.Message = message
		processInfo.AppNamesInfo = appNamesStr
		err = l.JavaProcessRepository.UpdateProcessInfo(processInfo)
		if err != nil {
			panic(err)
		}
	}
}

func (l LogController) handleRemoveJavaProcessLog(req vo.RaspLogRequest) {
	pid, err := strconv.ParseInt(req.Detail, 10, 32)
	if err != nil {
		panic(err)
	}
	err = l.JavaProcessRepository.DeleteProcessByPid(req.HostName, uint(pid))
	if err != nil {
		panic(err)
	}
}

func (l LogController) handleErrorLog(topic string, req vo.RaspLogRequest) {
	switch topic {
	case vo.JRASP_AGENT:
		l.handleAgentErrorLog(req)
	case vo.JRASP_DAEMON:
		l.handleDaemonErrorLog(req)
	case vo.JRASP_MODULE:
		l.handleModuleErrorLog(req)
	case vo.JRASP_ATTACK:
	default:
		panic(errors.New("unknown topic: " + req.Fields.KafkaTopic))
	}
}

func (l LogController) handleAgentErrorLog(req vo.RaspLogRequest) {
	Grok, _ := grok.New()
	maps, err := Grok.Parse(pattern, req.Message)
	if err != nil {
		// 不匹配的日志输出
		common.Log.Warnf(req.Message)
		return
	}
	if maps["level"] == "INFO" {
		return
	}
	errorLogs := &model.RaspErrorLogs{
		Topic:    vo.JRASP_AGENT,
		Time:     maps["time"],
		Level:    maps["level"],
		HostName: maps["host"],
		Message:  req.Message,
	}
	err = l.RaspErrorRepository.CreateRaspLogs(errorLogs)
	if err != nil {
		panic(err)
	}

}

func (l LogController) handleDaemonErrorLog(req vo.RaspLogRequest) {
	if req.Level == "INFO" {
		return
	}
	var message map[string]interface{}
	err := json.Unmarshal([]byte(req.Message), &message)
	if err != nil {
		panic(err)
	}
	errorLogs := &model.RaspErrorLogs{
		Topic:    vo.JRASP_DAEMON,
		Time:     req.Ts,
		Level:    req.Level,
		HostName: message["hostName"].(string),
		Message:  req.Message,
	}
	err = l.RaspErrorRepository.CreateRaspLogs(errorLogs)
	if err != nil {
		panic(err)
	}
}

func (l LogController) handleModuleErrorLog(req vo.RaspLogRequest) {
	Grok, _ := grok.New()
	maps, err := Grok.Parse(pattern, req.Message)
	if err != nil {
		// 不匹配的日志输出
		common.Log.Warnf(req.Message)
		return
	}
	if maps["level"] == "INFO" {
		return
	}
	errorLogs := &model.RaspErrorLogs{
		Topic:    vo.JRASP_MODULE,
		Time:     maps["time"],
		Level:    maps["level"],
		HostName: maps["host"],
		Message:  req.Message,
	}
	err = l.RaspErrorRepository.CreateRaspLogs(errorLogs)
	if err != nil {
		panic(err)
	}
}

func (l LogController) handleAttackLog(req vo.RaspLogRequest) {
	grok, _ := grok.New()
	maps, err := grok.Parse(pattern, req.Message)
	if err != nil {
		// 不匹配的日志输出
		common.Log.Warnf(req.Message)
		return
	}

	level := maps["level"]
	if level == "ERROR" || level == "ERR" {
		// 错误日志输出
		common.Log.Error(req.Message)
		return
	}

	attack := model.RaspAttack{
		HostName: maps["host"],
	}
	// 构件攻击详情对象
	detail := model.RaspAttackDetail{}

	// 攻击json
	msg := maps["message"]
	if msg != "" {
		var attackDetail = &vo.AttackDetail{}
		err := json.Unmarshal([]byte(msg), attackDetail)
		if err != nil {
			panic(err)
		}
		guid, _ := uuid.NewUUID()
		attack.RowGuid = guid.String()
		attack.HostIp = attackDetail.Context.LocalAddr
		attack.RemoteIp = attackDetail.Context.RemoteHost
		attack.RequestUri = attackDetail.Context.RequestURI
		attack.IsBlocked = attackDetail.IsBlocked
		attack.Level = attackDetail.Level
		attack.HandleResult = 0
		attack.AttackType = attackDetail.AttackType
		attack.AttackTime = time.Unix(attackDetail.AttackTime/1000, 0)

		// 构建攻击详情
		detail.ParentGuid = attack.RowGuid
		detail.Context = datatypes.JSON(util.Struct2Json(attackDetail.Context))
		detail.AppName = attackDetail.AppName
		detail.StackTrace = attackDetail.StackTrace
		detail.Payload = attackDetail.Payload
		detail.IsBlocked = attackDetail.IsBlocked
		detail.AttackType = attackDetail.AttackType
		detail.Algorithm = attackDetail.Algorithm
		detail.Extend = attackDetail.Extend
		detail.AttackTime = time.Unix(attackDetail.AttackTime/1000, 0)
		detail.Level = attackDetail.Level
		detail.MetaInfo = attackDetail.MetaInfo
	}

	err = l.RaspAttackRepository.CreateRaspAttack(&attack)
	if err != nil {
		return
	}
	err = l.RaspAttackRepository.CreateRaspAttackDetail(&detail)
	if err != nil {
		return
	}
}

func (l LogController) handleUpdateResourceName(req vo.RaspLogRequest) {
	detailMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(req.Detail), &detailMap)
	if err != nil {
		panic(err)
		return
	}
	resourceInfo, err := l.HostResourceRepository.GetResourceByNameAndIP(detailMap["hostName"].(string), detailMap["ip"].(string))
	if err != nil {
		panic(err)
	}
	if resourceInfo == nil {
		resourceInfo = &model.HostResource{
			HostName:     detailMap["hostName"].(string),
			Ip:           detailMap["ip"].(string),
			ResourceName: detailMap["resourceName"].(string),
		}
		err = l.HostResourceRepository.CreateResource(resourceInfo)
		if err != nil {
			panic(err)
		}
	}
}

func (l LogController) convertMessageToStatus(message string) int {
	status := 0
	if message == "success inject" || message == "success degrade" {
		status = 1
	} else if message == "not inject" {
		status = 0
	} else if message == "failed inject" || message == "failed uninstall agent" || message == "failed degrade" {
		status = 2
	}
	return status
}
