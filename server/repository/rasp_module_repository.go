package repository

import (
	"fmt"
	"server/common"
	"server/model"
	"server/vo"
	"strconv"
	"strings"
)

type IRaspModuleRepository interface {
	GetRaspModules(req *vo.RaspModuleListRequest) ([]*model.RaspModule, int64, error)
	GetRaspModuleById(id uint) (*model.RaspModule, error)
	UpdateRaspModule(module *model.RaspModule) error
	CreateRaspModule(module *model.RaspModule) error
	DeleteRaspModule(ids []uint) error
}

type RaspModuleRepository struct {
}

func NewRaspModuleRepository() IRaspModuleRepository {
	return RaspModuleRepository{}
}

func (a RaspModuleRepository) GetRaspModules(req *vo.RaspModuleListRequest) ([]*model.RaspModule, int64, error) {
	var list []*model.RaspModule
	db := common.DB.Model(&model.RaspModule{}).Order("created_at DESC")

	// 名称模糊查询
	name := strings.TrimSpace(req.Name)
	if name != "" {
		db = db.Where("module_name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}

	// 配置的状态
	status, err := strconv.ParseBool(req.Status)
	if err == nil {
		db = db.Where("status = ?", status)
	}

	//创建者
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		db = db.Where("creator = ?", creator)
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

func (a RaspModuleRepository) GetRaspModuleById(id uint) (*model.RaspModule, error) {
	var record *model.RaspModule
	err := common.DB.Find(&record, "id = ?", id).Error
	return record, err
}

func (a RaspModuleRepository) UpdateRaspModule(module *model.RaspModule) error {
	err := common.DB.Save(module).Error
	return err
}

func (a RaspModuleRepository) CreateRaspModule(module *model.RaspModule) error {
	err := common.DB.Create(module).Error
	return err
}

func (r RaspModuleRepository) DeleteRaspModule(ids []uint) error {
	err := common.DB.Where("id IN (?)", ids).Unscoped().Delete(&model.RaspModule{}).Error
	return err
}
