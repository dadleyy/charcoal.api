package models

type Photo struct {
	Common
	Label string `json:"label"`
	File uint `json:"file"`
}

func (photo *Photo) Marshal() interface{} {
	return photo
}
