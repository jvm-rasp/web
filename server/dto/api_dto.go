package dto

import "server/model"

type ApiTreeDto struct {
	ID       int          `json:"ID"`
	Desc     string       `json:"desc"`
	Category string       `json:"category"`
	Children []*model.Api `json:"children"`
}
