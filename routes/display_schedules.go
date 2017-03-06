package routes

import "fmt"
import "time"
import "regexp"

import "github.com/albrow/forms"

import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/services"

const APPROVAL_RE = "^APPROVED|REJECTED|PENDING$"

func FindDisplaySchedules(runtime *net.RequestRuntime) error {
	var results []models.DisplaySchedule
	blueprint := runtime.Blueprint()

	total, err := blueprint.Apply(&results)

	if err != nil {
		runtime.Debugf("bad schedule lookup query: %s", err.Error())
		return runtime.AddError(fmt.Errorf("BAD_QUERY"))
	}

	for _, item := range results {
		runtime.AddResult(item)
	}

	runtime.SetMeta("total", total)

	return nil
}

func UpdateDisplaySchedule(runtime *net.RequestRuntime) error {
	id, ok := runtime.IntParam("id")

	if ok != true {
		return runtime.AddError(fmt.Errorf("BAD_ID"))
	}

	manager := services.UserManager{runtime.DB, runtime.Logger}

	if admin := manager.IsAdmin(&runtime.User); admin != true {
		return runtime.AddError(fmt.Errorf("NON_ADMIN"))
	}

	var schedule models.DisplaySchedule

	if err := runtime.First(&schedule, id).Error; err != nil {
		runtime.Debugf("failed lookup: %s", err.Error())
		return runtime.AddError(fmt.Errorf("BAD_SCHEDULE"))
	}

	body, err := forms.Parse(runtime.Request)

	if err != nil {
		runtime.Debugf("unable to parse body: %s", err.Error())
		return runtime.AddError(fmt.Errorf("BAD_DATA_FORMAT"))
	}

	validate := body.Validator()
	updates := make(map[string]interface{})

	if body.KeyExists("approval") {
		re, err := regexp.Compile(APPROVAL_RE)

		if err != nil {
			runtime.Debugf("failed compiling regex: %s", err.Error())
			return runtime.AddError(fmt.Errorf("BAD_DATA_FORMAT"))
		}

		validate.Require("approval")
		validate.Match("approval", re)
		approval := body.Get("approval")
		updates["approval"] = approval

		// if we're updating the approval but a start/end date has not been set, we need to require
		// them in this body.
		if schedule.End == nil && approval == "APPROVED" {
			validate.Require("end")
		}

		if schedule.Start == nil && approval == "APPROVED" {
			validate.Require("start")
		}
	}

	// if the validator picked up errors, add them to the request runtime and then return
	if validate.HasErrors() == true {
		for _, m := range validate.Messages() {
			runtime.AddError(fmt.Errorf(m))
		}

		return fmt.Errorf("BAD_BODY")
	}

	var startd, endd time.Time

	// check to see if we're updating the start or end time. if we are, we need to make sure
	// that the start date is behind the end date.

	if body.KeyExists("start") {
		start := body.Get("start")
		startd, err = time.Parse(time.RFC3339, start)

		if err != nil {
			runtime.Debugf("bad start time: %s", err.Error())
			return runtime.AddError(fmt.Errorf("BAD_START"))
		}

		updates["start"] = &startd

		// if we are not updating the end and there already exists one, validate it
		if end := schedule.End; body.KeyExists("end") == false && end != nil && startd.After(*end) {
			return runtime.AddError(fmt.Errorf("START_BEFORE_END"))
		}
	}

	if body.KeyExists("end") {
		end := body.Get("end")
		endd, err = time.Parse(time.RFC3339, end)

		if err != nil {
			runtime.Debugf("bad end time: %s", err.Error())
			return runtime.AddError(fmt.Errorf("BAD_START"))
		}

		updates["end"] = &endd

		// if we are not updating the start and there already exists one, validate it
		if start := schedule.Start; body.KeyExists("start") == false && start != nil && endd.Before(*start) {
			return runtime.AddError(fmt.Errorf("END_BEFORE_START"))
		}
	}

	// if we're updating both, make sure they're valid
	if body.KeyExists("start") && body.KeyExists("end") && endd.Before(startd) {
		return runtime.AddError(fmt.Errorf("END_BEFORE_START"))
	}

	if err := runtime.Model(&schedule).Updates(updates).Error; err != nil {
		runtime.Debugf("error updating schedule: %s", err.Error())
		return runtime.AddError(fmt.Errorf("FAILED_SAVE"))
	}

	runtime.AddResult(schedule.Public())
	return nil
}
