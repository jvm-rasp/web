package repository

import (
	"fmt"
	"gorm.io/gorm"
	"server/common"
	"server/model"
	"server/vo"
	"strconv"
	"strings"
)

type IRaspModuleRepository interface {
	GetRaspModules(req *vo.RaspModuleListRequest) ([]*model.RaspModule, int64, error)
	GetRaspModuleById(id uint) (*model.RaspModule, error)
	GetRaspModuleByNameAndVersion(moduleName string, moduleVersion string) (*model.RaspModule, error)
	GetRaspModuleByName(moduleName string) (*model.RaspModule, error)
	GetRaspModuleByComponentName(componentName string) (*model.RaspModule, error)
	UpdateRaspModule(module *model.RaspModule) error
	CreateRaspModule(module *model.RaspModule) error
	DeleteRaspModule(ids []uint) error
}

type RaspModuleRepository struct {
	RaspComponentRepository IRaspComponentRepository
}

func NewRaspModuleRepository() IRaspModuleRepository {
	repo1 := NewRaspComponentRepository()
	return RaspModuleRepository{
		RaspComponentRepository: repo1,
	}
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

func (a RaspModuleRepository) GetRaspModuleByNameAndVersion(moduleName string, moduleVersion string) (*model.RaspModule, error) {
	var record *model.RaspModule
	err := common.DB.Where("module_name = ?", moduleName).Where("module_version = ?", moduleVersion).First(&record).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return record, err
}

func (a RaspModuleRepository) GetRaspModuleByName(moduleName string) (*model.RaspModule, error) {
	var record *model.RaspModule
	err := common.DB.Where("module_name = ?", moduleName).First(&record).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return record, err
}

func (a RaspModuleRepository) GetRaspModuleByComponentName(componentName string) (*model.RaspModule, error) {
	var record *model.RaspModule
	componentInfo, err := a.RaspComponentRepository.GetRaspComponentByName(componentName)
	if err != nil {
		return nil, err
	}
	if componentInfo == nil {
		return nil, nil
	}
	err = common.DB.Where("row_guid = ?", componentInfo.ParentGuid).First(&record).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
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
