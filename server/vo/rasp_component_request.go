package vo

type RaspComponentListRequest struct {
	ParentGuid string `json:"parentGuid" form:"parentGuid" validate:"required"`
}

type RaspComponentDeleteRequest struct {
	Id uint `json:"id" form:"id"`
}
