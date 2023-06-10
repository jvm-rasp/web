package socket

import (
	"bytes"
	"fmt"
	"net"
	"server/common"
	"server/config"
	"server/util"
	"strings"
	"time"
)

type UDPServer struct {
	conn    *net.UDPConn
	udpAddr *net.UDPAddr
	err     error
}

func NewUDPServer() *UDPServer {
	udpAddr, err := net.ResolveUDPAddr("udp4", "0.0.0.0:8080")
	if err != nil {
		common.Log.Errorf("resolve udp address error %v", err.Error())
		return nil
	}
	return &UDPServer{
		udpAddr: udpAddr,
	}
}

func (c *UDPServer) Start() {
	c.conn, c.err = net.ListenUDP("udp", c.udpAddr)
	if c.err != nil {
		common.Log.Errorf("bind udp port error: %v", c.err)
		return
	}
	for {
		buf := make([]byte, 2048)
		length, addr, err := c.conn.ReadFromUDP(buf[:]) // 读取数据，返回值依次为读取数据长度、远端地址、错误信息 // 读取操作会阻塞直至有数据可读取
		if err != nil {
			if strings.Index(err.Error(), "use of closed network connection") < 0 {
				common.Log.Errorf("read udp data error: %v", err.Error())
			}
			continue
		}
		data := buf[:length]
		scannedPack := new(Package)
		err = scannedPack.Unpack(bytes.NewReader(data))
		if err != nil {
			common.Log.Errorf("unpack udp data error: %v", err.Error())
			continue
		}
		if scannedPack.Type == SEARCH_SERVER {
			receivedMessage := string(scannedPack.Body)
			common.Log.Infof("received udp data from remote: %v, message: %v", addr, receivedMessage)
			dstAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%v:%v", addr.IP, 8080))
			if err != nil {
				continue
			}
			conn, err := net.DialUDP("udp", nil, dstAddr)
			sendMessage := fmt.Sprintf("%v://%v:%v/%v",
				util.Ternary(config.Conf.Ssl.Enable, "wss", "ws"),
				util.GetDefaultIp(),
				config.Conf.System.Port,
				config.Conf.System.UrlPathPrefix)
			pack := &Package{
				Magic:     MagicBytes,
				Version:   PROTOCOL_VERSION,
				Type:      UPDATE_SERVER,
				BodySize:  int32(len(sendMessage)),
				TimeStamp: time.Now().Unix(),
				Signature: EmptySignature,
				Body:      []byte((sendMessage)),
			}
			sendBuf := bytes.NewBuffer(nil)
			err = pack.Pack(sendBuf)
			if err != nil {
				continue
			}
			_, err = conn.Write(sendBuf.Bytes())
			if err != nil {
				continue
			}
		}
	}
}
