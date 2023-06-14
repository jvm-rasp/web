package report

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/fatih/structs"
	"github.com/go-playground/validator/v10"
	"github.com/gookit/goutil/fsutil"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	url2 "net/url"
	"os"
	"server/common"
	"server/config"
	"server/repository"
	"server/util"
	"strconv"
	"strings"
	"time"
)

const WS_URL_FORMAT = "%s://%s/log-collect/ws/%s/%s"

var UpdateManager *UpdateClient

type UpdateClient struct {
	SystemSettingRepository     repository.ISystemSettingRepository
	RaspHostRepository          repository.IRaspHostRepository
	JavaProcessInfoRepository   repository.IJavaProcessInfoRepository
	RaspFileRepository          repository.IRaspFileRepository
	RaspModuleRepository        repository.IRaspModuleRepository
	RaspComponentRepository     repository.IRaspComponentRepository
	RaspConfigRepository        repository.IRaspConfigRepository
	RaspConfigHistoryRepository repository.IRaspConfigHistoryRepository
	conn                        *websocket.Conn
	err                         error
	AutoUpdate                  bool
	isConnected                 bool
}

func InitUpdateManager() {
	// 判断bin目录下是否存在server.del文件，如果存在则删除
	exist := fsutil.PathExists(config.Conf.Env.BinFileName + ".del")
	if exist {
		err := os.Remove(config.Conf.Env.BinFileName + ".del")
		if err != nil {
			common.Log.Errorf(" delete server.del error: %v", err)
		}
	}

	SystemSettingRepository := repository.NewSystemSettingRepository()
	RaspHostRepository := repository.NewRaspHostRepository()
	JavaProcessInfoRepository := repository.NewJavaProcessInfoRepository(RaspHostRepository)
	RaspFileRepository := repository.NewRaspFileRepository()
	RaspModuleRepository := repository.NewRaspModuleRepository()
	RaspComponentRepository := repository.NewRaspComponentRepository()
	RaspConfigRepository := repository.NewRaspConfigRepository()
	RaspConfigHistoryRepository := repository.NewRaspConfigHistoryRepository()
	UpdateManager = &UpdateClient{
		isConnected:                 false,
		AutoUpdate:                  false,
		SystemSettingRepository:     SystemSettingRepository,
		RaspHostRepository:          RaspHostRepository,
		JavaProcessInfoRepository:   JavaProcessInfoRepository,
		RaspFileRepository:          RaspFileRepository,
		RaspModuleRepository:        RaspModuleRepository,
		RaspComponentRepository:     RaspComponentRepository,
		RaspConfigRepository:        RaspConfigRepository,
		RaspConfigHistoryRepository: RaspConfigHistoryRepository,
	}
}

func (this *UpdateClient) Start() {
	autoUpdate, err := this.SystemSettingRepository.GetSettingByName("autoUpdate")
	if err != nil {
		common.Log.Errorf("获取系统配置%v失败, error: %v", "autoUpdate", err)
		return
	}
	this.AutoUpdate, _ = strconv.ParseBool(autoUpdate.Value)
	for {
		if this.AutoUpdate && this.isConnected == false {
			this.Connect()
		}
		time.Sleep(time.Second * 10)
	}
}

func (this *UpdateClient) Connect() {
	updateUrl, err := this.SystemSettingRepository.GetSettingByName("updateUrl")
	if err != nil {
		common.Log.Errorf("获取系统配置%v失败, error: %v", "autoUpdate", err)
		return
	}
	url, err := url2.Parse(updateUrl.Value)
	if err != nil {
		common.Log.Errorf("解析服务器地址%v失败, error: %v", updateUrl.Value, err)
		return
	}
	projectGuid, err := this.SystemSettingRepository.GetSettingByName("projectGuid")
	if err != nil {
		common.Log.Errorf("获取系统配置%v失败, error: %v", "projectGuid", err)
		return
	}
	Url := fmt.Sprintf(WS_URL_FORMAT, url.Scheme, url.Host, projectGuid.Value, config.Conf.Env.HostName)
	// 判断协议是否加密
	var tlsConfig tls.Config
	if strings.HasPrefix(Url, "wss://") {
		tlsConfig = tls.Config{InsecureSkipVerify: true}
	}
	dialer := &websocket.Dialer{
		TLSClientConfig: &tlsConfig,
	}
	this.conn, _, this.err = dialer.Dial(Url, nil)
	defer this.DisConnect()
	if this.err != nil {
		common.Log.Errorf("连接远程更新服务端: %s 错误, error: %v", Url, this.err)
	} else {
		common.Log.Infof("连接远程更新服务端: %s 成功", Url)
		this.isConnected = true
		// 连接上远程服务端后立即更新client信息
		updateRequest, err := this.generateUpdateRequest()
		if err != nil {
			common.Log.Errorf("生成更新服务端数据失败, error: %v", err)
		} else {
			wsMessageRequest := &WebSocketMessageRequest{
				MessageType:    CLIENT_UPDATE,
				MessageContent: structs.Map(updateRequest),
			}
			err = this.conn.WriteJSON(wsMessageRequest)
			if err != nil {
				common.Log.Errorf("发送json消息失败, error: %v", err)
			}
		}
		// 开启心跳线程
		go this.Heartbeat()
		// 读取服务端发过来的信息
		for {
			messageType, message, err := this.conn.ReadMessage()
			if err != nil {
				common.Log.Errorf("读取远程更新服务端信息出错: %v", err)
				break
			}
			if messageType == websocket.TextMessage {
				common.Log.Infof("读取远程更新服务端信息成功: \n%v", string(message))
				var webSocketMessage WebSocketMessageRequest
				err = json.Unmarshal(message, &webSocketMessage)
				if err != nil {
					common.Log.Error("反序列化消息失败, %v", err)
					continue
				}
				// 参数校验
				if err = common.Validate.Struct(&webSocketMessage); err != nil {
					errStr := err.(validator.ValidationErrors)[0].Translate(common.Trans)
					if err != nil {
						common.Log.Error("校验数据格式出错 err: %s", errStr)
					}
					continue
				}
				switch webSocketMessage.MessageType {
				case CLIENT_UPGRADE:
					var messageContent ClientUpgradeRequest
					err = mapstructure.Decode(webSocketMessage.MessageContent, &messageContent)
					if err != nil {
						common.Log.Error("反序列化消息失败, %v", err)
						break
					}
					this.handleClientUpgrade(messageContent)
				}
			}

		}
	}
}

