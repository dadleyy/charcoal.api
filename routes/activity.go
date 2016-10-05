package routes

import "fmt"
import "github.com/labstack/echo"

import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/context"

func FindActivity(ectx echo.Context) error {
	runtime, _ := ectx.(*context.Runtime)

	var results []models.Activity
	blueprint := runtime.Blueprint()

	total, err := blueprint.Apply(&results, runtime.DB)

	if err != nil {
		runtime.Logger().Debugf("bad activity lookup query: %s", err.Error())
		return runtime.ErrorOut(fmt.Errorf("BAD_QUERY"))
	}

	for _, item := range results {
		runtime.AddResult(&item)
	}

	runtime.AddMeta("total", total)

	return nil
}
