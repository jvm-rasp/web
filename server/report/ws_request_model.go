package report

const (
	CLIENT_UPDATE         int = 1
	CLIENT_HEARTBEAT      int = 2
	CLIENT_UPGRADE        int = 3
	CLIENT_UPGRADE_RESULT int = 4
)

type WebSocketMessageRequest struct {
	MessageType    int                    `json:"messageType"`
	MessageContent map[string]interface{} `json:"messageContent"`
}

type ServerUpdateRequest struct {
	GroupType   string `json:"groupType"`
	Ip          string `json:"ip"`
	ExeFileHash string `json:"exeFileHash"`
	OsType      string `json:"osType"`
}

type ServerHeartRequest struct {
	HeartBeatTime string     `json:"heartBeatTime"`
	LogName       string     `json:"logName"`
	LogType       string     `json:"logType"`
	Version       string     `json:"version"`
	LogContent    ServerInfo `json:"logContent"`
}

type ServerInfo struct {
	HostName string      `json:"hostname" structs:"hostname"`
	Ip       string      `json:"ip" structs:"ip"`
	Time     string      `json:"time" structs:"time"`
	List     []AgentInfo `json:"list" structs:"list"`
}

type AgentInfo struct {
	HostName      string            `json:"hostName"  structs:"hostName"`
	Ip            string            `json:"ip"  structs:"ip"`
	AgentMode     string            `json:"agentMode"  structs:"agentMode"`
	CreateTime    string            `json:"createTime"  structs:"createTime"`
	HeartbeatTime string            `json:"heartbeatTime"  structs:"heartbeatTime"`
	Version       string            `json:"version"  structs:"version"`
	ExeFileHash   string            `json:"exeFileHash"  structs:"exeFileHash"`
	OsType        string            `json:"osType"  structs:"osType"`
	SuccessInject int64             `json:"successInject"  structs:"successInject"`
	FailedInject  int64             `json:"failedInject"  structs:"failedInject"`
	NotInject     int64             `json:"notInject"  structs:"notInject"`
	ProcessInfo   []JavaProcessInfo `json:"processInfo"  structs:"processInfo"`
}

type JavaProcessInfo struct {
	CmdlineInfo  string `json:"cmdlineInfo" structs:"cmdlineInfo"`
	AppNamesInfo string `json:"appNamesInfo" structs:"appNamesInfo"`
	StartTime    string `json:"startTime" structs:"startTime"`
	Pid          int    `json:"pid" structs:"pid"`
	Status       int    `json:"status" structs:"status"`
	Message      string `json:"message" structs:"message"`
}

type ClientUpgradeRequest struct {
	BatchGuid   string `json:"batchGuid" structs:"batchGuid"`
	ProjectGuid string `json:"projectGuid" structs:"projectGuid"`
	HostName    string `json:"hostName" structs:"hostName"`
	DownloadUrl string `json:"downloadUrl" structs:"downloadUrl"`
	Md5         string `json:"md5" structs:"md5"`
	Type        string `json:"type" structs:"type"`
	State       int    `json:"state" structs:"state"`
	Message     string `json:"message" structs:"message"`
}
