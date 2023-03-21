package model

type RaspConfig struct {
	ID            uint   `gorm:"primarykey;index" json:"id"`
	Name          string `gorm:"type:varchar(32);comment:'策略名称'" json:"name"`
	Desc          string `gorm:"type:varchar(128);comment:'策略说明'" json:"desc"`
	Status        bool   `gorm:"type:tinyint(1);default:1;comment:'是否禁用(0,禁用/1,启用, 默认启用)'" json:"status"`
	Creator       string `gorm:"type:varchar(20);comment:'创建人'" json:"creator"`
	Operator      string `gorm:"type:varchar(20);comment:'操作人'" json:"operator"`
	CreateTime    string `gorm:"type:varchar(50);comment:'创建时间'" json:"createTime"`
	UpdateTime    string `gorm:"type:varchar(50);comment:'更新时间'" json:"updateTime"`
	AgentMode     uint   `gorm:"type:tinyint(1);default:1;comment:'Agent模式'" json:"agentMode"`
	ModuleConfigs string `gorm:"type:text;comment:'检测插件配置'" json:"moduleConfigs"`
	LogPath       string `gorm:"type:varchar(255);comment:'日志路径'" json:"logPath"`
	AgentConfigs  string `gorm:"type:text;comment:'Agent配置信息'" json:"agentConfigs"`
	BinFileUrl    string `gorm:"type:varchar(255);comment:'守护进程下载地址'" json:"binFileUrl"`
	BinFileHash   string `gorm:"type:varchar(64);comment:'哈希值'" json:"binFileHash"`
}
