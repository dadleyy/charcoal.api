package server

import "os"
import "github.com/labstack/echo"
import "github.com/sizethree/miritos.api/db"
import "github.com/sizethree/miritos.api/context"
import "github.com/sizethree/miritos.api/activity"
import "github.com/sizethree/miritos.api/filestore"

type App struct {
	*echo.Echo
	Queue chan activity.Message
	DBConfig db.Config
}

func (app *App) Inject(handler echo.HandlerFunc) echo.HandlerFunc {
	inject := func(ctx echo.Context) error {
		// open a database connection for this request
		client, err := db.Open(app.DBConfig)

		if err != nil {
			ctx.Logger().Error(err)
			return err
		}

		// once finished handling, close db connection
		defer client.Close()

		errors  := make(context.ErrorList, 0)
		meta    := make(context.MetaData)
		results := make(context.ResultList, 0)

		var store context.FileSaver

		switch os.Getenv("FS_ENGINE") {
		case "s3":
			store = filestore.S3FileStore{}
		default:
			store = filestore.TempStore{}
		}

		app := &context.Runtime{
			Context: ctx,
			DB: client,
			Errors: errors,
			Meta: meta,
			Results: results,
			FS: store,
			ActivityStream: app.Queue,
		}

		result := handler(app)

		defer app.Finish()

		if result == nil {
			return nil
		}

		app.Error(result)
		return nil
	}

	return inject
}
