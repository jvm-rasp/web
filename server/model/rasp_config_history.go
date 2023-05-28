package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type RaspConfigHistory struct {
	gorm.Model
	ParentGuid    string         `gorm:"type:varchar(255);comment:'关联的配置guid'" json:"parentGuid"`
	Version       int            `gorm:"type:int;comment:'关联的配置ID'" json:"version"`
	AgentMode     uint           `gorm:"type:tinyint(1);comment:'Agent模式'" json:"agentMode"`
	ModuleConfigs datatypes.JSON `gorm:"type:text;comment:'检测插件配置'" json:"moduleConfigs"`
	LogPath       string         `gorm:"type:varchar(255);comment:'日志路径'" json:"logPath"`
	AgentConfigs  datatypes.JSON `gorm:"type:text;comment:'Agent配置信息'" json:"agentConfigs"`
	RaspBinInfo   datatypes.JSON `gorm:"type:text;comment:'Daemon下载信息'" json:"raspBinInfo"`
	RaspLibInfo   datatypes.JSON `gorm:"type:text;comment:'Agent Lib下载信息'" json:"raspLibInfo"`
	Desc          string         `gorm:"type:varchar(255);comment:'策略说明'" json:"desc"`
	Source        string         `gorm:"type:varchar(255);comment:'策略来源'" json:"source"`
}

type RaspModuleConfig struct {
	ModuleName string               `json:"moduleName"`
	Components []ComponentConfig    `json:"components"`
	Parameters RaspModuleUserConfig `json:"parameters"`
}

type RaspModuleUserConfig struct {
	Action map[string]int    `json:"action"`
	CnMap  map[string]string `json:"cn_map"`
}
