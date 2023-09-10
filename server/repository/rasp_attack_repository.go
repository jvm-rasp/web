package repository

import (
	"fmt"
	"gorm.io/gorm"
	"server/common"
	"server/config"
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
	DeleteAttackLogsByJob(maxSize int) error
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
	err := db.First(&detail).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return detail, nil
}

func (a RaspAttackRepository) GetRaspAttacks(req *vo.RaspAttackListRequest) ([]*model.RaspAttack, int64, error) {
	var list []*model.RaspAttack
	db := common.DB.Model(&model.RaspAttack{}).Order("created_at DESC")

	name := strings.TrimSpace(req.HostName)
	if name != "" {
		db = db.Where("host_name = ?", name)
	}
	url := strings.TrimSpace(req.Url)
	if url != "" {
		db = db.Where("request_uri like ?", fmt.Sprintf("%%%s%%", url))
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
	err := common.DB.First(&record, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return record, err
}

func (r RaspAttackRepository) UpdateRaspAttack(attack *model.RaspAttack) error {
	err := common.DB.Save(attack).Error
	return err
}

func (r RaspAttackRepository) DeleteAttackLogsByJob(maxSize int) error {
	var logs []model.RaspAttack
	if err := common.DB.Limit(1).Offset(maxSize).Order("id desc").Find(&logs).Error; err != nil {
		return err
	}
	if logs != nil && len(logs) > 0 {
		maxId := logs[0].ID
		if err := common.DB.Where("id <= ?", maxId).Delete(&model.RaspAttack{}).Error; err != nil {
			return err
		}
		// 删除rasp_attack_details表
		guid := logs[0].RowGuid
		detail, err := r.GetRaspAttackDetail(guid)
		if err == nil && detail != nil {
			if err = common.DB.Where("id <= ?", detail.ID).Delete(&model.RaspAttackDetail{}).Error; err != nil {
				return err
			}
		}
		// 释放空间
		switch config.Conf.Database.Driver {
		case "sqlite":
			if err := common.DB.Exec("vacuum").Error; err != nil {
				return err
			}
		case "mysql":
		}
	}
	return nil
}
