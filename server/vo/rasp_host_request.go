package vo

type RaspHostListRequest struct {
	HostName string `json:"name" form:"hostName"`
	Status   uint   `json:"status" form:"status"`
	PageNum  uint   `json:"pageNum" form:"pageNum"`
	PageSize uint   `json:"pageSize" form:"pageSize"`
}

type DeleteRaspHostRequest struct {
	Ids []uint `json:"ids" form:"ids"`
}
