package middleware

import "errors"
import "strings"
import "encoding/base64"
import "github.com/labstack/echo"
import "github.com/sizethree/miritos.api/context"

const ERR_BAD_RUNTIME = "BAD_RUNTIME"
const ERR_MISSING_CLIENT_ID = "MISSING_CLIENT_ID"
const ERR_BAD_CLIENT_ID = "BAD_CLIENT_ID"

func InjectClient(handler echo.HandlerFunc) echo.HandlerFunc {
	inject := func(ctx echo.Context) error {
		runtime, ok := ctx.(*context.Miritos)

		if ok != true {
			return errors.New(ERR_BAD_RUNTIME)
		}

		auth := runtime.RequestHeader("X-CLIENT-AUTH")

		if len(auth) < 1 {
			return handler(runtime)
		}

		decoded, err := base64.StdEncoding.DecodeString(auth)

		if err != nil {
			return handler(runtime)
		}

		parts := strings.Split(string(decoded), ":")

		if len(parts) < 2 || len(parts[0]) < 1 || len(parts[1]) < 1 {
			return handler(runtime)
		}

		where := runtime.DB.Where("client_id = ?", parts[0]).Where("client_secret = ?", parts[1])

		if e := where.First(&runtime.Client).Error; e != nil {
			runtime.Logger().Error(e)
			return handler(runtime)
		}

		return handler(runtime)
	}

	return inject
}

func RquireClient(handler echo.HandlerFunc) echo.HandlerFunc {
	require := func(ctx echo.Context) error {
		runtime, ok := ctx.(*context.Miritos)

		if ok != true {
			return errors.New(ERR_BAD_RUNTIME)
		}

		if valid := runtime.Client.ID >= 1; valid == false {
			return runtime.ErrorOut(errors.New(ERR_BAD_CLIENT_ID))
		}

		return handler(runtime)
	}

	return InjectClient(require)
}
