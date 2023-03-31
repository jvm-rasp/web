package vo

type RaspFileListRequest struct {
	FileName string `json:"fileName" form:"fileName"`
	Status   uint `json:"status" form:"status"`
	PageNum  uint `json:"pageNum" form:"pageNum"`
	PageSize uint `json:"pageSize" form:"pageSize"`
}

type FileInfo struct {
	FileName      string
	FileHash      string
	DiskPath      string
	DownLoadUrl   string
	ModuleName    string
	ModuleVersion string
	UpdateTime    string
}
