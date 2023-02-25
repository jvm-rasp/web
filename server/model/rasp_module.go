package model

import "gorm.io/gorm"

type RaspModule struct {
	gorm.Model
	Name              string `gorm:"type:varchar(32);index;comment:'模块名称'" json:"name"`
	Version           string `gorm:"type:varchar(128);comment:'模块版本: 1.0.4、1.0.5'" json:"version"`
	Url               string `gorm:"type:varchar(1024);index;comment:'模块下载地址'" json:"url"`
	Hash              string `gorm:"type:varchar(128);comment:'模块hash'" json:"hash"`
	Type              string `gorm:"type:varchar(32);comment:'模块类型: hook、algorithm和other'" json:"type"`
	MiddlewareName    string `gorm:"type:varchar(128);comment:'中间件名称：tomcat、jetty'" json:"middlewareName"`
	MiddlewareVersion string `gorm:"type:varchar(128);comment:'中间件版本:8.0、9.0'" json:"middlewareVersion"`
	Tag               string `gorm:"type:varchar(128);comment:'模块标签'" json:"tag"`
	Parameters        string `gorm:"type:varchar(4096);comment:'模块参数模版: json字符串'" json:"parameters"`
	Status            uint   `gorm:"type:tinyint(1);default:1;comment:'是否禁用(0,禁用/1,启用, 默认启用)'" json:"status"`
	Desc              string `gorm:"type:varchar(128);comment:'模块说明'" json:"desc"`
	Creator           string `gorm:"type:varchar(20);comment:'模块创建人'" json:"creator"`
}
