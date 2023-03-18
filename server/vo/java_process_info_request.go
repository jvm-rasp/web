package vo

type JavaProcessInfoListRequest struct {
	HostName string `json:"name" form:"hostName"`
	Status   string   `json:"status" form:"status"`
	PageNum  uint   `json:"pageNum" form:"pageNum"`
	PageSize uint   `json:"pageSize" form:"pageSize"`
}
