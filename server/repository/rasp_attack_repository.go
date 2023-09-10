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
	GetRaspAttacksAndDetailNoPage(req *vo.RaspAttackListRequest) ([]*model.RaspAttackWithDetail, int64, error)
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
func (a RaspAttackRepository) GetRaspAttacksAndDetailNoPage(req *vo.RaspAttackListRequest) ([]*model.RaspAttackWithDetail, int64, error) {
	var list []*model.RaspAttackWithDetail
	db := common.DB.Model(&model.RaspAttack{}).Order("rasp_attacks.created_at DESC")
	// 攻击时间范围
	attackStartTime := strings.TrimSpace(req.AttackStartTime)
	if attackStartTime != "" {
		db = db.Where("rasp_attacks.attack_time >= ?", attackStartTime)
	}
	attackEndTime := strings.TrimSpace(req.AttackEndTime)
	if attackStartTime != "" {
		db = db.Where("rasp_attacks.attack_time <= ?", attackEndTime)
	}
	// 主机IP
	hostIp := strings.TrimSpace(req.HostIp)
	if hostIp != "" {
		db = db.Where("rasp_attacks.host_ip = ?", hostIp)
	}
	// 攻击类型
	attackType := strings.TrimSpace(req.AttackType)
	if attackType != "" {
		db = db.Where("rasp_attacks.attack_type like ?", fmt.Sprintf("%%%s%%", attackType))
	}
	// 远程IP
	remoteIp := strings.TrimSpace(req.RemoteIp)
	if remoteIp != "" {
		db = db.Where("rasp_attacks.remote_ip = ?", remoteIp)
	}
	// 安全等级
	level := strings.TrimSpace(req.Level)
	if level != "" {
		if level == "2" {
			db = db.Where("rasp_attacks.level >= ?", 90)
		} else if level == "1" {
			db = db.Where("rasp_attacks.level < ?", 90)
		}
	}
	// url地址
	url := strings.TrimSpace(req.Url)
	if url != "" {
		db = db.Where("rasp_attacks.request_uri like ?", fmt.Sprintf("%%%s%%", url))
	}
	// 是否拦截
	isBlocked, err := strconv.ParseBool(req.IsBlocked)
	if err == nil {
		db = db.Where("rasp_attacks.is_blocked = ?", isBlocked)
	}
	// 处理结果
	result, err := strconv.ParseInt(req.HandleResult, 10, 32)
	if err == nil {
		db = db.Where("rasp_attacks.handle_result = ?", result)
	}

	db = db.Select("rasp_attacks.created_at," +
		"rasp_attacks.row_guid," +
		"rasp_attacks.host_name," +
		"rasp_attacks.host_ip," +
		"rasp_attacks.remote_ip," +
		"rasp_attacks.attack_type," +
		"rasp_attacks.is_blocked," +
		"rasp_attacks.level," +
		"rasp_attacks.handle_result," +
		"rasp_attacks.request_uri," +
		"rasp_attacks.attack_time," +
		"rasp_attack_details.parent_guid," +
		"rasp_attack_details.context," +
		"rasp_attack_details.app_name," +
		"rasp_attack_details.stack_trace," +
		"rasp_attack_details.payload," +
		"rasp_attack_details.algorithm," +
		"rasp_attack_details.extend," +
		"rasp_attack_details.meta_info").
		Joins("right join rasp_attack_details on rasp_attack_details.parent_guid=rasp_attacks.row_guid")

	//记录总条数
	var total int64
	err = db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	err = db.Find(&list).Error
	return list, total, err
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
