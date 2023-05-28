package vo

import "gorm.io/datatypes"

// CreateRaspConfigRequest 创建接口结构体
type CreateRaspConfigRequest struct {
	Name          string         `json:"name" form:"name" validate:"required"`
	Version       int            `json:"version" form:"version" validate:"gte=0"`
	Desc          string         `json:"desc" form:"desc" validate:"min=2,max=100"`
	Status        bool           `json:"status" form:"status" validate:"boolean"`
	AgentMode     uint           `json:"agentMode" form:"agentMode" validate:"oneof=0 1 2"`
	ModuleConfigs datatypes.JSON `json:"moduleConfigs" form:"moduleConfigs"`
	LogPath       string         `json:"logPath" form:"logPath" validate:"required"`
	AgentConfigs  datatypes.JSON `json:"agentConfigs" form:"agentConfigs"`
	RaspBinInfo   datatypes.JSON `json:"raspBinInfo" form:"raspBinInfo"`
	RaspLibInfo   datatypes.JSON `json:"raspLibInfo" form:"raspLibInfo"`
	HistoryDesc   string         `json:"historyDesc" form:"historyDesc" validate:"required"`
}

type UpdateRaspConfigRequest struct {
	ID            uint           `json:"id" form:"id" validate:"required,min=1"`
	Name          string         `json:"name" form:"name" validate:"required"`
	Version       int            `json:"version" form:"version" validate:"gte=0"`
	Desc          string         `json:"desc" form:"desc" validate:"min=2,max=100"`
	Status        bool           `json:"status" form:"status" validate:"boolean"`
	AgentMode     uint           `json:"agentMode" form:"agentMode" validate:"oneof=0 1 2"`
	ModuleConfigs datatypes.JSON `json:"moduleConfigs" form:"moduleConfigs"`
	LogPath       string         `json:"logPath" form:"logPath" validate:"required"`
	AgentConfigs  datatypes.JSON `json:"agentConfigs" form:"agentConfigs"`
	RaspBinInfo   datatypes.JSON `json:"raspBinInfo" form:"raspBinInfo"`
	RaspLibInfo   datatypes.JSON `json:"raspLibInfo" form:"raspLibInfo"`
	HistoryDesc   string         `json:"historyDesc" form:"historyDesc"`
	IsNewVersion  bool           `json:"IsNewVersion" form:"IsNewVersion" validate:"boolean"`
}

// 获取接口列表结构体
type RaspConfigListRequest struct {
	Name     string `json:"name" form:"name"`
	Status   string `json:"status" form:"status"`
	Creator  string `json:"creator" form:"creator"`
	PageNum  uint   `json:"pageNum" form:"pageNum"`
	PageSize uint   `json:"pageSize" form:"pageSize"`
}

type RaspConfigRequest struct {
	Key string `json:"key" form:"key"`
}

// 批量删除接口结构体
type DeleteRaspConfigRequest struct {
	Ids []uint `json:"ids" form:"ids"`
}

type UpdateRaspConfigStatusRequest struct {
	ID uint `json:"id" form:"id" validate:"required"`
}

type UpdateRaspConfigDefaultRequest struct {
	ID        uint `json:"id" form:"id" validate:"required"`
	IsDefault bool `json:"isDefault" form:"isDefault" validate:"boolean"`
}

type PushRaspConfigRequest struct {
	ID uint `json:"id" form:"id" validate:"required"`
}

type CopyRaspConfigRequest struct {
	ID uint `json:"id" form:"id" validate:"required"`
}

type RaspCheckboxModuleListRequest struct {
}

type ExportRaspConfigRequest struct {
	ID uint `json:"id" form:"id" validate:"required"`
}

type SyncRaspConfigRequest struct {
	SrcConfigId      uint `json:"srcConfigId" form:"srcConfigId" validate:"required"`
	SrcConfigVersion int  `json:"srcConfigVersion" form:"srcConfigVersion" validate:"required"`
	DstConfigId      uint `json:"dstConfigId" form:"dstConfigId" validate:"required"`
	SyncOptions      int  `json:"syncOptions" form:"syncOptions" validate:"required"`
}
