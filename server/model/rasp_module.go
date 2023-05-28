package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type RaspModule struct {
	gorm.Model
	RowGuid       string           `gorm:"type:varchar(255);index;comment:'标识guid'" json:"rowGuid"`
	ModuleName    string           `gorm:"type:varchar(32);comment:'模块名称'" json:"moduleName"`
	ModuleVersion string           `gorm:"type:varchar(10);comment:'模块版本'" json:"moduleVersion"`
	Desc          string           `gorm:"type:varchar(128);comment:'模块说明'" json:"desc"`
	Parameters    datatypes.JSON   `gorm:"type:varchar(4096);comment:'模块参数模版: json字符串'" json:"parameters"`
	Upgradable    bool             `gorm:"type:tinyint(0);default:0;comment:'模块是否可更新: json字符串'" json:"upgradable"`
	Creator       string           `gorm:"type:varchar(20);comment:'模块创建人'" json:"creator"`
	Operator      string           `gorm:"type:varchar(20);comment:'操作人'" json:"operator"`
	Components    []*RaspComponent `gorm:"-" json:"components"`
}

type RaspComponent struct {
	gorm.Model
	ParentGuid       string         `gorm:"type:varchar(255);index;comment:'关联的防护模块标识guid'" json:"parentGuid"`
	ComponentName    string         `gorm:"type:varchar(32);comment:'组件名称'" json:"componentName"`
	ComponentType    uint           `gorm:"type:tinyint(1);comment:'组件类型: hook、algorithm和other'" json:"componentType"`
	ComponentVersion string         `gorm:"type:varchar(10);comment:'组件版本'" json:"componentVersion"`
	DownLoadURL      string         `gorm:"type:varchar(1024);comment:'组件下载地址'" json:"downLoadURL"`
	Md5              string         `gorm:"type:varchar(128);comment:'组件hash'" json:"md5"`
	Parameters       datatypes.JSON `gorm:"type:varchar(4096);comment:'参数模版: json字符串'" json:"parameters"`
}

type RaspUpgradeComponent struct {
	gorm.Model
	ParentGuid string `gorm:"type:varchar(255);index;comment:'关联的防护模块标识guid'" json:"parentGuid"`
	Md5        string `gorm:"type:varchar(128);comment:'组件hash'" json:"md5"`
}
