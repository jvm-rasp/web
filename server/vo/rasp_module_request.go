package vo

import "gorm.io/datatypes"

type CreateRaspModuleRequest struct {
	ModuleName    string         `json:"moduleName" form:"moduleName" validate:"required,min=2,max=32"`
	ModuleVersion string         `json:"moduleVersion" form:"moduleVersion" validate:"required,min=2,max=32"`
	ModuleType    uint           `json:"moduleType" form:"moduleType" validate:"required,min=1,max=3"`
	DownLoadURL   string         `json:"downLoadURL" form:"downLoadURL" validate:"required,min=2,max=1024"`
	Md5           string         `json:"md5" form:"md5" validate:"required,min=2,max=256"`
	Parameters    datatypes.JSON `json:"parameters" form:"parameters" validate:"required"`
	Status        bool           `json:"status" form:"status" validate:"boolean"`
	Desc          string         `json:"desc" form:"desc" validate:"min=2,max=100"`
}

type RaspModuleListRequest struct {
	Name     string `json:"name" form:"name"`
	Status   string `json:"status" form:"status"`
	Creator  string `json:"creator" form:"creator"`
	PageNum  uint   `json:"pageNum" form:"pageNum"`
	PageSize uint   `json:"pageSize" form:"pageSize"`
}

type UpdateRaspModuleRequest struct {
	ID            uint           `json:"id" form:"id" validate:"required,min=1,max=32"`
	ModuleName    string         `json:"moduleName" form:"moduleName" validate:"required,min=2,max=32"`
	ModuleVersion string         `json:"moduleVersion" form:"moduleVersion" validate:"required,min=2,max=32"`
	ModuleType    uint           `json:"moduleType" form:"moduleType" validate:"required,min=1,max=3"`
	DownLoadURL   string         `json:"downLoadURL" form:"downLoadURL" validate:"required,min=2,max=1024"`
	Md5           string         `json:"md5" form:"md5" validate:"required,min=2,max=256"`
	Parameters    datatypes.JSON `json:"parameters" form:"parameters" validate:"required"`
	Status        bool           `json:"status" form:"status" validate:"boolean"`
	Desc          string         `json:"desc" form:"desc" validate:"min=2,max=100"`
}

type DeleteRaspModuleRequest struct {
	Id uint `json:"id" form:"id"`
}

type DeleteBatchRaspModuleRequest struct {
	Ids []uint `json:"ids" form:"ids"`
}

type UpdateRaspModuleStatusRequest struct {
	ID uint `json:"id" form:"id" validate:"required"`
}

type UpgradeRaspModuleRequest struct {
	ID uint `json:"id" form:"id" validate:"required"`
}
