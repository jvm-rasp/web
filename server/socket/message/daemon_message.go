package message

// daemon 1000~1999
const (
	DAEMON_STARTUP_LOGID  = 1000
	HOST_ENV_LOGID        = 1002
	HEART_BEAT_LOGID      = 1011
	DEPENDENCY_JAR_LOGID  = 1016
	AGENT_SUCCESS_UNLOAD  = 1017
	JAVA_PROCESS_STARTUP  = 1018
	JAVA_PROCESS_SHUTDOWN = 1019
	AGENT_SUCCESS_INIT    = 1020
	NACOS_INIT_INFO       = 1024
	Agent_CONFIG_UPDATE   = 1025
	CONFIG_ID             = 1030
)

type DaemonMessage struct {
	Level    string `json:"level"`
	Ts       string `json:"ts"`
	Caller   string `json:"caller"`
	Msg      string `json:"msg"`
	Logid    int    `json:"logId"`
	Ip       string `json:"ip"`
	HostName string `json:"hostName"`
	Pid      int    `json:"pid"`
	Detail   string `json:"detail"`
}
