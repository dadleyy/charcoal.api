package testutils

import "os"
import "github.com/jinzhu/gorm"
import "github.com/joho/godotenv"
import _ "github.com/jinzhu/gorm/dialects/mysql"

import "github.com/dadleyy/charcoal.api/db"

var connection *gorm.DB

func DBConfig() db.Config {
	_ = godotenv.Load("../.env")

	dbconf := db.Config{
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOSTNAME"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DEBUG") == "true",
	}

	return dbconf
}

func NewDB() *gorm.DB {
	if connection != nil {
		return connection
	}

	c := DBConfig()
	database, err := gorm.Open("mysql", c.String())

	if err != nil {
		panic(err)
	}

	connection = database

	return connection
}
