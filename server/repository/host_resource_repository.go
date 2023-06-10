package repository

import (
	"gorm.io/gorm"
	"server/common"
	"server/model"
)

type IHostResourceRepository interface {
	CreateResource(resource *model.HostResource) error
	GetResourceByName(hostName string) (*model.HostResource, error)
	GetResourceByNameAndIP(hostName string, ip string) (*model.HostResource, error)
}

type HostResourceRepository struct {
}

func NewHostResourceRepository() IHostResourceRepository {
	return HostResourceRepository{}
}

func (this HostResourceRepository) CreateResource(resource *model.HostResource) error {
	err := common.DB.Create(resource).Error
	return err
}

func (this HostResourceRepository) GetResourceByName(hostName string) (*model.HostResource, error) {
	var record *model.HostResource
	err := common.DB.First(&record, "host_name = ?", hostName).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return record, err
}

func (this HostResourceRepository) GetResourceByNameAndIP(hostName string, ip string) (*model.HostResource, error) {
	var record *model.HostResource
	err := common.DB.Where("host_name = ?", hostName).Where("ip = ?", ip).First(&record).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return record, err
}
