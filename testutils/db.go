package testutils

import "os"
import "github.com/jinzhu/gorm"
import "github.com/joho/godotenv"
import _ "github.com/jinzhu/gorm/dialects/mysql"

import "github.com/dadleyy/charcoal.api/db"

func NewDB() *gorm.DB {
	_ = godotenv.Load("../.env")

	dbconf := db.Config{
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOSTNAME"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DEBUG") == "true",
	}

	database, err := gorm.Open("mysql", dbconf.String())

	if err != nil {
		panic(err)
	}

	return database
}
