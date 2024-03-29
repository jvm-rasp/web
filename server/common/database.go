package common

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"server/config"
	"server/model"
)

// 全局mysql数据库变量
var DB *gorm.DB

func InitDB() {
	driver := config.Conf.Database.Driver
	if driver == "sqlite" {
		initSQLite()
	} else if driver == "mysql" {
		initMySQL()
	}
}

// 初始化mysql数据库
func initSQLite() {
	db, err := gorm.Open(sqlite.Open(config.Conf.Database.Source), &gorm.Config{
		Logger: logger.Default.LogMode(config.Conf.Database.LogMode),
	})
	if err != nil {
		Log.Panicf("初始化sqlite数据库异常: %v", err)
		panic(fmt.Errorf("初始化sqlite数据库异常: %v", err))
	}
	// 全局DB赋值
	DB = db
	// 自动迁移表结构
	dbAutoMigrate()
	Log.Infof("初始化sqlite数据库完成!")
}

func initMySQL() {
	dsn := config.Conf.Database.Source
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(config.Conf.Database.LogMode),
	})
	if err != nil {
		Log.Panicf("初始化mysql数据库异常: %v", err)
		panic(fmt.Errorf("初始化mysql数据库异常: %v", err))
	}
	// 全局DB赋值
	DB = db
	// 自动迁移表结构
	dbAutoMigrate()
	Log.Infof("初始化mysql数据库完成!")
}

func FreeDBSpace() {
	// 释放空间
	switch config.Conf.Database.Driver {
	case "sqlite":
		err := DB.Exec("vacuum").Error
		if err != nil {
			Log.Errorf("释放磁盘空间错误, error: %v", err)
			return
		}
	case "mysql":
	}
}

// 自动迁移表结构
func dbAutoMigrate() {
	err := DB.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Menu{},
		&model.Api{},
		&model.OperationLog{},
		&model.RaspConfig{},
		&model.RaspConfigHistory{},
		&model.RaspModule{},
		&model.RaspComponent{},
		&model.RaspUpgradeComponent{},
		&model.RaspHost{},
		&model.JavaProcessInfo{},
		&model.RaspAttack{},
		&model.RaspAttackDetail{},
		&model.RaspFile{},
		&model.RaspErrorLogs{},
		&model.HostResource{},
	)
	if err != nil {
		Log.Errorf("初始化数据表出错, errpr: %v", err)
	}

	// 先清空apis表和casbin_rule表并进行重建
	err = DB.Table("apis").Where("1=1").Unscoped().Delete(&model.Api{}).Error
	if err != nil {
		Log.Errorf("清空apis表出错, errpr: %v", err)
	}
	err = DB.Table("casbin_rule").Where("1=1").Unscoped().Delete(&model.Api{}).Error
	if err != nil {
		if err.Error() != "no such table: casbin_rule" {
			Log.Errorf("清空casbin_rule表出错, errpr: %v", err)
		}
	}
	FreeDBSpace()
}
