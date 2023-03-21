package vo

// CreateRaspConfigRequest 创建接口结构体
type CreateRaspConfigRequest struct {
	Name          string      `json:"name" form:"name" validate:"required,min=2,max=32"`
	Desc          string      `json:"desc" form:"desc" validate:"min=2,max=100"`
	Status        bool        `json:"status" form:"status" validate:"required"`
	AgentMode     uint        `json:"agentMode" form:"agentMode" validate:"oneof=1 2 3"`
	ModuleConfigs interface{} `json:"moduleConfigs" form:"moduleConfigs"`
	LogPath       string      `json:"logPath" form:"logPath" validate:"required,min=2,max=32"`
	AgentConfigs  interface{} `json:"agentConfigs" form:"agentConfigs"`
	BinFileUrl    string      `json:"binFileUrl" form:"binFileUrl"`
	BinFileHash   string      `json:"binFileHash" form:"binFileHash"`
}

// 获取接口列表结构体
type RaspConfigListRequest struct {
	Name     string `json:"name" form:"name"`
	Status   uint   `json:"status" form:"status"`
	Creator  string `json:"creator" form:"creator"`
	PageNum  uint   `json:"pageNum" form:"pageNum"`
	PageSize uint   `json:"pageSize" form:"pageSize"`
}

type RaspConfigRequest struct {
	Key     string `json:"key" form:"key"`
}

// 批量删除接口结构体
type DeleteRaspConfigRequest struct {
	Ids []uint `json:"ids" form:"ids"`
}
