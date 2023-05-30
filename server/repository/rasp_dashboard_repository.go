package repository

import (
	"server/common"
	"server/model"
)

type IRaspDashboardRepository interface {
	GetHighLevelCount(level int) (int64, error)
	GetBlockCount() (int64, error)
	GetAllCount() (int64, error)
	GetAttackTypes() ([]model.AttackTypes, error)
}

type RaspDashboardRepository struct {
}

func NewRaspDashboardRepository() IRaspDashboardRepository {
	return RaspDashboardRepository{}
}

func (r RaspDashboardRepository) GetHighLevelCount(level int) (int64, error) {
	var count int64
	err := common.DB.Model(&model.RaspAttack{}).Where("level >= ?", level).Count(&count).Error
	return count, err
}

func (r RaspDashboardRepository) GetBlockCount() (int64, error) {
	var count int64
	err := common.DB.Model(&model.RaspAttack{}).Where("is_blocked = ?", true).Count(&count).Error
	return count, err
}

func (r RaspDashboardRepository) GetAllCount() (int64, error) {
	var count int64
	err := common.DB.Model(&model.RaspAttack{}).Count(&count).Error
	return count, err
}

func (r RaspDashboardRepository) GetAttackTypes() ([]model.AttackTypes, error) {
	var results []model.AttackTypes
	err := common.DB.Model(&model.RaspAttack{}).Select("attack_type as name, count(*) as value").Group("attack_type").Find(&results).Error
	return results, err
}
