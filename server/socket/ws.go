package socket

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"server/common"
	"sync"
)

// 注册识别的key
const REGISTER_CLIENT_KEY = "hostName"

// messageData 单个发送数据信息
type MessageData struct {
	Id /*hostName*/, Group/*服务*/ string
	Message []byte
}

// Manager 所有 websocket 信息
type Manager struct {
	Group                map[string]map[string]*Client
	Lock                 sync.Mutex
	Register, UnRegister chan *Client
	Message              chan *MessageData
	LogChan              chan string // 日志通道
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
			}
			manager.Group[client.Group][client.Id] = client
			manager.Lock.Unlock()

		// client注销
		case client := <-manager.UnRegister:
			common.Log.Infof("client [%s] unregister from group [%s]", client.Id, client.Group)
			manager.Lock.Lock()
			if _, ok := manager.Group[client.Group]; ok {
				if _, ok := manager.Group[client.Group][client.Id]; ok {
					close(client.Message)
					delete(manager.Group[client.Group], client.Id)
					if len(manager.Group[client.Group]) == 0 {
						delete(manager.Group, client.Group)
					}
				}
			}
			manager.Lock.Unlock()
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

// 向指定的 client 发送数据
func (manager *Manager) Send(id string, group string, message []byte) {
	data := &MessageData{
		Id:      id,
		Group:   group,
		Message: message,
	}
	manager.Message <- data
}

// 注册
func (manager *Manager) RegisterClient(client *Client) {
	manager.Register <- client
}

// 注销
func (manager *Manager) UnRegisterClient(client *Client) {
	manager.UnRegister <- client
}

// 初始化 wsManager 管理器
var WebsocketManager = Manager{
	Group:      make(map[string]map[string]*Client),
	Register:   make(chan *Client, 128),
	UnRegister: make(chan *Client, 128),
	Message:    make(chan *MessageData, 128),
	LogChan:    common.LogChan,
}

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
		common.Log.Errorf("websocket connect error: %s", ctx.Param(REGISTER_CLIENT_KEY))
		return
	}

	client := &Client{
		Id:      ctx.Param(REGISTER_CLIENT_KEY),
		Group:   ctx.Param(REGISTER_CLIENT_KEY),
		Socket:  conn,
		Message: make(chan []byte, 1024),
		LogChan: manager.LogChan,
	}

	manager.RegisterClient(client)
	go client.Read()
	go client.Write()
}
