package model

import "gorm.io/gorm"

type RaspConfig struct {
	gorm.Model
	Name    string `gorm:"type:varchar(32);index;comment:'策略名称'" json:"name"`
	Content string `gorm:"type:varchar(128);comment:'策略内容'" json:"content"`
	Tag     string `gorm:"type:varchar(128);comment:'策略标签'" json:"tag"`
	Status  uint   `gorm:"type:tinyint(1);default:1;comment:'是否禁用(0,禁用/1,启用, 默认启用)'" json:"status"`
	Desc    string `gorm:"type:varchar(128);comment:'策略说明'" json:"desc"`
	Creator string `gorm:"type:varchar(20);comment:'创建人'" json:"creator"`
}
