package models

import "github.com/jinzhu/gorm"

type Photo struct {
	gorm.Model
}

func (photo *Photo) Marshal() interface{} {
	return photo
}
