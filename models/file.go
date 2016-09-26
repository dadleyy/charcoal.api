package models

import "github.com/jinzhu/gorm"

type File struct {
	gorm.Model
	Key string `json:"key"`
}

func (file *File) Marshal() interface{} {
	return file
}
