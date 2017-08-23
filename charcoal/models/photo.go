package models

import "fmt"
import "database/sql"

type Photo struct {
	Common
	Label  string        `json:"label"`
	File   uint          `json:"file"`
	Author sql.NullInt64 `json:"author"`
	Width  int           `json:"width"`
	Height int           `json:"height"`
}

type serializedPhoto struct {
	Photo
	Author interface{} `json:"author"`
	URL    string      `json:"url"`
}

func (photo *Photo) URL() string {
	return fmt.Sprintf("/photos?filter[id]=eq(%d)", photo.ID)
}

func (photo *Photo) Type() string {
	return "application/vnd.miritos.photo+json"
}

func (photo *Photo) Public() interface{} {
	var author interface{}

	author = nil
	if photo.Author.Valid {
		author = photo.Author.Int64
	}

	url := fmt.Sprintf("/photos/%d/view", photo.ID)

	result := serializedPhoto{*photo, author, url}
	return result
}
