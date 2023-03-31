package model

import "gorm.io/gorm"

type RaspFile struct {
	gorm.Model
	FileName      string `gorm:"type:varchar(255);index;comment:'文件名称'" json:"fileName"`
	FileHash      string `gorm:"type:varchar(1024);comment:'文件Hash'" json:"fileHash"`
	DiskPath      string `gorm:"type:varchar(1024);comment:'磁盘路径'" json:"diskPath"`
	DownLoadUrl   string `gorm:"type:varchar(1024);comment:'下载路径'" json:"downLoadUrl"`
	ModuleName    string `gorm:"type:varchar(255);index;comment:'模块名称'" json:"moduleName"`
	ModuleVersion string `gorm:"type:varchar(255);comment:'模块版本'" json:"moduleVersion"`
	Creator       string `gorm:"type:varchar(20);comment:'创建人'" json:"creator"`
	Tag           string `gorm:"type:varchar(128);comment:'标签'" json:"tag"`
	Desc          string `gorm:"type:varchar(128);comment:'说明'" json:"desc"`
}
