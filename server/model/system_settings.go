package model

import "gorm.io/gorm"

type SystemSetting struct {
	gorm.Model
	Name  string `gorm:"type:varchar(255);index;comment:'参数名称'" json:"name"`
	Type  string `gorm:"type:varchar(255);index;comment:'参数类型'" json:"type"`
	Value string `gorm:"type:varchar(255);index;comment:'参数值'" json:"value"`
}
