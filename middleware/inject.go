package middleware

import "os"
import "fmt"
import "github.com/jinzhu/gorm"
import "github.com/labstack/echo"
import _ "github.com/jinzhu/gorm/dialects/mysql"

import "github.com/sizethree/miritos.api/context"
import "github.com/sizethree/miritos.api/filestore"

const DSN_STR = "%v:%v@tcp(%v:%v)/%v?parseTime=true"

func Inject(handler echo.HandlerFunc) echo.HandlerFunc {
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	hostname := os.Getenv("DB_HOSTNAME")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_DATABASE")
	dsn := fmt.Sprintf(DSN_STR, username, password, hostname, port, database)

	inject := func(ctx echo.Context) error {
		db, err := gorm.Open("mysql", dsn)

		if err != nil {
			ctx.Logger().Error(err)
			return err
		}

		errors  := make(context.ErrorList, 0)
		meta    := make(context.MetaData)
		results := make(context.ResultList, 0)

		client := context.Database{db}
		var store context.FileSaver

		switch os.Getenv("FS_ENGINE") {
		case "s3":
			store = filestore.S3FileStore{"123", "456", "789"}
		default:
			store = filestore.TempStore{}
		}

		app := &context.Miritos{ctx, &client, errors, meta, results, store}

		result := handler(app)

		if result == nil {
			return app.Finish()
		}

		return result
	}

	return inject
}
