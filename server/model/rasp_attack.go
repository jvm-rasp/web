package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type RaspAttack struct {
	gorm.Model
	RowGuid      string    `gorm:"type:varchar(255);index;comment:'标识guid'" json:"rowGuid"`
	HostName     string    `gorm:"type:varchar(255);index;comment:'实例名称'" json:"hostName"`
	HostIp       string    `gorm:"type:varchar(255);index;comment:'实例ip'" json:"hostIp"`
	RemoteIp     string    `gorm:"type:varchar(255);index;comment:'攻击者ip地址'" json:"remoteIp"`
	AttackType   string    `gorm:"type:varchar(255);comment:'攻击类型'" json:"attackType"`
	IsBlocked    bool      `gorm:"type:boolean;comment:'是否阻断：0,放行；1：阻断'" json:"isBlocked"`
	Level        int       `gorm:"type:int(11);comment:'危险等级：0～100整数'" json:"level"`
	HandleResult int       `gorm:"type:int(11);comment:'处理状态，未处理=0(初始化状态)、确认漏洞=1、误报=2、忽略=3'" json:"handleResult"`
	RequestUri   string    `gorm:"type:varchar(255);comment:'请求路径'" json:"requestUri"`
	AttackTime   time.Time `gorm:"type:datetime;comment:'攻击时间'" json:"attackTime"`
	Tag          string    `gorm:"type:varchar(128);comment:'标签'" json:"tag"`
	Desc         string    `gorm:"type:varchar(128);comment:'说明'" json:"desc"`
}

type RaspAttackDetail struct {
	gorm.Model
	ParentGuid string         `gorm:"type:varchar(255);index;comment:'关联的攻击记录标识guid'" json:"parentGuid"`
	Context    datatypes.JSON `gorm:"type:text;comment:'上下文信息'" json:"context"`
	AppName    string         `gorm:"type:varchar(255);comment:'应用名称'" json:"appName"`
	StackTrace string         `gorm:"type:text;comment:'堆栈信息'" json:"stackTrace"`
	Payload    string         `gorm:"type:text;comment:'攻击payload'" json:"payload"`
	IsBlocked  bool           `gorm:"type:tinyint;comment:'是否阻断：0,放行；1：阻断'" json:"isBlocked"`
	AttackType string         `gorm:"type:varchar(255);comment:'攻击类型'" json:"attackType"`
	Algorithm  string         `gorm:"type:varchar(255);comment:'检测算法'" json:"algorithm"`
	Extend     string         `gorm:"type:text;comment:'告警信息'" json:"extend"`
	AttackTime time.Time      `gorm:"type:datetime;comment:'攻击时间'" json:"attackTime"`
	Level      int            `gorm:"type:int(11);comment:'危险等级：0～100整数'" json:"level"`
	MetaInfo   string         `gorm:"type:varchar(255);comment:'元数据信息'" json:"metaInfo"`
}
