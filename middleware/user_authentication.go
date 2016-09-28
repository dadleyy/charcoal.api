package middleware

import "errors"
import "github.com/labstack/echo"
import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/context"
import "github.com/sizethree/miritos.api/services"

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

		clientmgr := services.UserClientManager{runtime.DB}

		if err := clientmgr.Validate(bearer, &client); err != nil {
			runtime.Logger().Error(err)
			return runtime.ErrorOut(errors.New("BAD_BEARER"))
		}

		var token models.ClientToken
		if err := runtime.DB.Where("token = ?", bearer).First(&token).Error; err != nil {
			return runtime.ErrorOut(err)
		}

		if err := runtime.DB.First(&runtime.User, token.User).Error; err != nil {
			return runtime.ErrorOut(err)
		}

		return handler(runtime)
	}

	before := ClientAuthentication(auth)
	return before
}
