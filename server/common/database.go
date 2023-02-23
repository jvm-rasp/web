package common

import (
	"fmt"
	"server/config"
	"server/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 全局mysql数据库变量
var DB *gorm.DB

// 初始化mysql数据库
func InitMysql() {
	db, err := gorm.Open(sqlite.Open("jrasp.db"), &gorm.Config{})
	if err != nil {
		Log.Panicf("初始化mysql数据库异常: %v", err)
		panic(fmt.Errorf("初始化mysql数据库异常: %v", err))
	}

	// 开启mysql日志
	if config.Conf.Mysql.LogMode {
		db.Debug()
	}
	// 全局DB赋值
	DB = db
	// 自动迁移表结构
	dbAutoMigrate()
	Log.Infof("初始化数据库完成!")
}

// 自动迁移表结构
func dbAutoMigrate() {
	DB.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Menu{},
		&model.Api{},
		&model.OperationLog{},
	)
}