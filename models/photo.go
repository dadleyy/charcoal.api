package models

import "os"
import "fmt"
import "database/sql"

type Photo struct {
	Common
	Label string `json:"label"`
	File uint `json:"file"`
	Author sql.NullInt64 `json:"author"`
}

type serializedPhoto struct {
	Photo
	Author interface{} `json:"author"`
	Url string `json:"url"`
}

func (photo Photo) Url() string {
	root := os.Getenv("API_HOME")
	return fmt.Sprintf("%s/photos?filter[id]=eq(%d)", root, photo.ID)
}

func (photo Photo) Type() string {
	return "application/vnd.miritos.photo+json"
}

func (photo Photo) Public() interface{} {
	var author interface{}

	if photo.Author.Valid {
		author = photo.Author.Int64
	} else {
		author = nil
	}

	root := os.Getenv("API_HOME")
	url := fmt.Sprintf("%s/photos/%d/view", root, photo.ID)

	result := serializedPhoto{photo, author, url}
	return result
}
