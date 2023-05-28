package model

import "gorm.io/gorm"

type JavaProcessInfo struct {
	gorm.Model
	HostName       string `gorm:"type:varchar(255);index;comment:'实例名称'" json:"hostName"`
	CmdlineInfo    string `gorm:"type:text;comment:'cmdline'" json:"cmdlineInfo"`
	AppNamesInfo   string `gorm:"type:varchar(255);comment:'被保护应用名'" json:"appNamesInfo"`
	StartTime      string `gorm:"type:varchar(255);comment:'进程启动时间'" json:"startTime"`
	Pid            int    `gorm:"type:int(11);comment:'进程pid'" json:"pid"`
	Status         int    `gorm:"type:int(11);index;comment:'注入状态 0=未安装; 1=防护中; 2=安装失败'" json:"status"`
	Message        string `gorm:"type:varchar(255);index;comment:'信息'" json:"message"`
	DependencyInfo string `gorm:"type:text;comment:'jar包依赖信息'" json:"dependencyInfo"`
}
