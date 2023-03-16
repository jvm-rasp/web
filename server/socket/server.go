package socket

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"server/common"
)

// 日志topic 或者命令
const DAEMON_TOPIC byte = 0x01
const AGENT_TOPIC byte = 0x02
const MODULE_TOPIC byte = 0x03
const ATTACK_TOPIC byte = 0x04

// LOG_PACKAGE_CONSTANT_LENGTH 消息体的固定长度 149
const LOG_PACKAGE_CONSTANT_LENGTH = 149

func NewSockekServer(port int) error {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%v", port))
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}
	defer tcpListener.Close()
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handle(tcpConn)
	}
	return nil
}

func handle(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	// 长度分割
	scanner.Split(splitLog)
	for scanner.Scan() {
		scannedPack := new(LogPackage)
		scannedPack.Unpack(bytes.NewReader(scanner.Bytes()))
		Handler(scannedPack)
	}
}

// 处理日志
func Handler(p *LogPackage) {
	switch p.Type {
	case DAEMON_TOPIC:
		common.Log.Debug(fmt.Sprintf("handler log hostname: %s", p.HostName))
		handleDaemonLog(p)
	default:

	}
}

// splitLog 将字节码流分割为消息体
func splitLog(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if !atEOF && string(data[0:3]) == string(magicBytes[:]) {
		if len(data) > 13 {
			hostNamelength := int32(0)
			err := binary.Read(bytes.NewReader(data[5:9]), binary.BigEndian, &hostNamelength)
			if err != nil {
				return 0, nil, err
			}
			bodySize := int32(0)
			err = binary.Read(bytes.NewReader(data[9:13]), binary.BigEndian, &bodySize)
			if err != nil {
				return 0, nil, err
			}
			length := int(hostNamelength) + int(bodySize)
			if length+LOG_PACKAGE_CONSTANT_LENGTH <= len(data) {
				return LOG_PACKAGE_CONSTANT_LENGTH + length, data[:LOG_PACKAGE_CONSTANT_LENGTH+length], nil
			}
		}
	}
	return
}

func handleDaemonLog(p *LogPackage) {
	//　TODO
}
