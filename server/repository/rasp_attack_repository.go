package repository

import (
	"server/common"
	"server/model"
	"server/vo"
	"strings"
)

type IRaspAttackRepository interface {
	GetRaspAttacks(req *vo.RaspAttackListRequest) ([]*model.RaspAttack, int64, error)
	CreateRaspAttack(config *model.RaspAttack) error
	DeleteRaspAttack(ids []uint) error
}

type RaspAttackRepository struct {
}

func NewRaspAttackRepository() IRaspAttackRepository {
	return RaspAttackRepository{}
}

func (a RaspAttackRepository) CreateRaspAttack(attack *model.RaspAttack) error {
	err := common.DB.Create(attack).Error
	return err
}

func (a RaspAttackRepository) GetRaspAttacks(req *vo.RaspAttackListRequest) ([]*model.RaspAttack, int64, error) {
	var list []*model.RaspAttack
	db := common.DB.Model(&model.RaspAttack{}).Order("created_at DESC")

	name := strings.TrimSpace(req.HostName)
	if name != "" {
		db = db.Where("host_name = ?", name)
	}

	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	pageNum := int(req.PageNum)
	pageSize := int(req.PageSize)
	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}
	return list, total, err
}

func (r RaspAttackRepository) DeleteRaspAttack(ids []uint) error {
	err := common.DB.Where("id IN (?)", ids).Unscoped().Delete(&model.RaspAttack{}).Error
	return err
}
