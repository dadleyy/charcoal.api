package testutils

import "os"
import "fmt"

import _ "github.com/jinzhu/gorm/dialects/mysql"
import "github.com/jinzhu/gorm"
import "github.com/joho/godotenv"
import "github.com/go-sql-driver/mysql"

var connection *gorm.DB

func NewDB() *gorm.DB {
	if connection != nil {
		return connection
	}

	godotenv.Load("../.env")

	dbc := mysql.Config{
		User:                    os.Getenv("DB_USERNAME"),
		Passwd:                  os.Getenv("DB_PASSWORD"),
		Net:                     "tcp",
		Addr:                    fmt.Sprintf("%s:%s", os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT")),
		DBName:                  os.Getenv("DB_DATABASE"),
		AllowCleartextPasswords: true,
		ParseTime:               true,
	}

	db, err := gorm.Open("mysql", dbc.FormatDSN())

	if err != nil {
		panic(err)
	}

	connection = db

	return connection
}
