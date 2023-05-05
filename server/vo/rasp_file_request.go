package vo

type RaspFileListRequest struct {
	FileName string `json:"fileName" form:"fileName"`
	Creator  string `json:"creator" form:"creator"`
	FileHash string `json:"fileHash" form:"fileHash"`
	MimeType string `json:"mimeType" form:"mimeType"`
	PageNum  uint   `json:"pageNum" form:"pageNum"`
	PageSize uint   `json:"pageSize" form:"pageSize"`
}

type RaspFileDeleteRequest struct {
	Ids []uint `json:"ids" form:"ids"`
}

type RaspFileInfoRequest struct {
	Id uint `json:"id" form:"id"`
}

type RaspFileDownloadRequest struct {
	FileHash string `json:"hash" form:"hash"`
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
