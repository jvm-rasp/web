package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type RaspConfig struct {
	gorm.Model
	Name          string         `gorm:"type:varchar(32);comment:'策略名称'" json:"name"`
	Desc          string         `gorm:"type:varchar(128);comment:'策略说明'" json:"desc"`
	Status        bool           `gorm:"type:tinyint(1);default:1;comment:'是否禁用(0,禁用/1,启用, 默认启用)'" json:"status"`
	Creator       string         `gorm:"type:varchar(20);comment:'创建人'" json:"creator"`
	Operator      string         `gorm:"type:varchar(20);comment:'操作人'" json:"operator"`
	AgentMode     uint           `gorm:"type:tinyint(1);default:1;comment:'Agent模式'" json:"agentMode"`
	ModuleConfigs datatypes.JSON `gorm:"type:text;comment:'检测插件配置'" json:"moduleConfigs"`
	LogPath       string         `gorm:"type:varchar(255);comment:'日志路径'" json:"logPath"`
	AgentConfigs  datatypes.JSON `gorm:"type:text;comment:'Agent配置信息'" json:"agentConfigs"`
	BinFileUrl    string         `gorm:"type:varchar(255);comment:'守护进程下载地址'" json:"binFileUrl"`
	BinFileHash   string         `gorm:"type:varchar(64);comment:'哈希值'" json:"binFileHash"`
}
