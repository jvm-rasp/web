package vo

type RaspHostListRequest struct {
	HostName  string `json:"name" form:"hostName"`
	Ip        string `json:"ip" form:"ip"`
	AgentMode string `json:"agentMode" form:"agentMode"`
	Status    uint   `json:"status" form:"status"`
	PageNum   uint   `json:"pageNum" form:"pageNum"`
	PageSize  uint   `json:"pageSize" form:"pageSize"`
}

type DeleteRaspHostRequest struct {
	Ids []uint `json:"ids" form:"ids" validate:"required"`
}

type PushConfigRequest struct {
	// 配置id
	ConfigId uint   `json:"configId" form:"configId" validate:"required"`
	// 主机名称列表
	HostNames      []string `json:"hostNames" form:"hostNames" validate:"required"`
}
