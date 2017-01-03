package models

import "fmt"
import "database/sql"
import "github.com/jinzhu/gorm"
import "github.com/satori/go.uuid"

type Photo struct {
	Common
	Label  string        `json:"label"`
	File   uint          `json:"file"`
	Author sql.NullInt64 `json:"author"`
	Width  int           `json:"width"`
	Height int           `json:"height"`
	Uuid   string        `json:"uuid"`
}

type serializedPhoto struct {
	Photo
	Author interface{} `json:"author"`
	Url    string      `json:"url"`
}

func (photo *Photo) Url() string {
	return fmt.Sprintf("/photos?filter[id]=eq(%d)", photo.ID)
}

func (photo *Photo) Type() string {
	return "application/vnd.miritos.photo+json"
}

func (photo *Photo) Identifier() string {
	return photo.Uuid
}

func (photo *Photo) BeforeCreate(tx *gorm.DB) error {
	id := uuid.NewV4()
	photo.Uuid = id.String()
	return nil
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
