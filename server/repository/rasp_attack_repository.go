package repository

import (
	"server/common"
	"server/model"
	"server/vo"
	"strconv"
	"strings"
)

type IRaspAttackRepository interface {
	GetRaspAttacks(req *vo.RaspAttackListRequest) ([]*model.RaspAttack, int64, error)
	GetRaspAttackById(id uint) (*model.RaspAttack, error)
	GetRaspAttackDetail(parentGuid string) (*model.RaspAttackDetail, error)
	CreateRaspAttack(attack *model.RaspAttack) error
	CreateRaspAttackDetail(detail *model.RaspAttackDetail) error
	DeleteRaspAttack(guids []string) error
	DeleteRaspDetail(guids []string) error
	UpdateRaspAttack(attack *model.RaspAttack) error
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

func (a RaspAttackRepository) CreateRaspAttackDetail(detail *model.RaspAttackDetail) error {
	err := common.DB.Create(detail).Error
	return err
}

func (a RaspAttackRepository) GetRaspAttackDetail(parentGuid string) (*model.RaspAttackDetail, error) {
	var detail *model.RaspAttackDetail
	db := common.DB.Model(&model.RaspAttackDetail{}).Where("parent_guid = ?", parentGuid)
	err := db.Find(&detail).Error
	return detail, err
}

func (a RaspAttackRepository) GetRaspAttacks(req *vo.RaspAttackListRequest) ([]*model.RaspAttack, int64, error) {
	var list []*model.RaspAttack
	db := common.DB.Model(&model.RaspAttack{}).Order("created_at DESC")

	name := strings.TrimSpace(req.HostName)
	if name != "" {
		db = db.Where("host_name = ?", name)
	}
	isBlocked, err := strconv.ParseBool(req.IsBlocked)
	if err == nil {
		db = db.Where("is_blocked = ?", isBlocked)
	}
	result, err := strconv.ParseInt(req.HandleResult, 10, 32)
	if err == nil {
		db = db.Where("handle_result = ?", result)
	}

	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	var total int64
	err = db.Count(&total).Error
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

func (r RaspAttackRepository) DeleteRaspAttack(guids []string) error {
	err := common.DB.Where("row_guid IN (?)", guids).Unscoped().Delete(&model.RaspAttack{}).Error
	return err
}

func (r RaspAttackRepository) DeleteRaspDetail(guids []string) error {
	err := common.DB.Where("parent_guid IN (?)", guids).Unscoped().Delete(&model.RaspAttackDetail{}).Error
	return err
}

func (r RaspAttackRepository) GetRaspAttackById(id uint) (*model.RaspAttack, error) {
	var record *model.RaspAttack
	err := common.DB.Find(&record, "id = ?", id).Error
	return record, err
}

func (r RaspAttackRepository) UpdateRaspAttack(attack *model.RaspAttack) error {
	err := common.DB.Save(attack).Error
	return err
}
