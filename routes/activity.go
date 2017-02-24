package routes

import "fmt"
import "time"

import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/models"

func FindActivity(runtime *net.RequestRuntime) error {
	var results []models.Activity
	blueprint := runtime.Blueprint()

	total, err := blueprint.Apply(&results)

	if err != nil {
		runtime.Debugf("bad activity lookup query: %s", err.Error())
		return runtime.AddError(fmt.Errorf("BAD_QUERY"))
	}

	for _, item := range results {
		runtime.AddResult(item)
	}

	runtime.SetMeta("total", total)

	return nil
}

func FindLiveActivity(runtime *net.RequestRuntime) error {
	var schedules []models.DisplaySchedule
	today := time.Now()

	conditions := "start < ? AND end > ? AND approval = 'APPROVED'"
	cursor := runtime.Where(conditions, today, today).Select("distinct activity")
	blueprint := runtime.Blueprint(cursor)

	count, err := blueprint.Apply(&schedules)

	// select distinct activities
	if err != nil {
		runtime.Debugf("unable to load current feed: %s", err.Error())
		return runtime.AddError(fmt.Errorf("FAILED"))
	}

	ids := make([]uint, len(schedules))

	for index, item := range schedules {
		ids[index] = item.Activity
	}

	var activities []models.Activity

	if err := runtime.Where("id in (?)", ids).Find(&activities).Error; err != nil {
		runtime.Debugf("unable to load current feed: %s", err.Error())
		return runtime.AddError(fmt.Errorf("FAILED"))
	}

	for _, act := range activities {
		runtime.AddResult(act.Public())
	}

	runtime.SetMeta("total", count)

	return nil
}
