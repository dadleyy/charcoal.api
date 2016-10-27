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

type serializedInstagramPhoto struct {
	InstagramPhoto
	Client interface{} `json:"client"`
}

func (photo InstagramPhoto) Public() interface{} {
	var author interface{}

	author = nil
	if photo.Client.Valid {
		author = photo.Client.Int64
	}

	result := serializedInstagramPhoto{photo, author}
	return result
}
