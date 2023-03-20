package model

import "gorm.io/gorm"

type RaspAttack struct {
	gorm.Model
	HostName          string `gorm:"type:varchar(255);index;comment:'实例名称'" json:"hostName"`
	LocalIp           string `gorm:"type:varchar(255);index;comment:'被攻击者ip地址'" json:"localIp"`
	RemoteIp          string `gorm:"type:varchar(255);index;comment:'攻击者ip地址'" json:"remoteIp"`
	AttackType        string `gorm:"type:varchar(255);comment:'攻击类型'" json:"attackType"`
	CheckType         string `gorm:"type:varchar(255);comment:'检测算法类型以及信息'" json:"checkType"`
	IsBlocked         bool   `gorm:"type:boolean;comment:'是否阻断：0,放行；1：阻断'" json:"isBlocked"`
	Level             int    `gorm:"type:int(11);comment:'危险等级：0～100整数'" json:"level"`
	HandleStatus      string `gorm:"type:int(11);comment:'处理状态，未处理1(初始化状态)、处理中2(点击查看详情)、已经处理3'" json:"handleStatus"`
	HandleResult      string `gorm:"type:varchar(255);comment:'处理结果：确认漏洞、误拦截、忽略'" json:"handleResult"`
	StackTrace        string `gorm:"type:text;comment:'调用栈：json字符串格式'" json:"stackTrace"`
	HttpMethod        string `gorm:"type:text;comment:'http请求类型：get、post、put'" json:"httpMethod"`
	RequestProtocol   string `gorm:"type:varchar(16);comment:'http请求协议：rpc、http,dubbo'" json:"requestProtocol"`
	RequestUri        string `gorm:"type:varchar(255);comment:'请求路径'" json:"requestUri"`
	RequestParameters string `gorm:"type:varchar(2048);comment:'请求参数：json字符串格式'" json:"requestParameters"`
	AttackParameters  string `gorm:"type:varchar(2048);comment:'攻击参数：json字符串格式'" json:"attackParameters"`
	Cookies           string `gorm:"type:varchar(2048);comment:'请求cookie：json字符串格式'" json:"cookies"`
	Header            string `gorm:"type:varchar(2048);comment:'请求header：json字符串格式'" json:"header"`
	Body              string `gorm:"type:text;comment:'http body 信息'" json:"body"`
	AttackTime        string `gorm:"type:varchar(128);comment:'攻击时间'" json:"attackTime"`
	Tag               string `gorm:"type:varchar(128);comment:'标签'" json:"tag"`
	Desc              string `gorm:"type:varchar(128);comment:'说明'" json:"desc"`
}
