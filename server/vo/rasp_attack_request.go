package vo

type RaspAttackListRequest struct {
	HostName string `json:"hostName" form:"hostName"`
	PageNum  uint   `json:"pageNum" form:"pageNum"`
	PageSize uint   `json:"pageSize" form:"pageSize"`
}

type RaspAttackDetailRequest struct {
	Id string `json:"id" form:"id"`
}

type DeleteRaspAttackRequest struct {
	Ids []uint `json:"ids" form:"ids"`
}

type AttackDetail struct {
	Context    Context `json:"context"`
	StackTrace string  `json:"stackTrace"`
	Payload    string  `json:"payload"`
	IsBlocked  bool    `json:"isBlocked"`
	AttackType string  `json:"attackType"`
	Algorithm  string  `json:"algorithm"`
	Extend     string  `json:"extend"`
	AttackTime int64   `json:"attackTime"`
	Level      int     `json:"level"`
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
