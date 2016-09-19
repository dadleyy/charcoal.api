package routes

import "errors"
import "github.com/labstack/echo"
import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/context"

func FindUser(ectx echo.Context) error {
	runtime, _ := ectx.(*context.Miritos)
	blueprint := runtime.Blueprint()
	var users []models.User

	total, err := blueprint.Apply(&users, runtime.DB)

	if err != nil {
		return err
	}

	runtime.SetMeta("total", total)

	return nil
}

func UpdateUser(ectx echo.Context) error {
	return nil
}

func CreateUser(ectx echo.Context) error {
	runtime, ok := ectx.(*context.Miritos)

	if ok {
		return errors.New("unable to load miritos context")
	}

	runtime.Logger().Infof("attempting to create user...")

	return nil
}
