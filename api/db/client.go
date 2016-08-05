package db

import "os"
import "fmt"
import "errors"
import "github.com/jinzhu/gorm"

const DSN_STR = "%v:%v@tcp(%v:%v)/%v?parseTime=true"

type Client struct {
	*gorm.DB
}

func Get() (Client, error) {
	// get configuration information from the environment
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	hostname := os.Getenv("DB_HOSTNAME")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_DATABASE")

	// build our data source url
	dsn := fmt.Sprintf(DSN_STR, username, password, hostname, port, database)

	// attempt to connect to the mysql database
	db, err := gorm.Open("mysql", dsn)

	if err != nil {
		return Client{}, errors.New("BAD_CONNECTION")
	}

	// turn off gorm logging
	db.LogMode(false)

	return Client{db}, nil
}
