package common

import (
	"encoding/json"
	"server/config"
	"time"
)

func InitHeartBeat() {
	heartBeat := HeartBeat{
		HostName: config.Conf.Env.HostName,
		Ip:       config.Conf.Env.Ip,
	}
	go heartBeat.Start()
}

type HeartBeat struct {
	HostName string `json:"hostname"`
	Ip       string `json:"ip"`
	Time     string `json:"time"`
}

func (this *HeartBeat) Start() {
	for {
		this.Heart()
		time.Sleep(time.Second * 30)
	}
}

func (this *HeartBeat) Heart() {
	this.Time = time.Now().Format("2006-01-02 15:04:05")
	content, err := json.Marshal(this)
	if err != nil {
		Log.Errorf("序列化心跳包失败, err: %v", err)
		return
	}
	Log.Debug(string(content))
}
