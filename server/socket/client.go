package socket

import (
	"github.com/gorilla/websocket"
	"server/common"
)

// Client 单个 websocket 信息
type Client struct {
	Id /*ip*/, Group/*服务*/ string
	Socket  *websocket.Conn
	Message chan []byte // 命令通道
	LogChan chan string // 日志通道
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
		c.LogChan <- string(message)
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
				common.Log.Warnf("the chan of write message to remote is closed")
				return
			}
			err := c.Socket.WriteMessage(websocket.BinaryMessage, message)
			if err != nil {
				common.Log.Errorf("write message to remote client[%s] err: %s", c.Id, err)
			}
		}
	}
}
