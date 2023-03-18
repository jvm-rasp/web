package repository

import (
	"fmt"
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
	DeleteProcessByPid(hostName string, pid uint) error
}

type JavaProcessInfoRepository struct {
}

func NewJavaProcessInfoRepository() IJavaProcessInfoRepository {
	return JavaProcessInfoRepository{}
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
	err := common.DB.Where("id = ?", id).Unscoped().Delete(&model.JavaProcessInfo{}).Error
	return err
}

func (j JavaProcessInfoRepository) SaveProcessInfo(process *model.JavaProcessInfo) error {
	err := common.DB.Create(process).Error
	return err
}

func (j JavaProcessInfoRepository) DeleteProcessByPid(hostName string, pid uint) error {
	err := common.DB.Where("host_name = ?", hostName).
		Where("pid = ?", pid).Unscoped().Delete(&model.JavaProcessInfo{}).Error
	return err
}
