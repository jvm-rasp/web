package repository

import (
	"fmt"
	"server/common"
	"server/model"
	"server/vo"
	"strings"
)

type IRaspHostRepository interface {
	GetRaspHosts(req *vo.RaspHostListRequest) ([]*model.RaspHost, int64, error)
	CreateRaspHost(config *model.RaspHost) error
	DeleteRaspHost(ids []uint) error
}

type RaspHostRepository struct {
}

func NewRaspHostRepository() IRaspHostRepository {
	return RaspHostRepository{}
}

func (h RaspHostRepository) GetRaspHosts(req *vo.RaspHostListRequest) ([]*model.RaspHost, int64, error) {
	var list []*model.RaspHost
	db := common.DB.Model(&model.RaspHost{}).Order("created_at DESC")

	// 名称模糊查询
	name := strings.TrimSpace(req.HostName)
	if name != "" {
		db = db.Where("hostName LIKE ?", fmt.Sprintf("%%%s%%", name))
	}

	status := req.Status
	if status != 0 {
		db = db.Where("status = ?", status)
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

func (h RaspHostRepository) CreateRaspHost(raspHost *model.RaspHost) error {
	err := common.DB.Create(raspHost).Error
	return err
}

func (h RaspHostRepository) DeleteRaspHost(ids []uint) error {
	err := common.DB.Where("id IN (?)", ids).Unscoped().Delete(&model.RaspHost{}).Error
	return err
}
