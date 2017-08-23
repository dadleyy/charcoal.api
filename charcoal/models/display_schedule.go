package models

import "time"

// DisplaySchedule is a record that represents a time range that an activity record is scheduled to run.
type DisplaySchedule struct {
	Common
	ActivityID uint       `json:"activity_id"`
	Start      *time.Time `json:"start"`
	End        *time.Time `json:"end"`
	Approval   string     `json:"approval"`
}
