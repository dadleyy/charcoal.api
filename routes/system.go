package routes

import "errors"
import "github.com/labstack/echo"
import "github.com/sizethree/miritos.api/context"

func System(ectx echo.Context) error {
	_, ok := ectx.(*context.Runtime)

	if !ok {
		return errors.New("unable to load miritos context")
	}

	return nil
}
