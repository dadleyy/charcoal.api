package server

import "github.com/jinzhu/gorm"

type Database struct {
	*gorm.DB
}
