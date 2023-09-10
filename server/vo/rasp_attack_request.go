package vo

type RaspAttackListRequest struct {
	HostName     string `json:"hostName" form:"hostName"`
	Url          string `json:"url" form:"url"`
	IsBlocked    string `json:"isBlocked" form:"isBlocked"`
	HandleResult string `json:"handleResult" form:"handleResult"`
	PageNum      uint   `json:"pageNum" form:"pageNum"`
	PageSize     uint   `json:"pageSize" form:"pageSize"`

	AttackStartTime string `json:"attackStartTime" form:"attackStartTime"`
	AttackEndTime   string `json:"attackEndTime" form:"attackEndTime"`
	HostIp          string `json:"hostIp" form:"hostIp"`
	AttackType      string `json:"attackType" form:"attackType"`
	RemoteIp        string `json:"remoteIp" form:"remoteIp"`
	Level           string `json:"level" form:"level"`
}

type RaspAttackDetailRequest struct {
	Id string `json:"id" form:"id"`
}

type DeleteRaspAttackRequest struct {
	Guids []string `json:"guids" form:"guids"`
}

type UpdateRaspStatusRequest struct {
	Id     uint `json:"id" form:"id"`
	Result int  `json:"result" form:"result"`
}

type AttackDetail struct {
	Context    Context `json:"context"`
	AppName    string  `json:"appName"`
	StackTrace string  `json:"stackTrace"`
	Payload    string  `json:"payload"`
	IsBlocked  bool    `json:"isBlocked"`
	AttackType string  `json:"attackType"`
	Algorithm  string  `json:"algorithm"`
	Extend     string  `json:"extend"`
	AttackTime int64   `json:"attackTime"`
	Level      int     `json:"level"`
	MetaInfo   string  `json:"metaInfo"`
}

type Context struct {
	Method            string `json:"method"`
	Protocol          string `json:"protocol"`
	LocalAddr         string `json:"localAddr"`
	RemoteHost        string `json:"remoteHost"`
	RequestURL        string `json:"requestURL"`
	RequestURI        string `json:"requestURI"`
	ContentType       string `json:"contentType"`
	ContentLength     int    `json:"contentLength"`
	CharacterEncoding string `json:"characterEncoding"`
	Parameters        string `json:"parameters"`
	Header            string `json:"header"`
	QueryString       string `json:"queryString"`
	Marks             string `json:"marks"`
	Body              string `json:"body"`
}
