package model

import "gorm.io/gorm"

type RaspHost struct {
	gorm.Model
	HostName              string  `gorm:"type:varchar(255);index;comment:'实例名称'" json:"hostName"`
	Ip                    string  `gorm:"type:varchar(255);index;comment:'ip地址'" json:"ip"`
	AgentMode             string  `gorm:"type:varchar(255);comment:'agent接入模式:disable;static;dynamic'" json:"agentMode"`
	AgentInfo             string  `gorm:"type:varchar(65535);comment:'java agent信息'" json:"agentInfo"`
	ConfigId              uint64  `gorm:"type:int(11);index;comment:'策略id'" json:"configId"`
	HeatbeatTime          string  `gorm:"type:varchar(128);index;comment:'最近一次的心跳时间'" json:"heatbeatTime"`
	InstallDir            string  `gorm:"type:varchar(255);comment:'可执行文件安装的绝对路径'" json:"installDir"`
	Version               string  `gorm:"type:varchar(64);comment:'rasp 版本'" json:"version"`
	ExeFileHash           string  `gorm:"type:varchar(255);comment:'可执行文件的hash'" json:"exeFileHash"`
	AgentConfigUpdateTime string  `gorm:"type:varchar(128);comment:'agent 配置更新时间" json:"agentConfigUpdateTime"`
	OsType                string  `gorm:"type:varchar(128);comment:'实例操作系统类型:darwin、windowns、linux等'" json:"osType"`
	TotalMem              float64 `gorm:"type:float64;comment:'实例总内存,单位GB'" json:"totalMem"`
	CpuCounts             float64 `gorm:"type:float64;comment:'实例逻辑cpu数量'" json:"cpuCounts"`
	FreeDisk              float64 `gorm:"type:float64;comment:'实例可用磁盘容量,单位GB'" json:"freeDisk"`
	BuildDateTime         string  `gorm:"type:varchar(128);comment:'代码信息：编译时间'" json:"buildDateTime"`
	BuildGitBranch        string  `gorm:"type:varchar(128);comment:'代码信息：分支'" json:"buildGitBranch"`
	BuildGitCommit        string  `gorm:"type:varchar(128);comment:'代码信息：commit''" json:"buildGitCommit"`
	Tag                   string  `gorm:"type:varchar(128);comment:'标签'" json:"tag"`
	Desc                  string  `gorm:"type:varchar(128);comment:'说明'" json:"desc"`
}
