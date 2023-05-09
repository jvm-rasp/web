package repository

import (
	"fmt"
	"gorm.io/gorm"
	"server/common"
	"server/model"
	"server/vo"
	"strings"
)

type IJavaProcessInfoRepository interface {
	GetJavaProcessInfos(req *vo.JavaProcessInfoListRequest) ([]*model.JavaProcessInfo, int64, error)
	GetAllJavaProcessInfos(hostName string) ([]*model.JavaProcessInfo, error)
	DeleteProcess(id uint) error
	SaveProcessInfo(*model.JavaProcessInfo) error
	UpdateProcessInfo(info *model.JavaProcessInfo) error
	DeleteProcessByPid(hostName string, pid uint) error
	GetProcessByPid(hostName string, pid uint) (*model.JavaProcessInfo, error)
	GetSuccessInjectCount(hostName string) (int64, error)
	GetFailedInjectCount(hostName string) (int64, error)
	GetNotInjectCount(hostName string) (int64, error)
	UpdateRaspHostInjectCounts(hostName string) error
}

type JavaProcessInfoRepository struct {
	RaspHostRepository IRaspHostRepository
}

func NewJavaProcessInfoRepository(raspHostRepository IRaspHostRepository) IJavaProcessInfoRepository {
	return JavaProcessInfoRepository{
		RaspHostRepository: raspHostRepository,
	}
}

func (j JavaProcessInfoRepository) GetJavaProcessInfos(req *vo.JavaProcessInfoListRequest) ([]*model.JavaProcessInfo, int64, error) {
	var list []*model.JavaProcessInfo
	db := common.DB.Model(&model.JavaProcessInfo{}).Order("created_at DESC")

	// 名称模糊查询
	name := strings.TrimSpace(req.HostName)
	if name != "" {
		db = db.Where("host_name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}

	// 进程注入状态
	status := strings.TrimSpace(req.Status)
	if status != "" {
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

func (j JavaProcessInfoRepository) GetAllJavaProcessInfos(hostName string) ([]*model.JavaProcessInfo, error) {
	var list []*model.JavaProcessInfo
	db := common.DB.Model(&model.JavaProcessInfo{}).Order("created_at DESC")

	// 名称模糊查询
	name := strings.TrimSpace(hostName)
	if name != "" {
		db = db.Where("host_name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	err := db.Find(&list).Error
	return list, err
}

func (j JavaProcessInfoRepository) DeleteProcess(id uint) error {
	var javaProcessInfo model.JavaProcessInfo
	err := common.DB.Where("id = ?", id).Find(&javaProcessInfo).Error
	if err != nil {
		return err
	}
	hostName := javaProcessInfo.HostName
	err = common.DB.Model(&model.JavaProcessInfo{}).Unscoped().Delete(javaProcessInfo).Error
	if err != nil {
		return err
	}
	err = j.UpdateRaspHostInjectCounts(hostName)
	return err
}

func (j JavaProcessInfoRepository) SaveProcessInfo(process *model.JavaProcessInfo) error {
	err := common.DB.Create(process).Error
	if err != nil {
		return err
	}
	err = j.UpdateRaspHostInjectCounts(process.HostName)
	return err
}

func (j JavaProcessInfoRepository) UpdateProcessInfo(info *model.JavaProcessInfo) error {
	err := common.DB.Save(info).Error
	if err != nil {
		return err
	}
	err = j.UpdateRaspHostInjectCounts(info.HostName)
	return err
}

func (j JavaProcessInfoRepository) DeleteProcessByPid(hostName string, pid uint) error {
	javaProcessInfo, err := j.GetProcessByPid(hostName, pid)
	if err != nil {
		return err
	}
	if javaProcessInfo != nil {
		hostName = javaProcessInfo.HostName
		err = common.DB.Unscoped().Delete(javaProcessInfo).Error
		if err != nil {
			return err
		}
		err = j.UpdateRaspHostInjectCounts(hostName)
		return err
	} else {
		return nil
	}
}

func (j JavaProcessInfoRepository) GetProcessByPid(hostName string, pid uint) (*model.JavaProcessInfo, error) {
	var javaProcessInfo model.JavaProcessInfo
	err := common.DB.Model(&model.JavaProcessInfo{}).Where("host_name = ?", hostName).Where("pid = ?", pid).First(&javaProcessInfo).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &javaProcessInfo, err
}

func (j JavaProcessInfoRepository) GetSuccessInjectCount(hostName string) (int64, error) {
	var count int64
	err := common.DB.Model(&model.JavaProcessInfo{}).Where("host_name = ?", hostName).Where("status = ?", 1).Count(&count).Error
	return count, err
}

func (j JavaProcessInfoRepository) GetFailedInjectCount(hostName string) (int64, error) {
	var count int64
	err := common.DB.Model(&model.JavaProcessInfo{}).Where("host_name = ?", hostName).Where("status = ?", 2).Count(&count).Error
	return count, err
}

func (j JavaProcessInfoRepository) GetNotInjectCount(hostName string) (int64, error) {
	var count int64
	err := common.DB.Model(&model.JavaProcessInfo{}).Where("host_name = ?", hostName).Where("status = ?", 0).Count(&count).Error
	return count, err
}

func (j JavaProcessInfoRepository) UpdateRaspHostInjectCounts(hostName string) error {
	// 更新保护应用数量
	raspHost, err := j.RaspHostRepository.GetRaspHostByHostName(hostName)
	if err != nil {
		return err
	}
	raspHost.NotInject, err = j.GetNotInjectCount(hostName)
	if err != nil {
		return err
	}
	raspHost.SuccessInject, err = j.GetSuccessInjectCount(hostName)
	if err != nil {
		return err
	}
	raspHost.FailedInject, err = j.GetFailedInjectCount(hostName)
	if err != nil {
		return err
	}
	err = j.RaspHostRepository.UpdateRaspHost(raspHost)
	if err != nil {
		return err
	}
	return nil
}
