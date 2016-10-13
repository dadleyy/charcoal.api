package routes
/*

import "fmt"
import "github.com/labstack/echo"

import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/context"

func FindDisplaySchedules(ectx echo.Context) error {
	runtime, _ := ectx.(*context.Runtime)

	var results []models.DisplaySchedule
	blueprint := runtime.Blueprint()

	total, err := blueprint.Apply(&results, runtime.DB)

	if err != nil {
		runtime.Logger().Debugf("bad wschedule lookup query: %s", err.Error())
		return fmt.Errorf("BAD_QUERY")
	}

	for _, item := range results {
		runtime.AddResult(item)
	}

	runtime.AddMeta("total", total)

	return nil
}
*/
