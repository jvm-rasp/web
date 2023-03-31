package repository

import (
	"fmt"
	"server/common"
	"server/model"
	"server/vo"
	"strings"
)

type IRaspFileRepository interface {
	GetRaspFiles(req *vo.RaspFileListRequest) ([]*model.RaspFile, int64, error)
	CreateRaspFile(files *model.RaspFile) error
	DeleteRaspFile(ids []uint) error
}

type RaspFileRepository struct {
}

func NewRaspFileRepository() IRaspFileRepository {
	return RaspFileRepository{}
}

func (h RaspFileRepository) GetRaspFiles(req *vo.RaspFileListRequest) ([]*model.RaspFile, int64, error) {
	var list []*model.RaspFile
	db := common.DB.Model(&model.RaspFile{}).Order("created_at DESC")

	// 名称模糊查询
	name := strings.TrimSpace(req.FileName)
	if name != "" {
		db = db.Where("file_name LIKE ?", fmt.Sprintf("%%%s%%", name))
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

func (h RaspFileRepository) CreateRaspFile(raspHost *model.RaspFile) error {
	err := common.DB.Create(raspHost).Error
	return err
}

func (h RaspFileRepository) DeleteRaspFile(ids []uint) error {
	err := common.DB.Where("id IN (?)", ids).Unscoped().Delete(&model.RaspFile{}).Error
	return err
}
