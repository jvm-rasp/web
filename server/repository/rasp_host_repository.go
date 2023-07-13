package repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"server/common"
	"server/model"
	"server/util"
	"server/vo"
	"strings"
)

type IRaspHostRepository interface {
	GetRaspHosts(req *vo.RaspHostListRequest) ([]*model.RaspHost, int64, error)
	GetRaspHostList() ([]*model.RaspHost, int64, error)
	GetRaspHostsByConfigId(id uint) ([]*model.RaspHost, int64, error)
	CreateRaspHost(host *model.RaspHost) (uint, error)
	DeleteRaspHost(ids []uint) error
	DeleteRaspHostById(id uint) error
	QueryRaspHost(hostName string) ([]*model.RaspHost, error)
	UpdateRaspHostByHostName(host *model.RaspHost) error
	UpdateRaspHost(host *model.RaspHost) error
	GetRaspHostById(id uint) (*model.RaspHost, error)
	GetRaspHostByHostName(hostName string) (*model.RaspHost, error)
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

	// 保护状态查询
	status := req.Status
	if status != "" {
		switch status {
		case "0":
			db = db.Where("not_inject > 0")
		case "1":
			db = db.Where("success_inject > 0")
		case "2":
			db = db.Where("failed_inject > 0")
		}
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

func (h RaspHostRepository) GetRaspHostList() ([]*model.RaspHost, int64, error) {
	var list []*model.RaspHost
	db := common.DB.Model(&model.RaspHost{}).Order("created_at DESC")
	var total int64
	err := db.Count(&total).Find(&list).Error
	return list, total, err
}

func (h RaspHostRepository) GetRaspHostsByConfigId(id uint) ([]*model.RaspHost, int64, error) {
	var list []*model.RaspHost
	db := common.DB.Model(&model.RaspHost{}).Order("created_at DESC")
	var total int64
	err := db.Where("config_id = ?", id).Count(&total).Find(&list).Error
	return list, total, err
}

func (h RaspHostRepository) CreateRaspHost(raspHost *model.RaspHost) (uint, error) {
	if raspHost.HostName == "" || raspHost.Ip == "" {
		common.Log.Errorf("hostInfo信息为空, stacktrace: %v", util.GetCallers())
	}
	// 先获取默认策略
	var list []*model.RaspConfig
	result := common.DB.Model(&model.RaspConfig{}).Where("is_default = ?", true).Find(&list)
	if result.Error != nil {
		return 0, result.Error
	}
	if len(list) > 0 {
		record := list[0]
		raspHost.ConfigId = record.ID
	}
	err := common.DB.Create(raspHost).Error
	return raspHost.ConfigId, err
}

func (h RaspHostRepository) DeleteRaspHost(ids []uint) error {
	err := common.DB.Where("id IN (?)", ids).Unscoped().Delete(&model.RaspHost{}).Error
	return err
}

func (h RaspHostRepository) DeleteRaspHostById(id uint) error {
	err := common.DB.Where("id = ?", id).Unscoped().Delete(&model.RaspHost{}).Error
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
	err := common.DB.Save(host).Error
	return err
}

func (h RaspHostRepository) GetRaspHostById(id uint) (*model.RaspHost, error) {
	var host *model.RaspHost
	err := common.DB.Where("id = ?", id).First(&host).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return host, err
}

func (h RaspHostRepository) GetRaspHostByHostName(hostName string) (*model.RaspHost, error) {
	var host *model.RaspHost
	err := common.DB.Where("host_name = ?", hostName).First(&host).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return host, err
}
