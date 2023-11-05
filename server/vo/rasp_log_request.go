package vo

import "time"

// topic 名称
const (
	JRASP_DAEMON = "jrasp-daemon"
	JRASP_AGENT  = "jrasp-agent"
	JRASP_MODULE = "jrasp-module"
	JRASP_ATTACK = "jrasp-attack"
)

type RaspLogRequest struct {
	Timestamp time.Time `json:"@timestamp"`
	Fields    Fields    `json:"fields"`
	Host      Host      `json:"host"`
	Message   string    `json:"message"`

	Caller   string `json:"caller"`
	Detail   string `json:"detail"`
	HostName string `json:"hostName"`
	Ip       string `json:"ip"`
	Level    string `json:"level"`
	LogId    int    `json:"logId"`
	Msg      string `json:"msg"`
	Pid      int    `json:"pid"`
	Ts       string `json:"ts"`
}

type Host struct {
	Name string `json:"name"`
}

type Fields struct {
	KafkaTopic string `json:"kafka_topic"`
}
