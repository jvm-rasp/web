package socket

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"server/common"
	"server/global"
	"server/vo"

	"server/model"
	"server/repository"
	"sync"
	"time"
)

// 注册识别的key
const REGISTER_CLIENT_KEY = "hostName"

// Manager 所有 websocket 信息
type Manager struct {
	Group                   map[string]map[string]*Client
	groupCount, clientCount uint
	Lock                    sync.Mutex
	Register, UnRegister    chan *Client
	Message                 chan *MessageData
	GroupMessage            chan *GroupMessageData
	BroadCastMessage        chan *BroadCastMessageData
}

// 启动 websocket 管理器
func (manager *Manager) Start() {
	common.Log.Info("websocket manage start")
	for {
		select {
		// client注册
		case client := <-manager.Register:
			common.Log.Infof("client [%s] register from group [%s]", client.Id, client.Group)
			manager.Lock.Lock()
			if manager.Group[client.Group] == nil {
				manager.Group[client.Group] = make(map[string]*Client)
				manager.groupCount += 1
			}
			manager.Group[client.Group][client.Id] = client
			manager.clientCount += 1
			manager.Lock.Unlock()

		// client注销
		case client := <-manager.UnRegister:
			common.Log.Infof("client [%s] unregister from group [%s]", client.Id, client.Group)
			manager.Lock.Lock()
			if _, ok := manager.Group[client.Group]; ok {
				if _, ok := manager.Group[client.Group][client.Id]; ok {
					close(client.Message)
					delete(manager.Group[client.Group], client.Id)
					manager.clientCount -= 1
					if len(manager.Group[client.Group]) == 0 {
						delete(manager.Group, client.Group)
						manager.groupCount -= 1
					}
				}
			}
			manager.Lock.Unlock()
		}
	}
}

// Client 单个 websocket 信息
type Client struct {
	Id /*ip*/, Group/*服务*/ string
	Socket  *websocket.Conn
	Message chan []byte
}

// messageData 单个发送数据信息
type MessageData struct {
	Id /*hostName*/, Group/*服务*/ string
	Message []byte
}

// groupMessageData 组广播数据信息
type GroupMessageData struct {
	Group   string
	Message []byte
}

// 广播发送数据信息
type BroadCastMessageData struct {
	Message []byte
}

// 读信息，从 websocket 连接直接读取数据
func (c *Client) Read() {
	defer func() {
		WebsocketManager.UnRegister <- c
		common.Log.Infof("client [%s] disconnect", c.Id)
		if err := c.Socket.Close(); err != nil {
			common.Log.Warnf("client [%s] disconnect err: %s", c.Id, err)
		}
	}()

	for {
		messageType, message, err := c.Socket.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			break
		}
		// c.Message <- message
		if messageType == websocket.TextMessage {
			hostName := c.Id
			// bugfix:  修复上报ip bug
			ip := c.Socket.RemoteAddr().(*net.TCPAddr).IP.String()
			heartbeatTime := string(message)
			hostInfo, err := RaspHostRepository.GetRaspHostByHostName(hostName)
			if err != nil {
				common.Log.Warnf("update host [%s] heartbeat err: %s", hostName, err)
			}
			if hostInfo == nil {
				// 如果是第一次连接则注册进库中
				hostInfo = &model.RaspHost{
					HostName:      hostName,
					Ip:            ip,
					HeartbeatTime: time.Now().Format("2006-01-02 15:04:05.000"),
				}
				configId, err := RaspHostRepository.CreateRaspHost(hostInfo)
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
				// 如果库中已有则更新
				hostInfo.HeartbeatTime = heartbeatTime
				err = RaspHostRepository.UpdateRaspHost(hostInfo)
				if err != nil {
					common.Log.Warnf("update host [%s] heartbeat err: %s", hostName, err)
				}
			}
		}
	}
}

// 给 client 发送消息
func (c *Client) Write() {
	defer func() {
		common.Log.Infof("client [%s] disconnect", c.Id)
		if err := c.Socket.Close(); err != nil {
			common.Log.Warnf("client [%s] disconnect err: %s", c.Id, err)
		}
	}()

	for {
		select {
		case message, ok := <-c.Message:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				common.Log.Warn("the chan of write message to remote is closed")
				return
			}
			err := c.Socket.WriteMessage(websocket.BinaryMessage, message)
			if err != nil {
				common.Log.Error("write message to remote client[%s] err: %s", c.Id, err)
			}
		}
	}
}

