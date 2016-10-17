package routes

import "fmt"

import "github.com/sizethree/miritos.api/net"
import "github.com/sizethree/miritos.api/models"

func FindActivity(runtime *net.RequestRuntime) error {
	var results []models.Activity
	blueprint := runtime.Blueprint()

	total, err := blueprint.Apply(&results, runtime.Database())

	if err != nil {
		runtime.Debugf("bad activity lookup query: %s", err.Error())
		return runtime.AddError(fmt.Errorf("BAD_QUERY"))
	}

	for _, item := range results {
		runtime.Debugf("adding item %d, actor_url %s", item.ID, item.ObjectUrl)
		runtime.AddResult(item)
	}

	runtime.SetMeta("total", total)

	return nil
}
