package models

import "github.com/jinzhu/gorm"

type Activity struct {
	gorm.Model
}

func (item *Activity) Marshal() interface{} {
	return item
}
