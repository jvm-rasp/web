package vo

type CreateRaspModuleRequest struct {
	Name              string `json:"name" form:"name" validate:"required,min=2,max=32"`
	Version           string `json:"version" form:"version" validate:"required,min=2,max=32"`
	Url               string `json:"url" form:"url" validate:"required,min=2,max=1024"`
	Hash              string `json:"hash" form:"hash" validate:"required,min=2,max=256"`
	Type              string `json:"type" form:"type" validate:"required,min=2,max=256"`
	MiddlewareName    string `json:"middlewareName" form:"middlewareName" validate:"required,min=2,max=32"`
	MiddlewareVersion string `json:"middlewareVersion" form:"middlewareVersion" validate:"required,min=2,max=32"`
	Tag               string `json:"tag" form:"tag" validate:"required,min=2,max=32"`
	Parameters        string `json:"parameters" form:"parameters" validate:"required,min=2,max=4096"`
	Status            uint   `json:"status" form:"status" validate:"oneof=0 1"`
	Desc              string `json:"desc" form:"desc" validate:"min=2,max=100"`
}

type RaspModuleListRequest struct {
	Name     string `json:"name" form:"name"`
	Status   uint   `json:"status" form:"status"`
	Creator  string `json:"creator" form:"creator"`
	PageNum  uint   `json:"pageNum" form:"pageNum"`
	PageSize uint   `json:"pageSize" form:"pageSize"`
}

type DeleteRaspModuleRequest struct {
	Ids []uint `json:"ids" form:"ids"`
}