// 处理单个 client 发送数据
func (manager *Manager) SendService() {
	for {
		select {
		case data := <-manager.Message:
			if groupMap, ok := manager.Group[data.Group]; ok {
				if conn, ok := groupMap[data.Id]; ok {
					conn.Message <- data.Message
				}
			}
		}
	}
}

// 处理 group 广播数据
func (manager *Manager) SendGroupService() {
	for {
		select {
		// 发送广播数据到某个组的 channel 变量 Send 中
		case data := <-manager.GroupMessage:
			if groupMap, ok := manager.Group[data.Group]; ok {
				for _, conn := range groupMap {
					conn.Message <- data.Message
				}
			}
		}
	}
}

// 处理广播数据
func (manager *Manager) SendAllService() {
	for {
		select {
		case data := <-manager.BroadCastMessage:
			for _, v := range manager.Group {
				for _, conn := range v {
					conn.Message <- data.Message
				}
			}
		}
	}
}

// 向指定的 client 发送数据
func (manager *Manager) Send(id string, group string, message []byte) {
	data := &MessageData{
		Id:      id,
		Group:   group,
		Message: message,
	}
	manager.Message <- data
}

// 向指定的 Group 广播
func (manager *Manager) SendGroup(group string, message []byte) {
	data := &GroupMessageData{
		Group:   group,
		Message: message,
	}
	manager.GroupMessage <- data
}

// 广播
func (manager *Manager) SendAll(message []byte) {
	data := &BroadCastMessageData{
		Message: message,
	}
	manager.BroadCastMessage <- data
}

// 注册
func (manager *Manager) RegisterClient(client *Client) {
	manager.Register <- client
}

// 注销
func (manager *Manager) UnRegisterClient(client *Client) {
	manager.UnRegister <- client
}

// 当前组个数
func (manager *Manager) LenGroup() uint {
	return manager.groupCount
}

// 当前连接个数
func (manager *Manager) LenClient() uint {
	return manager.clientCount
}

// 获取 wsManager 管理器信息
func (manager *Manager) Info() map[string]interface{} {
	managerInfo := make(map[string]interface{})
	managerInfo["groupLen"] = manager.LenGroup()
	managerInfo["clientLen"] = manager.LenClient()
	managerInfo["chanRegisterLen"] = len(manager.Register)
	managerInfo["chanUnregisterLen"] = len(manager.UnRegister)
	managerInfo["chanMessageLen"] = len(manager.Message)
	managerInfo["chanGroupMessageLen"] = len(manager.GroupMessage)
	managerInfo["chanBroadCastMessageLen"] = len(manager.BroadCastMessage)
	return managerInfo
}

// 初始化 wsManager 管理器
var WebsocketManager = Manager{
	Group:            make(map[string]map[string]*Client),
	Register:         make(chan *Client, 128),
	UnRegister:       make(chan *Client, 128),
	GroupMessage:     make(chan *GroupMessageData, 128),
	Message:          make(chan *MessageData, 128),
	BroadCastMessage: make(chan *BroadCastMessageData, 128),
	groupCount:       0,
	clientCount:      0,
}

var RaspHostRepository = repository.NewRaspHostRepository()

// gin 处理 websocket handler
func (manager *Manager) WsClient(ctx *gin.Context) {
	upGrader := websocket.Upgrader{
		// cross origin domain
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		// 处理 Sec-WebSocket-Protocol Header
		Subprotocols: []string{ctx.GetHeader("Sec-WebSocket-Protocol")},
	}

	conn, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		common.Log.Error("websocket connect error: %s", ctx.Param(REGISTER_CLIENT_KEY))
		return
	}

	client := &Client{
		Id:      ctx.Param(REGISTER_CLIENT_KEY),
		Group:   ctx.Param(REGISTER_CLIENT_KEY),
		Socket:  conn,
		Message: make(chan []byte, 1024),
	}

	manager.RegisterClient(client)
	go client.Read()
	go client.Write()
}
