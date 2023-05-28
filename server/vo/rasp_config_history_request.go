package vo

type RaspConfigHistoryListRequest struct {
	ConfigGuid string `json:"configGuid" form:"configGuid"`
}

type RaspConfigHistoryDataRequest struct {
	ConfigGuid string `json:"configGuid" form:"configGuid"`
	Version    int    `json:"version" form:"version"`
}
