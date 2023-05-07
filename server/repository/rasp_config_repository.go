package repository

import (
	"errors"
	"fmt"
	"server/common"
	"server/model"
	"server/vo"
	"strconv"
	"strings"
)

type IRaspConfigRepository interface {
	GetRaspConfigs(req *vo.RaspConfigListRequest) ([]*model.RaspConfig, int64, error) // 获取配置列表列表
	GetRaspConfigById(id uint) (*model.RaspConfig, error)
	CreateRaspConfig(config *model.RaspConfig) error // 创建接口
	UpdateRaspConfig(config *model.RaspConfig) error
	DeleteRaspConfig(ids []uint) error
	GetRaspConfig(hostName string) (*model.RaspConfig, error)
	GetRaspDefaultConfig() (*model.RaspConfig, error)
}

type RaspConfigRepository struct {
}

func NewRaspConfigRepository() IRaspConfigRepository {
	return RaspConfigRepository{}
}

// CreateRaspConfig 创建配置接口
func (a RaspConfigRepository) CreateRaspConfig(config *model.RaspConfig) error {
	err := common.DB.Create(config).Error
	return err
}

func (a RaspConfigRepository) UpdateRaspConfig(config *model.RaspConfig) error {
	err := common.DB.Save(config).Error
	return err
}

// GetRaspConfigs 查询配置
func (a RaspConfigRepository) GetRaspConfigs(req *vo.RaspConfigListRequest) ([]*model.RaspConfig, int64, error) {
	var list []*model.RaspConfig
	db := common.DB.Model(&model.RaspConfig{}).Order("created_at DESC")

	// 名称模糊查询
	name := strings.TrimSpace(req.Name)
	if name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
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

func (a RaspConfigRepository) GetRaspConfigById(id uint) (*model.RaspConfig, error) {
	var record *model.RaspConfig
	err := common.DB.Find(&record, "id = ?", id).Error
	return record, err
}

func (r RaspConfigRepository) DeleteRaspConfig(ids []uint) error {
	err := common.DB.Where("id IN (?)", ids).Unscoped().Delete(&model.RaspConfig{}).Error
	return err
}

func (a RaspConfigRepository) GetRaspConfig(hostName string) (*model.RaspConfig, error) {
	var list []*model.RaspConfig
	db := common.DB.Model(&model.RaspConfig{}).Order("created_at DESC")
	name := strings.TrimSpace(hostName)
	if name != "" {
		db = db.Where("name = ?", name)
	}

	err := db.Find(&list).Error
	if len(list) == 0 {
		return nil, errors.New("no config find in db, hostName: " + hostName)
	}
	return list[0], err
}

func (a RaspConfigRepository) GetRaspDefaultConfig() (*model.RaspConfig, error) {
	var record *model.RaspConfig
	err := common.DB.Find(&record, "is_default = ?", true).Error
	return record, err
}
