package routes

import "errors"
import "github.com/labstack/echo"
import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/context"

const ERR_BAD_SESSION = "BAD_SESSION"

func PrintAuth(ectx echo.Context) error {
	runtime, ok := ectx.(*context.Miritos)

	if ok != true {
		return runtime.ErrorOut(errors.New(ERR_BAD_SESSION))
	}

	runtime.AddResult(&runtime.User)

	return nil
}

func PrintClientTokens(ectx echo.Context) error {
	runtime, ok := ectx.(*context.Miritos)

	if ok != true {
		return runtime.ErrorOut(errors.New(ERR_BAD_SESSION))
	}

	blueprint := runtime.Blueprint()
	var tokens []models.ClientToken

	total, err := blueprint.Apply(&tokens, runtime.DB)

	if err != nil {
		return err
	}

	runtime.Logger().Infof("looking up auth info")

	for _, token := range tokens {
		runtime.AddResult(&token)
	}

	runtime.AddMeta("total", total)

	return nil
}
