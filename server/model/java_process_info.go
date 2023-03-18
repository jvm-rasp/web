package model

import "gorm.io/gorm"

type JavaProcessInfo struct {
	gorm.Model
	HostName       string `gorm:"type:varchar(255);index;comment:'实例名称'" json:"hostName"`
	CmdlineInfo    string `gorm:"type:text;index;comment:'cmdline'" json:"cmdlineInfo"`
	StartTime      string `gorm:"type:varchar(255);comment:'进程启动时间'" json:"startTime"`
	Pid            int    `gorm:"type:int(11);comment:'进程pid'" json:"pid"`
	Status         string `gorm:"type:int(11);index;comment:'注入状态'" json:"status"`
	DependencyInfo string `gorm:"type:text;index;comment:'jar包依赖信息'" json:"dependencyInfo"`
}
