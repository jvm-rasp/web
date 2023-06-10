package model

import "gorm.io/gorm"

type HostResource struct {
	gorm.Model
	HostName     string `gorm:"type:varchar(255);index;comment:'实例名称'" json:"hostName"`
	Ip           string `gorm:"type:varchar(255);index;comment:'ip地址'" json:"ip"`
	ResourceName string `gorm:"type:varchar(255);index;comment:'资源名称'" json:"resourceName"`
}
