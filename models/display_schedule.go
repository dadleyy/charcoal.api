package models

import "time"

type DisplaySchedule struct {
	Common
	Activity uint `json:"activity"`
	Start *time.Time
	End *time.Time
	Approval string `json:"approval"`
}

