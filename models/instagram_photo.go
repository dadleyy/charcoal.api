package models

import "database/sql"

type InstagramPhoto struct {
	Common
	Owner       string        `json:"owner"`
	InstagramID string        `json:"instagram_id"`
	Caption     string        `json:"caption"`
	Photo       uint          `json:"photo"`
	Client      sql.NullInt64 `json:"client"`
}
