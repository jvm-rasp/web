package repository

import (
	"errors"
	"gorm.io/gorm"
	"server/common"
	"server/model"
)

type ISystemSettingRepository interface {
	CreateSetting(setting *model.SystemSetting) error
	DeleteSettingByName(name string) error
	SaveSetting(setting *model.SystemSetting) error
	GetSettingByName(name string) (*model.SystemSetting, error)
	GetSettings() ([]*model.SystemSetting, error)
	SetSetting(name string, value string) error
}

type SystemSettingRepository struct {
}

func NewSystemSettingRepository() ISystemSettingRepository {
	return SystemSettingRepository{}
}

func (this SystemSettingRepository) CreateSetting(setting *model.SystemSetting) error {
	err := common.DB.Create(setting).Error
	return err
}

func (this SystemSettingRepository) DeleteSettingByName(name string) error {
	err := common.DB.Where("name = ?", name).Unscoped().Delete(&model.SystemSetting{}).Error
	return err
}

func (this SystemSettingRepository) SaveSetting(setting *model.SystemSetting) error {
	err := common.DB.Save(setting).Error
	return err
}

func (this SystemSettingRepository) GetSettingByName(name string) (*model.SystemSetting, error) {
	var record *model.SystemSetting
	err := common.DB.Where("name = ?", name).First(&record).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return record, err
}

func (this SystemSettingRepository) GetSettings() ([]*model.SystemSetting, error) {
	var list []*model.SystemSetting
	db := common.DB.Model(&model.SystemSetting{}).Order("created_at DESC")
	err := db.Find(&list).Error
	return list, err
}

func (this SystemSettingRepository) SetSetting(name string, value string) error {
	setting, err := this.GetSettingByName(name)
	if err != nil {
		return err
	}
	if setting == nil {
		return errors.New("未找到配置: " + name)
	}
	setting.Value = value
	err = this.SaveSetting(setting)
	return err
}
