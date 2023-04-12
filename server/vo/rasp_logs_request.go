package vo

type RaspLogsListRequest struct {
	Topic     string `json:"topic" form:"topic"`
	StartDate string `json:"startDate" form:"startDate"`
	EndDate   string `json:"endDate" form:"endDate"`
	Level     string `json:"level" form:"level"`
	PageNum   uint   `json:"pageNum" form:"pageNum"`
	PageSize  uint   `json:"pageSize" form:"pageSize"`
}

// RaspLogDeleteRequest 批量删除接口结构体
type RaspLogsDeleteRequest struct {
	Ids []uint `json:"ids" form:"ids"`
}
