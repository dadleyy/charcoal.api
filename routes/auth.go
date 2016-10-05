package routes

import "fmt"
import "github.com/labstack/echo"
import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/context"

const ERR_BAD_SESSION = "BAD_SESSION"

func PrintAuth(ectx echo.Context) error {
	runtime, ok := ectx.(*context.Runtime)

	if ok != true {
		return fmt.Errorf(ERR_BAD_SESSION)
	}

	runtime.AddResult(&runtime.User)

	return nil
}

func PrintClientTokens(ectx echo.Context) error {
	runtime, ok := ectx.(*context.Runtime)

	if ok != true {
		return fmt.Errorf(ERR_BAD_SESSION)
	}

	if runtime.Client.ID == 0 {
		return runtime.ErrorOut(fmt.Errorf("BAD_CLIENT"))
	}

	blueprint := runtime.Blueprint()
	var tokens []models.ClientToken

	blueprint.Filter("filter[client]", fmt.Sprint("eq(%d)", runtime.Client.ID))

	total, err := blueprint.Apply(&tokens, runtime.DB)

	if err != nil {
		runtime.Logger().Debugf("unable to apply client tokens: %s", err.Error())
		return runtime.ErrorOut(fmt.Errorf("NO_TOKENS"))
	}

	for _, token := range tokens {
		runtime.AddResult(&token)
	}

	runtime.AddMeta("total", total)

	return nil
}
