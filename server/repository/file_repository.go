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
	GetRaspFileById(id uint) (*model.RaspFile, error)
	GetRaspFileByName(fileName string) (*model.RaspFile, int64, error)
	GetRaspFileByHash(fileHash string) (*model.RaspFile, error)
	CreateRaspFile(files *model.RaspFile) error
	DeleteRaspFile(ids []uint) error
	UpdateRaspFile(raspHost *model.RaspFile) error
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
	// 创建人模糊查询
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		db = db.Where("creator LIKE ?", fmt.Sprintf("%%%s%%", creator))
	}
	// hash模糊查询
	hash := strings.TrimSpace(req.FileHash)
	if hash != "" {
		db = db.Where("file_hash LIKE ?", fmt.Sprintf("%%%s%%", hash))
	}
	// 文件类型模糊查询
	mimeType := strings.TrimSpace(req.MimeType)
	if mimeType != "" {
		db = db.Where("mime_type LIKE ?", fmt.Sprintf("%%%s%%", mimeType))
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

func (h RaspFileRepository) GetRaspFileById(id uint) (*model.RaspFile, error) {
	var record *model.RaspFile
	err := common.DB.Find(&record, "id = ?", id).Error
	return record, err
}

func (h RaspFileRepository) GetRaspFileByName(fileName string) (*model.RaspFile, int64, error) {
	var record *model.RaspFile
	var count int64
	err := common.DB.Find(&record, "file_name = ?", fileName).Count(&count).Error
	return record, count, err
}

func (h RaspFileRepository) GetRaspFileByHash(fileHash string) (*model.RaspFile, error) {
	var record *model.RaspFile
	err := common.DB.Find(&record, "file_hash = ?", fileHash).Error
	return record, err
}

func (h RaspFileRepository) CreateRaspFile(raspHost *model.RaspFile) error {
	err := common.DB.Create(raspHost).Error
	return err
}

func (h RaspFileRepository) DeleteRaspFile(ids []uint) error {
	err := common.DB.Where("id IN (?)", ids).Unscoped().Delete(&model.RaspFile{}).Error
	return err
}

func (h RaspFileRepository) UpdateRaspFile(raspHost *model.RaspFile) error {
	err := common.DB.Save(raspHost).Error
	return err
}
