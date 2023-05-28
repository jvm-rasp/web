package repository

import (
	"gorm.io/gorm"
	"server/common"
	"server/model"
)

type IRaspComponentRepository interface {
	CreateRaspComponent(component *model.RaspComponent) error
	CreateRaspUpgradeComponent(component *model.RaspUpgradeComponent) error
	GetRaspComponentByName(componentName string) (*model.RaspComponent, error)
	GetRaspComponentsByGuid(parentGuid string) ([]*model.RaspComponent, int64, error)
	GetRaspUpgradeComponentsByGuid(parentGuid string) ([]*model.RaspUpgradeComponent, int64, error)
	GetRaspComponentsByGuidAndName(parentGuid string, name string) (*model.RaspComponent, error)
	DeleteRaspComponentByIds(ids []uint) error
	DeleteRaspUpgradeComponentById(id uint) error
	DeleteRaspComponentByGuid(parentGuid string) error
	DeleteRaspComponentByGuids(parentGuid []string) error
	UpdateRaspComponent(component *model.RaspComponent) error
}

type RaspComponentRepository struct {
}

func NewRaspComponentRepository() IRaspComponentRepository {
	return RaspComponentRepository{}
}

func (a RaspComponentRepository) CreateRaspComponent(component *model.RaspComponent) error {
	err := common.DB.Create(component).Error
	return err
}

func (a RaspComponentRepository) CreateRaspUpgradeComponent(component *model.RaspUpgradeComponent) error {
	err := common.DB.Create(component).Error
	return err
}

func (a RaspComponentRepository) GetRaspComponentByName(componentName string) (*model.RaspComponent, error) {
	var record *model.RaspComponent
	err := common.DB.Model(&model.RaspComponent{}).Where("component_name = ?", componentName).First(&record).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return record, err
}

func (a RaspComponentRepository) GetRaspComponentsByGuid(parentGuid string) ([]*model.RaspComponent, int64, error) {
	var list []*model.RaspComponent
	var total int64
	err := common.DB.Model(&model.RaspComponent{}).Where("parent_guid = ?", parentGuid).Count(&total).Find(&list).Error
	return list, total, err
}

func (a RaspComponentRepository) GetRaspUpgradeComponentsByGuid(parentGuid string) ([]*model.RaspUpgradeComponent, int64, error) {
	var list []*model.RaspUpgradeComponent
	var total int64
	err := common.DB.Model(&model.RaspUpgradeComponent{}).Where("parent_guid = ?", parentGuid).Count(&total).Find(&list).Error
	return list, total, err
}

func (a RaspComponentRepository) GetRaspComponentsByGuidAndName(parentGuid string, name string) (*model.RaspComponent, error) {
	var record *model.RaspComponent
	err := common.DB.Model(&model.RaspComponent{}).Where("parent_guid = ?", parentGuid).Where("component_name = ?", name).First(&record).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return record, err
}

func (a RaspComponentRepository) DeleteRaspComponentByIds(ids []uint) error {
	err := common.DB.Where("id IN (?)", ids).Unscoped().Delete(&model.RaspComponent{}).Error
	return err
}

func (a RaspComponentRepository) DeleteRaspComponentByGuid(parentGuid string) error {
	err := common.DB.Where("parent_guid IN (?)", parentGuid).Unscoped().Delete(&model.RaspComponent{}).Error
	return err
}

func (a RaspComponentRepository) DeleteRaspComponentByGuids(parentGuid []string) error {
	err := common.DB.Where("parent_guid IN (?)", parentGuid).Unscoped().Delete(&model.RaspComponent{}).Error
	return err
}

func (a RaspComponentRepository) DeleteRaspUpgradeComponentById(id uint) error {
	err := common.DB.Where("id IN (?)", id).Unscoped().Delete(&model.RaspUpgradeComponent{}).Error
	return err
}

func (a RaspComponentRepository) UpdateRaspComponent(component *model.RaspComponent) error {
	err := common.DB.Save(component).Error
	return err
}
