package context

import "os"
import "fmt"
import "github.com/jinzhu/gorm"
import "github.com/labstack/echo"
import _ "github.com/jinzhu/gorm/dialects/mysql"

const DSN_STR = "%v:%v@tcp(%v:%v)/%v?parseTime=true"

type Miritos struct {
	echo.Context
	DB *gorm.DB
}

func New(echoContext echo.Context) (*Miritos, error) {
	var result *Miritos
	log := echoContext.Logger()

	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	hostname := os.Getenv("DB_HOSTNAME")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_DATABASE")

	dsn := fmt.Sprintf(DSN_STR, username, password, hostname, port, database)
	log.Infof("db connection: %s", dsn)

	db, err := gorm.Open("mysql", dsn)

	if err != nil {
		return result, err
	}

	result = &Miritos{echoContext, db}

	return result, nil
}