func (this *UpdateClient) Heartbeat() {
	for {
		if this.conn == nil || this.AutoUpdate == false || this.isConnected == false {
			break
		}
		heartBeatRequest, err := this.generateHeartbeat()
		if err != nil {
			common.Log.Errorf("生成心跳数据失败, error: %v", err)
		} else {
			heartMessage := WebSocketMessageRequest{
				MessageType:    CLIENT_HEARTBEAT,
				MessageContent: structs.Map(heartBeatRequest),
			}
			err = this.conn.WriteJSON(heartMessage)
			if err != nil {
				common.Log.Errorf("发送心跳数据失败, error: %v", err)
				break
			}
		}
		time.Sleep(time.Second * 180)
	}
}

func (this *UpdateClient) generateUpdateRequest() (*ServerUpdateRequest, error) {
	projectType, err := this.SystemSettingRepository.GetSettingByName("projectType")
	if err != nil {
		common.Log.Errorf("获取系统配置%v失败, error: %v", "projectType", err)
		return nil, err
	}
	updateRequest := &ServerUpdateRequest{
		GroupType:   projectType.Value,
		Ip:          config.Conf.Env.Ip,
		ExeFileHash: config.Conf.Env.BinFileHash,
		OsType:      config.Conf.Env.OsType,
	}
	return updateRequest, nil
}

func (this *UpdateClient) generateHeartbeat() (*ServerHeartRequest, error) {
	// 生成时间
	location := time.FixedZone("CST", 8*3600)
	t := time.Now().In(location)
	now := t.Format("2006-01-02 15:04:05.000")
	// 构造host list
	var agentInfoList = []AgentInfo{}
	hostList, _, err := this.RaspHostRepository.GetRaspHostList()
	if err != nil {
		common.Log.Errorf("获取Agent列表失败, error: %v", err)
		return nil, err
	}
	for _, item := range hostList {
		var agentInfo AgentInfo
		agentInfo.HostName = item.HostName
		agentInfo.Ip = item.Ip
		agentInfo.AgentMode = item.AgentMode
		agentInfo.CreateTime = item.CreatedAt.Format("2006-01-02 15:04:05.000")
		agentInfo.HeartbeatTime = item.HeartbeatTime
		agentInfo.Version = item.Version
		agentInfo.ExeFileHash = item.ExeFileHash
		agentInfo.OsType = item.OsType
		agentInfo.SuccessInject = item.SuccessInject
		agentInfo.FailedInject = item.FailedInject
		agentInfo.NotInject = item.NotInject
		var processInfo = []JavaProcessInfo{}
		procList, err := this.JavaProcessInfoRepository.GetAllJavaProcessInfos(item.HostName)
		if err != nil {
			common.Log.Errorf("获取java进程失败, error: %v", err)
			return nil, err
		}
		for _, proc := range procList {
			var process JavaProcessInfo
			util.Struct2Struct(proc, &process)
			processInfo = append(processInfo, process)
		}
		agentInfo.ProcessInfo = processInfo
		agentInfoList = append(agentInfoList, agentInfo)
	}

	serverInfo := ServerInfo{
		HostName: config.Conf.Env.HostName,
		Ip:       config.Conf.Env.Ip,
		Time:     now,
		List:     agentInfoList,
	}

	heartRequest := &ServerHeartRequest{
		HeartBeatTime: now,
		LogName:       "heart",
		LogType:       "rasp",
		Version:       "v2",
		LogContent:    serverInfo,
	}
	return heartRequest, nil
}

func (this *UpdateClient) DisConnect() {
	if this.conn != nil {
		err := this.conn.Close()
		if err != nil {
			common.Log.Errorf("关闭服务端连接失败, error: %v", err)
		}
		this.conn = nil
	}
	this.isConnected = false
}
