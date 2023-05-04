package util

import (
	"errors"
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

func GetDefaultIface() (*net.Interface, error) {
	defaultIP := GetDefaultIp()
	if defaultIP == "" {
		return nil, errors.New("not found default ip")
	}
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, item := range ifaces {
		addrs, err := item.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := addr.(*net.IPNet)
			if ip.IP.String() == defaultIP {
				return &item, nil
			}
		}
	}
	return nil, errors.New("not found default ifaces")
}
