package common

import (
	"fmt"
	"github.com/hashicorp/mdns"
	"net"
	"os"
	"server/config"
	"server/util"
)

type MDNSServer struct {
	c       chan os.Signal
	server  *mdns.Server
	service *mdns.MDNSService
}

func InitMDNSService() {
	if config.Conf.Mdns.Enable {
		server := NewMDNSServer()
		server.NewService()
		go server.Start()
		Log.Infof("初始化mdns服务完成")
	}
}

func NewMDNSServer() *MDNSServer {
	return new(MDNSServer)
}

func (s *MDNSServer) NewService() {
	// Setup our service export
	instance := "admin"
	port := config.Conf.System.Port
	hostName, _ := os.Hostname()
	domain := ""
	serviceName := "jrasp"

	protocol := util.Ternary(config.Conf.Ssl.Enable, "wss", "ws")
	ip := util.GetDefaultIp()
	prefix := config.Conf.System.UrlPathPrefix
	txt := []string{fmt.Sprintf("%v://%v:%v/%v", protocol, ip, port, prefix)}
	ips := []net.IP{
		net.ParseIP(ip),
	}
	var err error
	s.service, err = mdns.NewMDNSService(instance, serviceName, domain, hostName+".", port, ips, txt)
	if err != nil {
		panic(err)
	}
}

func (s *MDNSServer) Start() {
	iface, err := util.GetDefaultIface()
	if err != nil {
		panic(err)
	}
	dnsConfig := &mdns.Config{
		Zone:              s.service,
		Iface:             iface,
		LogEmptyResponses: false,
	}
	s.server, err = mdns.NewServer(dnsConfig)
	if err != nil {
		panic(err)
	}
	defer s.ShutDown()
	select {
	case sig := <-s.c:
		Log.Debugf("Got %s signal. Aborting...\n", sig)
		break
	}
}

func (s *MDNSServer) ShutDown() {
	err := s.server.Shutdown()
	if err != nil {
		panic(err)
	}
	close(s.c)
}
