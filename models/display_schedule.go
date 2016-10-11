package models

import "time"

type DisplaySchedule struct {
	Common
	Activity uint `json:"activity"`
	Start *time.Time `json:"start"`
	End *time.Time `json:"end"`
	Approval string `json:"approval"`
}

func (schedule DisplaySchedule) Public() interface{} {
	return schedule
}
