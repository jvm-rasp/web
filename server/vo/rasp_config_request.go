package vo

// CreateRaspConfigRequest 创建接口结构体
type CreateRaspConfigRequest struct {
	Name     string `json:"name" form:"name" validate:"required,min=2,max=32"`
	Content  string `json:"content" form:"content" validate:"required"`
	Tag      string `json:"tag" form:"tag"`
	Status   uint   `json:"status" form:"status" validate:"oneof=0 1"`
	Desc     string `json:"desc" form:"desc" validate:"min=2,max=100"`
}

// 获取接口列表结构体
type RaspConfigListRequest struct {
	Name     string `json:"name" form:"name"`
	Status   uint   `json:"status" form:"status"`
	Creator  string `json:"creator" form:"creator"`
	PageNum  uint   `json:"pageNum" form:"pageNum"`
	PageSize uint   `json:"pageSize" form:"pageSize"`
}

// 批量删除接口结构体
type DeleteRaspConfigRequest struct {
	Ids []uint `json:"ids" form:"ids"`
}
