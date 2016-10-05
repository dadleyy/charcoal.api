package middleware

import "os"
import "github.com/jinzhu/gorm"
import "github.com/labstack/echo"
import _ "github.com/jinzhu/gorm/dialects/mysql"

import "github.com/sizethree/miritos.api/server"
import "github.com/sizethree/miritos.api/context"
import "github.com/sizethree/miritos.api/activity"
import "github.com/sizethree/miritos.api/filestore"

func Inject(activities chan activity.Message, dbconf server.DatabaseConfig) echo.MiddlewareFunc {
	injector := func(handler echo.HandlerFunc) echo.HandlerFunc {
		inject := func(ctx echo.Context) error {
			db, err := gorm.Open("mysql", dbconf.String())

			if err != nil {
				ctx.Logger().Error(err)
				return err
			}

			if dbconf.Debug == true {
				db.LogMode(true)
			}

			errors  := make(context.ErrorList, 0)
			meta    := make(context.MetaData)
			results := make(context.ResultList, 0)

			client := context.Database{db}
			var store context.FileSaver

			switch os.Getenv("FS_ENGINE") {
			case "s3":
				store = filestore.S3FileStore{}
			default:
				store = filestore.TempStore{}
			}

			app := &context.Runtime{
				Context: ctx,
				DB: &client,
				Errors: errors,
				Meta: meta,
				Results: results,
				FS: store,
				ActivityStream: activities,
			}

			result := handler(app)

			if result == nil {
				return app.Finish()
			}

			return result
		}

		return inject
	}

	return injector
}
