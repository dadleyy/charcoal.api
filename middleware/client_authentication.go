package middleware

import "os"
import "github.com/labstack/echo"

func ClientAuthentication(handler echo.HandlerFunc) echo.HandlerFunc {
	secret := os.Getenv("APP_SECRET")

	auth := func(ctx echo.Context) error {
		ctx.Logger().Infof("client authentication (%s): running client auth middleware...", secret)
		return handler(ctx)
	}

	return auth
}
