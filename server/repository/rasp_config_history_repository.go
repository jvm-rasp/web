package repository

import (
	"gorm.io/gorm"
	"server/common"
	"server/model"
)

type IRaspConfigHistoryRepository interface {
	CreateRaspConfigHistory(config *model.RaspConfigHistory) error // 创建接口
	GetRaspConfigHistoryByGuid(configGuid string) ([]*model.RaspConfigHistory, int64, error)
	GetRaspConfigHistoryDataByGuid(configGuid string, version int) (*model.RaspConfigHistory, error)
	DeleteRaspConfigHistory(configGuid string) error
}

type RaspConfigHistoryRepository struct {
}

func NewRaspConfigHistoryRepository() IRaspConfigHistoryRepository {
	return RaspConfigHistoryRepository{}
}

func (this RaspConfigHistoryRepository) CreateRaspConfigHistory(config *model.RaspConfigHistory) error {
	err := common.DB.Create(config).Error
	return err
}

func (this RaspConfigHistoryRepository) GetRaspConfigHistoryByGuid(configGuid string) ([]*model.RaspConfigHistory, int64, error) {
	var list []*model.RaspConfigHistory
	var total int64
	db := common.DB.Model(&model.RaspConfigHistory{}).Order("version DESC")
	err := db.Where("parent_guid = ?", configGuid).Limit(20).Count(&total).Find(&list).Error
	return list, total, err
}

func (this RaspConfigHistoryRepository) GetRaspConfigHistoryDataByGuid(configGuid string, version int) (*model.RaspConfigHistory, error) {
	var record *model.RaspConfigHistory
	err := common.DB.Model(&model.RaspConfigHistory{}).Where("parent_guid", configGuid).Where("version", version).First(&record).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return record, err
}

func (this RaspConfigHistoryRepository) DeleteRaspConfigHistory(configGuid string) error {
	err := common.DB.Where("parent_guid IN (?)", configGuid).Unscoped().Delete(&model.RaspConfigHistory{}).Error
	return err
}
