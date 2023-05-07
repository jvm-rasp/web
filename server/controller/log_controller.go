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
const pattern = "%{TIMESTAMP_ISO8601:time}\\s*%{LOGLEVEL:level}\\s*\\[%{DATA:thread}\\]\\s*\\[%{DATA:api}\\]\\s*%{GREEDYDATA:message}"

type ILogController interface {
	ReportLog(c *gin.Context)
	GetAttackLogs(c *gin.Context)
	GetAttackDetail(c *gin.Context)
	BatchDeleteLogByIds(c *gin.Context)
	UpdateStatusById(c *gin.Context)
}

type LogController struct {
	RaspHostRepository    repository.IRaspHostRepository
	JavaProcessRepository repository.IJavaProcessInfoRepository
	RaspAttackRepository  repository.IRaspAttackRepository
	RaspErrorRepository   repository.IRaspErrorLogsRepository
	RaspConfigRepository  repository.IRaspConfigRepository
}

func NewLogController() ILogController {
	repository1 := repository.NewRaspHostRepository()
	repository2 := repository.NewJavaProcessInfoRepository()
	repository3 := repository.NewRaspAttackRepository()
	repository4 := repository.NewRaspErrorLogsRepository()
	repository5 := repository.NewRaspConfigRepository()
	controller := LogController{
		RaspHostRepository:    repository1,
		JavaProcessRepository: repository2,
		RaspAttackRepository:  repository3,
		RaspErrorRepository:   repository4,
		RaspConfigRepository:  repository5,
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
		l.handleAgentErrorLog(req)
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
		configId, err := l.RaspHostRepository.CreateRaspHost(host)
		if err != nil {
			panic(err)
		}
		// 推送默认配置
		if configId != 0 {
			hostController := NewRaspHostController()
			content, err := hostController.GeneratePushConfig(configId)
			if err != nil {
				panic(err)
			}
			hostController.PushHostsConfig([]string{host.HostName}, content)
		}
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
	host.HeartbeatTime = req.Ts

	if len(dbData) == 0 {
		configId, err := l.RaspHostRepository.CreateRaspHost(host)
		if err != nil {
			panic(err)
		}
		// 推送默认配置
		if configId != 0 {
			hostController := NewRaspHostController()
			content, err := hostController.GeneratePushConfig(configId)
			if err != nil {
				panic(err)
			}
			hostController.PushHostsConfig([]string{host.HostName}, content)
		}
	} else {
		err := l.RaspHostRepository.UpdateRaspHostByHostName(host)
		if err != nil {
			panic(err)
		}
	}
}

func (l LogController) handleHeartbeatLog(req vo.RaspLogRequest) {
	dbData, err := l.RaspHostRepository.QueryRaspHost(req.HostName)
	if err != nil {
		panic(err)
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

	host.HeartbeatTime = req.Ts

	if len(dbData) == 0 {
		configId, err := l.RaspHostRepository.CreateRaspHost(host)
		if err != nil {
			panic(err)
		}
		// 推送默认配置
		if configId != 0 {
			hostController := NewRaspHostController()
			content, err := hostController.GeneratePushConfig(configId)
			if err != nil {
				panic(err)
			}
			hostController.PushHostsConfig([]string{host.HostName}, content)
		}
	} else {
		err := l.RaspHostRepository.UpdateRaspHostByHostName(host)
		if err != nil {
			panic(err)
		}
	}

	// db 数据
	dblist, err := l.JavaProcessRepository.GetAllJavaProcessInfos(req.HostName)
	if err != nil {
		panic(err)
	}

	// 上报数据
	detailMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(req.Detail), &detailMap)
	if err != nil {
		panic(err)
	}

	for _, v := range dblist {
		// db有，上报数据中没有，删除db数据
		if len(detailMap) == 0 || detailMap[string(v.Pid)] == nil {
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
		for _, v2 := range dblist {
			if string(rune(v2.Pid)) == k {
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
		status := processDetail["status"].(string)
		// 缺少命令后信息
		process := &model.JavaProcessInfo{HostName: req.HostName, Pid: int(pid), StartTime: startTime, Status: status}
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
	status := detailMap["status"].(string)

	process := &model.JavaProcessInfo{
		Status:    status,
		Pid:       int(pid),
		StartTime: startTime,
		HostName:  req.HostName,
	}
	err = l.JavaProcessRepository.UpdateProcessByHostName(process)
	if err != nil {
		panic(err)
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
		if dbConfigId != configId {
			hostController := NewRaspHostController()
			content, err := hostController.GeneratePushConfig(dbConfigId)
			if err != nil {
				panic(err)
			}
			hostController.PushHostsConfig([]string{req.HostName}, content)
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
	status := detailMap["injectedStatus"].(string)

	cmdLines := detailMap["cmdLines"].([]interface{})
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

	process := &model.JavaProcessInfo{
		Status:      status,
		Pid:         int(pid),
		StartTime:   startTime,
		CmdlineInfo: strings.Join(paramSlice, ","),
		HostName:    req.HostName,
	}
	err = l.JavaProcessRepository.SaveProcessInfo(process)
	if err != nil {
		panic(err)
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
		Topic:   vo.JRASP_AGENT,
		Time:    maps["time"],
		Level:   maps["level"],
		Message: req.Message,
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
	errorLogs := &model.RaspErrorLogs{
		Topic:   vo.JRASP_DAEMON,
		Time:    req.Ts,
		Level:   req.Level,
		Message: req.Message,
	}
	err := l.RaspErrorRepository.CreateRaspLogs(errorLogs)
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
		Topic:   vo.JRASP_MODULE,
		Time:    maps["time"],
		Level:   maps["level"],
		Message: req.Message,
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

	// 构建对象
	attack := model.RaspAttack{
		HostName: req.Host.Name,
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
		detail.StackTrace = attackDetail.StackTrace
		detail.Payload = attackDetail.Payload
		detail.IsBlocked = attackDetail.IsBlocked
		detail.AttackType = attackDetail.AttackType
		detail.Algorithm = attackDetail.Algorithm
		detail.Extend = attackDetail.Extend
		detail.AttackTime = time.Unix(attackDetail.AttackTime/1000, 0)
		detail.Level = attackDetail.Level
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
