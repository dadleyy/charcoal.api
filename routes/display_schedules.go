package routes

import "fmt"

import "github.com/sizethree/miritos.api/net"
import "github.com/sizethree/miritos.api/models"

func FindDisplaySchedules(runtime *net.RequestRuntime) error {
	var results []models.DisplaySchedule
	blueprint := runtime.Blueprint()

	total, err := blueprint.Apply(&results, runtime.Database())

	if err != nil {
		runtime.Debugf("bad wschedule lookup query: %s", err.Error())
		return runtime.AddError(fmt.Errorf("BAD_QUERY"))
	}

	for _, item := range results {
		runtime.AddResult(item)
	}

	runtime.SetMeta("total", total)

	return nil
}
