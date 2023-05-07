package model

import (
	"gorm.io/gorm"
)

type OperationLog struct {
	gorm.Model
	Username   string    `gorm:"type:varchar(20);comment:'用户登录名'" json:"username"`
	Ip         string    `gorm:"type:varchar(20);comment:'Ip地址'" json:"ip"`
	IpLocation string    `gorm:"type:varchar(20);comment:'Ip所在地'" json:"ipLocation"`
	Method     string    `gorm:"type:varchar(20);comment:'请求方式'" json:"method"`
	Path       string    `gorm:"type:varchar(100);comment:'访问路径'" json:"path"`
	Desc       string    `gorm:"type:varchar(100);comment:'说明'" json:"desc"`
	Status     int       `gorm:"type:int(4);comment:'响应状态码'" json:"status"`
	StartTime  string `gorm:"type:varchar(128);comment:'发起时间'" json:"startTime"`
	TimeCost   int64     `gorm:"type:int(6);comment:'请求耗时(ms)'" json:"timeCost"`
	UserAgent  string    `gorm:"type:varchar(20);comment:'浏览器标识'" json:"userAgent"`
}
