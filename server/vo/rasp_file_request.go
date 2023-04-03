package vo

type RaspFileListRequest struct {
	ModuleName string `json:"moduleName" form:"moduleName"`
	FileHash   string `json:"fileHash" form:"fileHash"`
	Creator    string `json:"creator" form:"creator"`
	PageNum    uint   `json:"pageNum" form:"pageNum"`
	PageSize   uint   `json:"pageSize" form:"pageSize"`
}

type RaspFileDeleteRequest struct {
	Ids []uint `json:"ids" form:"ids"`
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
