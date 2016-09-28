package middleware

import "errors"
import "github.com/labstack/echo"
import "github.com/sizethree/miritos.api/models"
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

		var token models.ClientToken

		where := runtime.DB.Where("token = ?", bearer)

		if err := where.First(&token).Error; err != nil{
			return runtime.ErrorOut(errors.New("NO_BEARER_TOKEN"))
		}

		if token.Client != client.ID {
			return runtime.ErrorOut(errors.New("BAD_BEARER_TOKEN"))
		}

		if err := runtime.DB.First(&runtime.User).Error; err != nil {
			return runtime.ErrorOut(errors.New("INVALID_USER"))
		}

		return handler(runtime)
	}

	before := ClientAuthentication(auth)
	return before
}
