package api

import "github.com/jinzhu/gorm"
import _ "github.com/jinzhu/gorm/dialects/mysql"

type Runtime struct {
	DB *gorm.DB
}
