package middleware

import "errors"
import "strings"
import "encoding/base64"
import "github.com/labstack/echo"
import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/context"

const ERR_BAD_RUNTIME = "BAD_RUNTIME"
const ERR_MISSING_CLIENT_ID = "MISSING_CLIENT_ID"
const ERR_BAD_CLIENT_ID = "BAD_CLIENT_ID"

func ClientAuthentication(handler echo.HandlerFunc) echo.HandlerFunc {
	auth := func(ctx echo.Context) error {
		runtime, ok := ctx.(*context.Miritos)

		if ok != true {
			return errors.New(ERR_BAD_RUNTIME)
		}

		auth := runtime.RequestHeader("X-CLIENT-AUTH")

		if len(auth) < 1 {
			return runtime.ErrorOut(errors.New(ERR_MISSING_CLIENT_ID))
		}

		decoded, err := base64.StdEncoding.DecodeString(auth)

		if err != nil {
			return runtime.ErrorOut(errors.New(ERR_BAD_CLIENT_ID))
		}

		var client models.Client

		parts := strings.Split(string(decoded), ":")

		if len(parts) < 2 || len(parts[0]) < 1 || len(parts[1]) < 1 {
			return runtime.ErrorOut(errors.New(ERR_BAD_CLIENT_ID))
		}

		where := runtime.DB.Where("client_id = ?", parts[0]).Where("client_secret = ?", parts[1])

		if e := where.First(&client).Error; e != nil {
			return runtime.ErrorOut(errors.New(ERR_BAD_CLIENT_ID))
		}

		runtime.Client = client

		return handler(runtime)
	}

	return auth
}
