package vo

type RaspLogRecord struct {
	Level    string `json:"level"`
	Ts       int64  `json:"ts"`
	Caller   string `json:"caller"`
	Msg      string `json:"msg"`
	LogId    int    `json:"logId"`
	Ip       string `json:"ip"`
	HostName string `json:"hostName"`
	Detail   string `json:"detail"`
	Pid      int    `json:"pid"`
	Topic    string `json:"topic"`
	// agent 日志
	ProcessId  string `json:"processId"`
	Thread     string `json:"thread"`
	StackTrace string `json:"stackTrace"`
}
