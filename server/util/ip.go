package util

import (
	"net"
	"strings"
)

func GetDefaultIp() string {
	conn, err := net.Dial("udp", "114.114.114.114:53")
	if err != nil {
		return ""
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip := strings.Split(localAddr.IP.String(), ":")[0]
	return ip
}
