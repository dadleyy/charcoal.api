package routes

import "fmt"
import "time"

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

func FindLiveActivity(runtime *net.RequestRuntime) error {
	blueprint := runtime.Blueprint()
	today := time.Now()
	var schedules []models.DisplaySchedule
	count := uint(0)

	offset := blueprint.Limit() * blueprint.Page()

	cursor := runtime.Database().Limit(blueprint.Limit()).Offset(offset)

	// add the clauses that will give us only approved schedules that are currently running
	cursor = cursor.Where("start < ? AND end > ? AND approval = 'APPROVED'", today, today)

	// select distinct activities
	if err := cursor.Select("distinct activity").Find(&schedules).Count(&count).Error; err != nil {
		runtime.Debugf("unable to load current feed: %s", err.Error())
		return runtime.AddError(fmt.Errorf("FAILED"))
	}

	ids := make([]uint, len(schedules))

	for index, item := range schedules {
		ids[index] = item.Activity
	}

	var activities []models.Activity

	if err := runtime.Database().Where("id in (?)", ids).Find(&activities).Error; err != nil {
		runtime.Debugf("unable to load current feed: %s", err.Error())
		return runtime.AddError(fmt.Errorf("FAILED"))
	}

	for _, act := range activities {
		runtime.AddResult(act.Public())
	}

	runtime.SetMeta("total", count)

	return nil
}
