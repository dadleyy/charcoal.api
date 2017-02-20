package db

import "github.com/jinzhu/gorm"

type Connection struct {
	*gorm.DB
}

func Open(config Config, connection *Connection) error {
	conn, err := gorm.Open("mysql", config.String())

	if err != nil {
		return err
	}

	conn.LogMode(config.Debug == true)
	*connection = Connection{conn}
	return nil
}
