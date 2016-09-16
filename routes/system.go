package routes

import "errors"
import "net/http"

import "github.com/labstack/echo"
import "github.com/sizethree/miritos.api/context"

func System(ectx echo.Context) error {
	_, ok := ectx.(*context.Miritos)
	logger := ectx.Logger()

	if !ok {
		return errors.New("unable to load miritos context")
	}

	logger.Info("HI!");

	return ectx.String(http.StatusOK, "hi!");
}
