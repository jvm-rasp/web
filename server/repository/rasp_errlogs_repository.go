package repository

import (
	"fmt"
	"server/common"
	"server/model"
	"server/vo"
	"strings"
)

type IRaspErrorLogsRepository interface {
	GetRaspLogs(req *vo.RaspLogsListRequest) ([]*model.RaspErrorLogs, int64, error)
	CreateRaspLogs(errorLogs *model.RaspErrorLogs) error
	DeleteRaspLogs(ids []uint) error
}

type RaspErrorLogsRepository struct {
}

func NewRaspErrorLogsRepository() IRaspErrorLogsRepository {
	return RaspErrorLogsRepository{}
}

func (r RaspErrorLogsRepository) GetRaspLogs(req *vo.RaspLogsListRequest) ([]*model.RaspErrorLogs, int64, error) {
	var list []*model.RaspErrorLogs
	db := common.DB.Model(&model.RaspErrorLogs{}).Order("time DESC")

	// topic模糊查询
	topic := strings.TrimSpace(req.Topic)
	if topic != "" {
		db = db.Where("topic LIKE ?", fmt.Sprintf("%%%s%%", topic))
	}

	// level模糊查询
	level := strings.TrimSpace(req.Level)
	if level != "" {
		db = db.Where("level LIKE ?", fmt.Sprintf("%%%s%%", level))
	}
	// 时间范围查询
	startDate := strings.TrimSpace(req.StartDate)
	endDate := strings.TrimSpace(req.EndDate)
	if startDate != "" && endDate != "" {
		db = db.Where("time >= ? and time <= ?", startDate+" 00:00:00", endDate+" 23:59:59")
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

func (r RaspErrorLogsRepository) CreateRaspLogs(errorLogs *model.RaspErrorLogs) error {
	err := common.DB.Create(errorLogs).Error
	return err
}

func (r RaspErrorLogsRepository) DeleteRaspLogs(ids []uint) error {
	err := common.DB.Where("id IN (?)", ids).Unscoped().Delete(&model.RaspErrorLogs{}).Error
	return err
}
