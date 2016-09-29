package middleware

import "errors"
import "github.com/labstack/echo"
import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/context"
import "github.com/sizethree/miritos.api/services"

const ERR_BAD_BEARER = "ERR_BAD_BEARER"

func InjectUser(handler echo.HandlerFunc) echo.HandlerFunc {
	inject := func(ctx echo.Context) error {
		runtime, ok := ctx.(*context.Miritos)

		if ok != true {
			return errors.New("BAD_RUNTIME")
		}

		client := runtime.Client

		if valid := client.ID >= 1; valid != true {
			return handler(runtime)
		}

		bearer := runtime.RequestHeader("X-CLIENT-BEARER-TOKEN")

		if len(bearer) < 1 {
			return handler(runtime)
		}

		clientmgr := services.UserClientManager{runtime.DB}

		if err := clientmgr.Validate(bearer, &client); err != nil {
			runtime.Logger().Error(err)
			return handler(runtime)
		}

		var token models.ClientToken

		if err := runtime.DB.Where("token = ?", bearer).First(&token).Error; err != nil {
			runtime.Logger().Error(err)
			return handler(runtime)
		}

		if err := runtime.DB.First(&runtime.User, token.User).Error; err != nil {
			runtime.Logger().Error(err)
			return handler(runtime)
		}

		return handler(runtime)
	}

	return InjectClient(inject)
}

func RequireUser(handler echo.HandlerFunc) echo.HandlerFunc {
	require := func(ctx echo.Context) error {
		runtime, ok := ctx.(*context.Miritos)

		if ok != true {
			return errors.New("BAD_RUNTIME")
		}

		if valid := runtime.User.ID >= 1; valid != true {
			return runtime.ErrorOut(errors.New(ERR_BAD_BEARER))
		}

		return handler(runtime)
	}

	return InjectUser(require)
}