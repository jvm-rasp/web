package model

import "gorm.io/gorm"

type RaspErrorLogs struct {
	gorm.Model
	Topic   string `gorm:"type:varchar(255);index;comment:'日志类型'" json:"topic"`
	Time    string `gorm:"type:varchar(128);index;comment:'日志时间'" json:"time"`
	Level   string `gorm:"type:varchar(255);comment:'日志等级'" json:"level"`
	Message string `gorm:"type:text;comment:'日志内容'" json:"message"`
}
