package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type RaspModule struct {
	gorm.Model
	ModuleName        string         `gorm:"type:varchar(32);comment:'模块名称'" json:"moduleName"`
	ModuleType        uint           `gorm:"type:tinyint(1);comment:'模块类型: hook、algorithm和other'" json:"moduleType"`
	ModuleVersion     string         `gorm:"type:varchar(10);comment:'模块版本'" json:"moduleVersion"`
	DownLoadURL       string         `gorm:"type:varchar(1024);comment:'模块下载地址'" json:"downLoadURL"`
	Md5               string         `gorm:"type:varchar(128);comment:'模块hash'" json:"md5"`
	Desc              string         `gorm:"type:varchar(128);comment:'模块说明'" json:"desc"`
	Status            bool           `gorm:"type:tinyint(1);default:1;comment:'是否禁用(0,禁用/1,启用, 默认启用)'" json:"status"`
	Parameters        datatypes.JSON `gorm:"type:varchar(4096);comment:'模块参数模版: json字符串'" json:"parameters"`
	Creator           string         `gorm:"type:varchar(20);comment:'模块创建人'" json:"creator"`
	Operator          string         `gorm:"type:varchar(20);comment:'操作人'" json:"operator"`
}
