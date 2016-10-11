package db

import "github.com/jinzhu/gorm"

type Connection struct {
	*gorm.DB
}

func Open(config Config) (*Connection, error) {
	conn, err := gorm.Open("mysql", config.String())

	if err != nil {
		return nil, err
	}

	return &Connection{conn}, nil
}
