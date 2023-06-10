package vo

type UpdateSettingRequest struct {
	Name  string      `json:"name" validate:"required"`
	Value interface{} `json:"value" validate:"required"`
}

type GetProjectInfoRequest struct {
	ReportUrl string `json:"reportUrl" validate:"required"`
}

type ProjectInfo struct {
	ProjectGuid string `json:"projectGuid"`
}
