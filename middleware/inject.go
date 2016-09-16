package middleware

import "github.com/labstack/echo"
import "github.com/sizethree/miritos.api/context"

func Inject(handler echo.HandlerFunc) echo.HandlerFunc {
	inject := func(ctx echo.Context) error {
		app, err := context.New(ctx);

		if err != nil {
			ctx.Logger().Error(err)
			return err
		}

		return handler(app)
	}

	return inject
}
