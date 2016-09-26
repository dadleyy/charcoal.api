package models

import "github.com/jinzhu/gorm"

type Object struct {
	gorm.Model
}

func (item *Object) Marshal() interface{} {
	return item
}
