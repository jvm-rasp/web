package repository

import (
	"errors"
	"fmt"
	"server/common"
	"server/model"
	"server/vo"
	"strings"
)

type IRaspHostRepository interface {
	GetRaspHosts(req *vo.RaspHostListRequest) ([]*model.RaspHost, int64, error)
	CreateRaspHost(host *model.RaspHost) error
	DeleteRaspHost(ids []uint) error
	QueryRaspHost(hostName string) ([]*model.RaspHost, error)
	UpdateRaspHostByHostName(host *model.RaspHost) error
	UpdateRaspHost(host *model.RaspHost) error
	GetRaspHostById(id uint) (*model.RaspHost, error)
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
		db = db.Where("host_name LIKE ?", fmt.Sprintf("%%%s%%", name))
	} else {
		// ip 模糊查询
		ip := strings.TrimSpace(req.Ip)
		if ip != "" {
			db = db.Where("ip LIKE ?", fmt.Sprintf("%%%s%%", ip))
		}
	}

	// agent 模式查询
	agentMode := strings.TrimSpace(req.AgentMode)
	if agentMode != "" {
		db = db.Where("agent_mode = ?", agentMode)
	}

	// todo 在线状态查询
	//status := req.Status
	//if status != 0 {
	//	db = db.Where("status = ?", status)
	//}

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

func (h RaspHostRepository) QueryRaspHost(hostName string) ([]*model.RaspHost, error) {
	var list []*model.RaspHost
	db := common.DB.Model(&model.RaspHost{}).Order("created_at DESC")
	name := strings.TrimSpace(hostName)
	if name != "" {
		db = db.Where("host_name = ?", name)
	}
	err := db.Find(&list).Limit(1).Error
	if err != nil {
		return list, err
	}
	return list, nil
}

// UpdateRaspHostByHostName 通过id更新节点信息
func (h RaspHostRepository) UpdateRaspHostByHostName(host *model.RaspHost) error {
	if host == nil || host.HostName == "" {
		return errors.New("host object is nil or host_name is nil")
	}
	err := common.DB.Model(host).Where("host_name = ?", host.HostName).Updates(host).Error
	return err
}

func (h RaspHostRepository) UpdateRaspHost(host *model.RaspHost) error {
	err := common.DB.Model(host).Updates(host).Error
	return err
}

func (h RaspHostRepository) GetRaspHostById(id uint) (*model.RaspHost, error) {
	var host *model.RaspHost
	err := common.DB.Model(&model.RaspHost{}).Where("id = ?", id).Find(&host).Error
	return host, err
}
