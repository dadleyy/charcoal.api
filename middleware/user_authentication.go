package middleware

import "errors"
import "github.com/labstack/echo"
import "github.com/sizethree/miritos.api/context"

func UserAuthentication(handler echo.HandlerFunc) echo.HandlerFunc {
	auth := func(ctx echo.Context) error {
		runtime, ok := ctx.(*context.Miritos)

		if ok != true {
			return errors.New("BAD_RUNTIME")
		}

		client := runtime.Client

		if valid := client.ID >= 1; valid != true {
			return runtime.ErrorOut(errors.New("NO_CLIENT"))
		}

		bearer := runtime.RequestHeader("X-CLIENT-BEARER-TOKEN")

		if len(bearer) < 1 {
			return runtime.ErrorOut(errors.New("NO_BEARER_TOKEN"))
		}

		runtime.Logger().Infof("looking up user auth info based on client[%d]", client.ID)

		return nil
	}

	before := ClientAuthentication(auth)
	return before
}
